package apps

import (
	"os"
	"time"

	"github.com/MamangRust/monolith-pointofsale-pkg/adapter"
	"github.com/MamangRust/monolith-pointofsale-pkg/kafka"
	"github.com/MamangRust/monolith-pointofsale-pkg/outbox"
	"github.com/MamangRust/monolith-pointofsale-pkg/resilience"
	"github.com/MamangRust/monolith-pointofsale-pkg/server"

	"github.com/MamangRust/monolith-pointofsale-shared/convert"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"

	pb "github.com/MamangRust/monolith-pointofsale-pb"
	pbcashier "github.com/MamangRust/monolith-pointofsale-pb/cashier"
	pbmerchant "github.com/MamangRust/monolith-pointofsale-pb/merchant"
	pborder "github.com/MamangRust/monolith-pointofsale-pb/order"
	pbtransaction "github.com/MamangRust/monolith-pointofsale-pb/transaction"
	mencache "github.com/MamangRust/monolith-pointofsale-transacton/cache"
	"github.com/MamangRust/monolith-pointofsale-transacton/handler"
	"github.com/MamangRust/monolith-pointofsale-transacton/repository"
	"github.com/MamangRust/monolith-pointofsale-transacton/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	cashierAddr := convert.EnvOr("GRPC_CASHIER_ADDR", "localhost:50055")
	merchantAddr := convert.EnvOr("GRPC_MERCHANT_ADDR", "localhost:50056")
	orderAddr := convert.EnvOr("GRPC_ORDER_ADDR", "localhost:50058")
	orderItemAddr := convert.EnvOr("GRPC_ORDERITEM_ADDR", "localhost:50057")

	srv.Logger.Info("Connecting to gRPC microservices from Transaction microservice",
		zap.String("cashier", cashierAddr),
		zap.String("merchant", merchantAddr),
		zap.String("order", orderAddr),
		zap.String("order_item", orderItemAddr),
	)

	cashierConn, err := server.NewGRPCClient(cashierAddr, resilience.GuardDialOption("cashier", srv.Logger))
	if err != nil {
		return nil, err
	}

	merchantConn, err := server.NewGRPCClient(merchantAddr, resilience.GuardDialOption("merchant", srv.Logger))
	if err != nil {
		cashierConn.Close()
		return nil, err
	}

	orderConn, err := server.NewGRPCClient(orderAddr, resilience.GuardDialOption("order", srv.Logger))
	if err != nil {
		cashierConn.Close()
		merchantConn.Close()
		return nil, err
	}

	orderItemConn, err := server.NewGRPCClient(orderItemAddr, resilience.GuardDialOption("order_item", srv.Logger))
	if err != nil {
		cashierConn.Close()
		merchantConn.Close()
		orderConn.Close()
		return nil, err
	}

	go func() {
		<-srv.Ctx.Done()
		srv.Logger.Info("Closing gRPC client connections in Transaction microservice")
		cashierConn.Close()
		merchantConn.Close()
		orderConn.Close()
		orderItemConn.Close()
	}()

	cashierClient := pbcashier.NewCashierQueryServiceClient(cashierConn)
	merchantClient := pbmerchant.NewMerchantQueryServiceClient(merchantConn)
	orderClient := pborder.NewOrderQueryServiceClient(orderConn)
	orderItemClient := pb.NewOrderItemQueryServiceClient(orderItemConn)

	guardCashier := resilience.NewDependencyGuard("cashier", 5, 30, 100, 3*time.Second, srv.Logger)
	guardMerchant := resilience.NewDependencyGuard("merchant", 5, 30, 100, 3*time.Second, srv.Logger)
	guardOrder := resilience.NewDependencyGuard("order", 5, 30, 100, 3*time.Second, srv.Logger)
	guardOrderItem := resilience.NewDependencyGuard("order_item", 5, 30, 100, 3*time.Second, srv.Logger)

	repos := repository.NewRepositories(srv.DB, cashierClient, merchantClient, orderClient, orderItemClient,
		repository.GuardOptions{
			Cashier:   []adapter.GuardOption{adapter.WithDependencyGuard(guardCashier)},
			Merchant:  []adapter.GuardOption{adapter.WithDependencyGuard(guardMerchant)},
			Order:     []adapter.GuardOption{adapter.WithDependencyGuard(guardOrder)},
			OrderItem: []adapter.GuardOption{adapter.WithDependencyGuard(guardOrderItem)},
		},
	)

	var myKafka *kafka.Kafka
	if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		myKafka = kafka.NewKafka(srv.Logger, []string{brokers})
	}

	mencacheObj := mencache.NewMencache(srv.CacheStore)

	traceLoggerObservability := observability.NewTraceLoggerObservability(srv.Logger)

	outboxService := outbox.NewOutboxService(srv.DB, myKafka, srv.Logger)

	services := service.NewService(&service.Deps{
		Mencache:      mencacheObj,
		Repositories:  repos,
		Logger:        srv.Logger,
		Kafka:         myKafka,
		Pool:          srv.DB,
		Outbox:        outboxService,
		Observability: traceLoggerObservability,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbtransaction.RegisterTransactionQueryServiceServer(gs, handlers.Transaction)
		pbtransaction.RegisterTransactionCommandServiceServer(gs, handlers.TransactionCommand)
		pbtransaction.RegisterTransactionStatsServiceServer(gs, handlers.TransactionStats)
		pbtransaction.RegisterTransactionStatsByMerchantServiceServer(gs, handlers.TransactionStatsByMerchant)
	}

	// Start the outbox relay so events committed with the business writes are
	// published to Kafka with durable retry and dead-letter semantics.
	go outboxService.Start(srv.Ctx, outbox.OutboxRelayInterval, outbox.OutboxRelayBatchSize)

	return srv, nil
}
