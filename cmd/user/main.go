package main

import (
	"log"
	"syscall"

	"github.com/razorpay/foundation"
	"github.com/razorpay/foundation/providers"

	"github.com/razorpay/go-foundation-v2/internal/config"
	"github.com/razorpay/go-foundation-v2/internal/user"
	"github.com/razorpay/go-foundation-v2/internal/user/repo"
	userservice "github.com/razorpay/go-foundation-v2/internal/user/service"
)

func main() {
	config := &config.AppConfig{}
	server, err := foundation.NewServer("./config/user",
		providers.WithCustomConfig(config),
	)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	tel := server.Container().Telemetry()
	dbCollection := server.Container().GetDatabaseCollection()

	primaryDB, err := dbCollection.Get("primary_db")
	if err != nil {
		log.Fatalf("failed to get primary db: %v", err)
	}

	userServer, err := user.New(
		tel,
		userservice.New(
			tel,
			repo.New(tel, primaryDB),
			config.AWS,
		),
	)
	if err != nil {
		log.Fatalf("initialize user service: %v", err)
	}

	if err := server.Start(
		// register grpc handlers
		foundation.WithGRPCHandlers(
			userServer.GRPCHandler,
		),

		// register http handlers (handled via grpc gateway)
		foundation.WithHTTPHandlers(
			userServer.HTTPHandler(server.Context()),
		),

		// database healthchecks
		foundation.WithHealthChecks(
			dbCollection.NewHealthCheck("primary_db", false),
			dbCollection.NewHealthCheck("replica_db", true),
		),

		// additional shutdown signals
		foundation.WithShutdownSignals(
			syscall.SIGHUP,
		),
	); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
