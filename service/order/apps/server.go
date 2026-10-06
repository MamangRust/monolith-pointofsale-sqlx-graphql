package apps

import (
	"context"
	"time"

	mencache "github.com/MamangRust/monolith-pointofsale-order/cache"
	"github.com/MamangRust/monolith-pointofsale-order/handler"
	"github.com/MamangRust/monolith-pointofsale-order/repository"
	"github.com/MamangRust/monolith-pointofsale-order/service"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter"
	"github.com/MamangRust/monolith-pointofsale-pkg/resilience"
	"github.com/MamangRust/monolith-pointofsale-pkg/server"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
	"github.com/MamangRust/monolith-pointofsale-shared/convert"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"

	pb "github.com/MamangRust/monolith-pointofsale-pb"
	pbcashier "github.com/MamangRust/monolith-pointofsale-pb/cashier"
	pbmerchant "github.com/MamangRust/monolith-pointofsale-pb/merchant"
	pborder "github.com/MamangRust/monolith-pointofsale-pb/order"
	pbproduct "github.com/MamangRust/monolith-pointofsale-pb/product"
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
	productAddr := convert.EnvOr("GRPC_PRODUCT_ADDR", "localhost:50059")
	orderItemAddr := convert.EnvOr("GRPC_ORDERITEM_ADDR", "localhost:50057")

	srv.Logger.Info("Connecting to gRPC microservices from Order microservice",
		zap.String("cashier", cashierAddr),
		zap.String("merchant", merchantAddr),
		zap.String("product", productAddr),
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

	productConn, err := server.NewGRPCClient(productAddr, resilience.GuardDialOption("product", srv.Logger))
	if err != nil {
		cashierConn.Close()
		merchantConn.Close()
		return nil, err
	}

	orderItemConn, err := server.NewGRPCClient(orderItemAddr, resilience.GuardDialOption("order_item", srv.Logger))
	if err != nil {
		cashierConn.Close()
		merchantConn.Close()
		productConn.Close()
		return nil, err
	}

	go func() {
		<-srv.Ctx.Done()
		srv.Logger.Info("Closing gRPC client connections in Order microservice")
		cashierConn.Close()
		merchantConn.Close()
		productConn.Close()
		orderItemConn.Close()
	}()

	cashierClient := pbcashier.NewCashierQueryServiceClient(cashierConn)
	merchantClient := pbmerchant.NewMerchantQueryServiceClient(merchantConn)

	productQueryClient := pbproduct.NewProductQueryServiceClient(productConn)
	productCommandClient := pbproduct.NewProductCommandServiceClient(productConn)
	orderItemQueryClient := pb.NewOrderItemQueryServiceClient(orderItemConn)
	orderItemCommandClient := pb.NewOrderItemCommandServiceClient(orderItemConn)

	guardCashier := resilience.NewDependencyGuard("cashier", 5, 30, 100, 3*time.Second, srv.Logger)
	guardMerchant := resilience.NewDependencyGuard("merchant", 5, 30, 100, 3*time.Second, srv.Logger)
	guardProduct := resilience.NewDependencyGuard("product", 5, 30, 100, 3*time.Second, srv.Logger)
	guardOrderItem := resilience.NewDependencyGuard("order_item", 5, 30, 100, 3*time.Second, srv.Logger)

	repos := repository.NewRepositories(&repository.Deps{
		Db:                     srv.DB,
		CashierQueryClient:     cashierClient,
		MerchantQueryClient:    merchantClient,
		ProductQueryClient:     productQueryClient,
		ProductCommandClient:   productCommandClient,
		OrderItemQueryClient:   orderItemQueryClient,
		OrderItemCommandClient: orderItemCommandClient,
		Guards: repository.GuardOptions{
			Cashier:   []adapter.GuardOption{adapter.WithDependencyGuard(guardCashier)},
			Merchant:  []adapter.GuardOption{adapter.WithDependencyGuard(guardMerchant)},
			Product:   []adapter.GuardOption{adapter.WithDependencyGuard(guardProduct)},
			OrderItem: []adapter.GuardOption{adapter.WithDependencyGuard(guardOrderItem)},
		},
	})
	traceLoggerObservability := observability.NewTraceLoggerObservability(srv.Logger)

	cacheMetrics, err := observability.NewCacheMetrics("order")
	if err != nil {
		return nil, err
	}

	cacheStore := cache.NewCacheStore(srv.Redis, srv.Logger, cacheMetrics)
	mencacheObj := mencache.NewMencache(cacheStore)

	services := service.NewService(&service.Deps{
		Mencache:      mencacheObj,
		Ctx:           context.Background(),
		Repositories:  repos,
		Logger:        srv.Logger,
		Observability: traceLoggerObservability,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pborder.RegisterOrderQueryServiceServer(gs, handlers.Order)
		pborder.RegisterOrderCommandServiceServer(gs, handlers.OrderCommand)
		pborder.RegisterOrderStatsServiceServer(gs, handlers.OrderStats)
		pborder.RegisterOrderStatsByMerchantServiceServer(gs, handlers.OrderStatsByMerchant)
	}

	return srv, nil
}
