package main

import (
	"github.com/MamangRust/monolith-pointofsale-auth/apps"
	"github.com/MamangRust/monolith-pointofsale-pkg/server"
	"github.com/spf13/viper"
)

func main() {
	port := viper.GetInt("GRPC_AUTH_PORT")
	if port == 0 {
		port = 50051
	}

	srv, err := apps.NewServer(&server.Config{
		ServiceName:    "auth-service",
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
