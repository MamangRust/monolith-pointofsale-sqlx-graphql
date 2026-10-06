package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type Kafka struct {
	logger   logger.LoggerInterface
	producer sarama.SyncProducer
	brokers  []string
}

func NewKafka(logger logger.LoggerInterface, brokers []string) *Kafka {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}

	log.Println("Kafka producer connected successfully")

	return &Kafka{
		producer: producer,
		brokers:  brokers,
		logger:   logger,
	}
}

// SendMessage publishes a message to the given topic, carrying the active
// OpenTelemetry trace context (traceparent/tracestate) as Kafka message
// headers so consumers can continue the same trace.
func (k *Kafka) SendMessage(ctx context.Context, topic string, key string, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	if len(carrier) > 0 {
		for kk, vv := range carrier {
			msg.Headers = append(msg.Headers, sarama.RecordHeader{
				Key:   []byte(kk),
				Value: []byte(vv),
			})
		}
	}

	partition, offset, err := k.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}

// SendMessageWithHeaders publishes a message with explicit record headers in
// addition to the payload, still injecting the active OpenTelemetry trace
// context. Used by the email service to attach retry/DLQ metadata (Phase 4)
// while keeping the payload an unchanged envelope.
func (k *Kafka) SendMessageWithHeaders(ctx context.Context, topic string, key string, value []byte, headers []sarama.RecordHeader) error {
	msg := &sarama.ProducerMessage{
		Topic:   topic,
		Key:     sarama.StringEncoder(key),
		Value:   sarama.ByteEncoder(value),
		Headers: headers,
	}

	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	if len(carrier) > 0 {
		for kk, vv := range carrier {
			msg.Headers = append(msg.Headers, sarama.RecordHeader{
				Key:   []byte(kk),
				Value: []byte(vv),
			})
		}
	}

	partition, offset, err := k.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}

// StartConsumers starts a consumer group with Sarama auto-commit disabled.
// The handler must call session.Commit() explicitly after reaching a terminal
// outcome (e.g. an email was sent or the message was rejected as
// invalid/duplicate); messages that fail transiently stay uncommitted and are
// redelivered on the next rebalance (at-least-once).
func (k *Kafka) StartConsumers(topics []string, groupID string, handler sarama.ConsumerGroupHandler) error {
	return k.StartConsumersWithContext(context.Background(), topics, groupID, handler)
}

// StartConsumersWithContext starts a consumer group that stops cleanly when
// ctx is cancelled: the consume loop returns and the group is closed. Same
// contract as StartConsumers (manual commit, start from the oldest offset).
func (k *Kafka) StartConsumersWithContext(ctx context.Context, topics []string, groupID string, handler sarama.ConsumerGroupHandler) error {
	if ctx == nil {
		ctx = context.Background()
	}
	config := sarama.NewConfig()
	config.Consumer.Offsets.AutoCommit.Enable = false
	// Start fresh consumer groups at the oldest available offset so retry
	// messages published while the service was offline are not skipped.
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true

	consumerGroup, err := sarama.NewConsumerGroup(k.brokers, groupID, config)
	if err != nil {
		return err
	}

	go func() {
		defer consumerGroup.Close()
		for {
			if err := consumerGroup.Consume(ctx, topics, handler); err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Printf("Error from consumer: %v", err)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case err, ok := <-consumerGroup.Errors():
				if !ok {
					return
				}
				log.Printf("Consumer group error: %v", err)
			}
		}
	}()

	return nil
}

// Close releases the Kafka producer. Consumers stop via their context; Close
// is idempotent and safe to call more than once.
func (k *Kafka) Close() error {
	if k == nil || k.producer == nil {
		return nil
	}
	err := k.producer.Close()
	k.producer = nil
	return err
}
