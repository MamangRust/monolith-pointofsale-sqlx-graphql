package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MamangRust/monolith-pointofsale-email/config"
	"github.com/MamangRust/monolith-pointofsale-email/handler"
	"github.com/MamangRust/monolith-pointofsale-email/mailer"
	"github.com/MamangRust/monolith-pointofsale-email/metrics"
	"github.com/MamangRust/monolith-pointofsale-pkg/database"
	"github.com/MamangRust/monolith-pointofsale-pkg/dotenv"
	"github.com/MamangRust/monolith-pointofsale-pkg/emailretry"
	"github.com/MamangRust/monolith-pointofsale-pkg/kafka"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	otel_pkg "github.com/MamangRust/monolith-pointofsale-pkg/otel"
	"github.com/MamangRust/monolith-pointofsale-pkg/outbox"
	"github.com/MamangRust/monolith-pointofsale-shared/convert"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	telemetry := otel_pkg.NewTelemetry(otel_pkg.Config{
		ServiceName:          "email-service",
		ServiceVersion:       "1.0.0",
		Environment:          "production",
		Endpoint:             convert.EnvOr("OTEL_ENDPOINT", "otel-collector:4317"),
		Insecure:             true,
		EnableRuntimeMetrics: os.Getenv("OTEL_ENABLED") != "false",
		Disabled:             os.Getenv("OTEL_ENABLED") == "false",
	})

	if err := telemetry.Init(context.Background()); err != nil {
		log.Fatalf("Error initializing telemetry: %v", err)
	}
	defer func() {
		_ = telemetry.Shutdown(context.Background())
	}()

	logger, err := logger.NewLogger("email-service", telemetry.GetLogger())
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	if err := dotenv.Viper(); err != nil {
		logger.Fatal("Failed to load .env file", zap.Error(err))
	}

	cfg := config.Config{
		KafkaBrokers: []string{viper.GetString("KAFKA_BROKERS")},
		SMTPServer:   viper.GetString("SMTP_SERVER"),
		SMTPPort:     viper.GetInt("SMTP_PORT"),
		SMTPUser:     viper.GetString("SMTP_USER"),
		SMTPPass:     viper.GetString("SMTP_PASS"),
		MaxRetries:   viper.GetInt("EMAIL_MAX_RETRIES"),
		RetryBackoff: viper.GetDuration("EMAIL_RETRY_BACKOFF"),
	}
	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = emailretry.DefaultMaxAttempts
	}
	if cfg.RetryBackoff <= 0 {
		cfg.RetryBackoff = emailretry.DefaultBackoff
	}

	// Register OTel metric instruments after the SDK is initialized so they
	// are bound to the real meter provider (exported via OTLP), not the noop
	// default from package init.
	metrics.Register()

	m := &mailer.Mailer{
		Server:   cfg.SMTPServer,
		Port:     cfg.SMTPPort,
		User:     cfg.SMTPUser,
		Password: cfg.SMTPPass,
	}

	dbPool, err := database.NewClient(logger)
	if err != nil {
		logger.Fatal("Failed to connect to database for consumer inbox", zap.Error(err))
	}
	defer dbPool.Close()

	inbox, err := outbox.NewPostgresInbox(dbPool)
	if err != nil {
		logger.Fatal("Failed to initialize consumer inbox", zap.Error(err))
	}
	myKafka := kafka.NewKafka(logger, cfg.KafkaBrokers)

	h := handler.NewEmailHandlerWithInbox(m, inbox, "email-service-group", myKafka, cfg.RetryBackoff)

	err = myKafka.StartConsumersWithContext(ctx, []string{
		"email-service-topic-auth-register",
		"email-service-topic-auth-forgot-password",
		"email-service-topic-auth-verify-code-success",
		"email-service-topic-merchant-create",
		"email-service-topic-merchant-update-status",
		"email-service-topic-merchant-document-create",
		"email-service-topic-merchant-document-update-status",
		"email-service-topic-transaction-create",
	}, "email-service-group", h)

	if err != nil {
		log.Fatalf("Error starting consumer: %v", err)
	}

	retryH := handler.NewRetryHandler(m, inbox, "email-service-group", myKafka, cfg.MaxRetries, cfg.RetryBackoff)
	if err := myKafka.StartConsumersWithContext(ctx, []string{emailretry.RetryTopic}, emailretry.RetryGroup, retryH); err != nil {
		log.Fatalf("Error starting retry consumer: %v", err)
	}

	logger.Info("Email service started", zap.String("retry_topic", emailretry.RetryTopic), zap.String("dlq_topic", emailretry.DLQTopic))

	<-ctx.Done()
	logger.Info("Shutting down email service")
	if err := myKafka.Close(); err != nil {
		logger.Error("Failed to close Kafka resources", zap.Error(err))
	}
}
