package apps

import (
	"context"
	"os"
	"time"

	mencache "github.com/MamangRust/monolith-pointofsale-merchant/cache"
	"github.com/MamangRust/monolith-pointofsale-merchant/handler"
	"github.com/MamangRust/monolith-pointofsale-merchant/repository"
	"github.com/MamangRust/monolith-pointofsale-merchant/service"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter"
	"github.com/MamangRust/monolith-pointofsale-pkg/kafka"
	"github.com/MamangRust/monolith-pointofsale-pkg/outbox"
	"github.com/MamangRust/monolith-pointofsale-pkg/resilience"
	"github.com/MamangRust/monolith-pointofsale-pkg/server"

	"github.com/MamangRust/monolith-pointofsale-shared/observability"

	pbmerchant "github.com/MamangRust/monolith-pointofsale-pb/merchant"
	pbmerchant_document "github.com/MamangRust/monolith-pointofsale-pb/merchant_document"
	pbuser "github.com/MamangRust/monolith-pointofsale-pb/user"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	userAddr := os.Getenv("GRPC_USER_ADDR")
	if userAddr == "" {
		userAddr = "localhost:50053"
	}

	srv.Logger.Info("Connecting to User service via gRPC", zap.String("addr", userAddr))
	userConn, err := server.NewGRPCClient(userAddr, resilience.GuardDialOption("user", srv.Logger))
	if err != nil {
		return nil, err
	}

	go func() {
		<-srv.Ctx.Done()
		srv.Logger.Info("Closing merchant service remote gRPC connections")
		userConn.Close()
	}()

	userClient := pbuser.NewUserQueryServiceClient(userConn)

	repos := repository.NewRepositories(srv.DB, userClient,
		repository.GuardOptions{
			User: []adapter.GuardOption{adapter.WithDependencyGuard(resilience.NewDependencyGuard("user", 5, 30, 100, 3*time.Second, srv.Logger))},
		},
	)

	var myKafka *kafka.Kafka
	if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		myKafka = kafka.NewKafka(srv.Logger, []string{brokers})
	}
	traceLoggerObservability := observability.NewTraceLoggerObservability(srv.Logger)

	mencacheObj := mencache.NewMencache(srv.CacheStore)

	outboxService := outbox.NewOutboxService(srv.DB, myKafka, srv.Logger)

	services := service.NewService(&service.Deps{
		Ctx:           context.Background(),
		Mencache:      mencacheObj,
		Kafka:         myKafka,
		Repositories:  repos,
		Pool:          srv.DB,
		Outbox:        outboxService,
		Logger:        srv.Logger,
		Observability: traceLoggerObservability,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbmerchant.RegisterMerchantQueryServiceServer(gs, handlers.Merchant)
		pbmerchant.RegisterMerchantCommandServiceServer(gs, handlers.MerchantCommand)
		pbmerchant_document.RegisterMerchantDocumentQueryServiceServer(gs, handlers.MerchantDocument)
		pbmerchant_document.RegisterMerchantDocumentCommandServiceServer(gs, handlers.MerchantDocumentCommand)
	}

	// Start the outbox relay so events committed with the business writes are
	// published to Kafka with durable retry and dead-letter semantics.
	go outboxService.Start(srv.Ctx, outbox.OutboxRelayInterval, outbox.OutboxRelayBatchSize)

	return srv, nil
}
