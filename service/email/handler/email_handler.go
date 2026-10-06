package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/MamangRust/monolith-pointofsale-email/mailer"
	"github.com/MamangRust/monolith-pointofsale-email/metrics"
	"github.com/MamangRust/monolith-pointofsale-pkg/emailretry"
	"github.com/MamangRust/monolith-pointofsale-pkg/event"
	"github.com/MamangRust/monolith-pointofsale-pkg/outbox"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// retryPublisher is implemented by *kafka.Kafka and lets the handlers publish
// retry/DLQ messages with metadata headers (Phase 4).
type retryPublisher interface {
	SendMessageWithHeaders(ctx context.Context, topic, key string, value []byte, headers []sarama.RecordHeader) error
}

type EmailHandler struct {
	Mailer       *mailer.Mailer
	inbox        outbox.ConsumerInbox
	consumerName string
	producer     retryPublisher
	retryBackoff time.Duration
}

func NewEmailHandler(m *mailer.Mailer) *EmailHandler {
	return newEmailHandler(m, nil, "", nil, emailretry.DefaultBackoff)
}

// NewEmailHandlerWithInbox enables durable PostgreSQL-backed deduplication
// and retry-topic offloading for transient SMTP failures.
func NewEmailHandlerWithInbox(m *mailer.Mailer, inbox outbox.ConsumerInbox, consumerName string, producer retryPublisher, retryBackoff time.Duration) *EmailHandler {
	return newEmailHandler(m, inbox, consumerName, producer, retryBackoff)
}

func newEmailHandler(m *mailer.Mailer, inbox outbox.ConsumerInbox, consumerName string, producer retryPublisher, retryBackoff time.Duration) *EmailHandler {
	return &EmailHandler{
		Mailer:       m,
		inbox:        inbox,
		consumerName: consumerName,
		producer:     producer,
		retryBackoff: retryBackoff,
	}
}

func (h *EmailHandler) Setup(_ sarama.ConsumerGroupSession) error {
	// Readiness marker for the trace smoke test (tests/smoke/trace_smoke.sh):
	// logged once the consumer group has joined and is ready to receive
	// messages, so the test knows it is safe to publish a tracked event.
	log.Printf("Consumer session ready — consuming email topics")
	return nil
}
func (h *EmailHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *EmailHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// Continue the distributed trace: extract traceparent/tracestate from
		// the Kafka message headers (injected by the publishing service) and
		// create a child span for this consumed event.
		ctx := otel.GetTextMapPropagator().Extract(context.Background(), kafkaHeaderCarrier(msg.Headers))
		ctx, span := otel.Tracer("email-service").Start(ctx, "consume:"+msg.Topic)
		span.SetAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.destination", msg.Topic),
			attribute.Int64("messaging.kafka.partition", int64(msg.Partition)),
			attribute.Int64("messaging.kafka.offset", msg.Offset),
		)

		var payload event.EmailEnvelope
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			log.Printf("Failed to unmarshal message (topic=%s): %v", msg.Topic, err)
			metrics.EmailInvalid.Add(ctx, 1)
			span.SetAttributes(attribute.Bool("email.invalid", true))
			span.End()
			sess.MarkMessage(msg, "")
			sess.Commit()
			continue
		}

		if !payload.IsValid() {
			log.Printf("Skipping message with invalid email envelope (topic=%s)", msg.Topic)
			metrics.EmailInvalid.Add(ctx, 1)
			span.SetAttributes(attribute.Bool("email.invalid", true))
			span.End()
			sess.MarkMessage(msg, "")
			sess.Commit()
			continue
		}

		// Dedupe on the envelope's event_id (unique per event) rather than the
		// Kafka message key (a business entity ID): two distinct events for the
		// same entity must both be delivered.
		eventKey := fmt.Sprintf("%s:%s", msg.Topic, payload.EventID)
		var (
			reserved, processed bool
			reservationVersion  int64
			err                 error
		)
		if h.inbox != nil {
			reserved, processed, reservationVersion, err = h.inbox.Reserve(ctx, h.consumerName, eventKey, msg.Topic, msg.Partition, msg.Offset)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				span.End()
				return err
			}
			if processed {
				// The side effect was already completed by this consumer.
				span.SetAttributes(attribute.Bool("email.duplicated", true))
				span.End()
				sess.MarkMessage(msg, "")
				sess.Commit()
				continue
			}
			if !reserved {
				err := fmt.Errorf("consumer inbox lease is active for event %s", eventKey)
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				span.End()
				return err
			}
		}

		err = h.Mailer.Send(payload.Email, payload.Subject, payload.Body)
		if err != nil {
			log.Printf("Failed to send email; offloading to retry topic: %v", err)
			metrics.EmailFailed.Add(ctx, 1)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			span.SetAttributes(attribute.Bool("email.failed", true))
			// Release the inbox lease so the retry path can re-claim the event.
			if h.inbox != nil {
				if releaseErr := h.inbox.Release(ctx, h.consumerName, eventKey, reservationVersion, err); releaseErr != nil {
					log.Printf("failed to release consumer inbox lease: %v", releaseErr)
				}
			}

			if h.producer == nil {
				span.End()
				return err
			}
			meta := emailretry.RetryMeta{
				SrcTopic: msg.Topic, SrcPartition: msg.Partition, SrcOffset: msg.Offset,
				Attempt: 1, RetryAt: emailretry.NextRetryAt(1, h.retryBackoff), Reason: err.Error(),
			}
			if pubErr := h.producer.SendMessageWithHeaders(ctx, emailretry.RetryTopic, payload.EventID, msg.Value, emailretry.BuildHeaders(meta)); pubErr != nil {
				span.End()
				return fmt.Errorf("publish to retry topic: %w", pubErr)
			}
			metrics.EmailRetried.Add(ctx, 1)
			span.End()
			sess.MarkMessage(msg, "")
			sess.Commit()
			continue
		}

		metrics.EmailSent.Add(ctx, 1)
		if h.inbox != nil {
			if err := h.inbox.MarkProcessed(ctx, h.consumerName, eventKey, reservationVersion); err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				span.End()
				return err
			}
		}
		span.End()
		sess.MarkMessage(msg, "")
		sess.Commit()
	}
	return nil
}

// kafkaHeaderCarrier adapts sarama record headers to the OTel propagation
// HeaderCarrier interface so trace context can be extracted.
type kafkaHeaderCarrier []*sarama.RecordHeader

func (c kafkaHeaderCarrier) Get(key string) string {
	for _, h := range c {
		if h != nil && string(h.Key) == key {
			return string(h.Value)
		}
	}
	return ""
}

func (c kafkaHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for _, h := range c {
		if h != nil {
			keys = append(keys, string(h.Key))
		}
	}
	return keys
}

func (c kafkaHeaderCarrier) Set(string, string) {}
