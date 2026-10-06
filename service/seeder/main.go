package main

import (
	"context"
	"log"
	"time"

	"github.com/MamangRust/monolith-pointofsale-pkg/database"
	"github.com/MamangRust/monolith-pointofsale-pkg/database/seeder"
	"github.com/MamangRust/monolith-pointofsale-pkg/dotenv"
	"github.com/MamangRust/monolith-pointofsale-pkg/hash"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	otel_pkg "github.com/MamangRust/monolith-pointofsale-pkg/otel"
	"github.com/MamangRust/monolith-pointofsale-shared/convert"
)

func main() {
	telemetry := otel_pkg.NewTelemetry(otel_pkg.Config{
		ServiceName:            "seeder",
		ServiceVersion:         "1.0.0",
		Environment:            "production",
		Endpoint:               convert.EnvOr("OTEL_ENDPOINT", "otel-collector:4317"),
		Insecure:               true,
		EnableRuntimeMetrics:   false,
		RuntimeMetricsInterval: 15 * time.Second,
		Disabled:               true,
	})

	_ = telemetry.Init(context.Background())

	if err := dotenv.Viper(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	l, err := logger.NewLogger("seeder", telemetry.GetLogger())
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	pool, err := database.NewClient(l)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer pool.Close()

	ctx := context.Background()

	s := seeder.NewSeeder(seeder.Deps{
		Db:     pool,
		Ctx:    ctx,
		Logger: l,
		Hash:   hash.NewHashingPassword(),
	})

	if err := s.Run(); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	l.Info("Seeding completed successfully.")
	_ = telemetry.Shutdown(context.Background())
}
