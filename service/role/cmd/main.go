package main

import (
	"github.com/MamangRust/monolith-pointofsale-pkg/server"
	"github.com/MamangRust/monolith-pointofsale-role/apps"
	"github.com/spf13/viper"
)

func main() {
	port := viper.GetInt("GRPC_ROLE_ADDR")
	if port == 0 {
		port = 50052
	}

	srv, err := apps.NewServer(&server.Config{
		ServiceName:    "role-service",
		ServiceVersion: "1.0.0",
		Environment:    "production",
		OtelEndpoint:   "otel-collector:4317",
		Port:           port,
	})
	if err != nil {
		panic(err)
	}

	if err := srv.Run(); err != nil {
		panic(err)
	}
}
