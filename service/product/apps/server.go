package apps

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-pointofsale-pkg/adapter"
	"github.com/MamangRust/monolith-pointofsale-pkg/resilience"
	"github.com/MamangRust/monolith-pointofsale-pkg/server"

	mencache "github.com/MamangRust/monolith-pointofsale-product/cache"
	"github.com/MamangRust/monolith-pointofsale-product/handler"
	"github.com/MamangRust/monolith-pointofsale-product/repository"
	"github.com/MamangRust/monolith-pointofsale-product/service"
	"github.com/MamangRust/monolith-pointofsale-shared/convert"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"

	pbcategory "github.com/MamangRust/monolith-pointofsale-pb/category"
	pbmerchant "github.com/MamangRust/monolith-pointofsale-pb/merchant"
	pbproduct "github.com/MamangRust/monolith-pointofsale-pb/product"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	categoryAddr := convert.EnvOr("GRPC_CATEGORY_ADDR", "localhost:50054")
	merchantAddr := convert.EnvOr("GRPC_MERCHANT_ADDR", "localhost:50056")

	srv.Logger.Info("Connecting to gRPC microservices from Product microservice",
		zap.String("category", categoryAddr),
		zap.String("merchant", merchantAddr),
	)

	categoryConn, err := server.NewGRPCClient(categoryAddr, resilience.GuardDialOption("category", srv.Logger))
	if err != nil {
		return nil, err
	}

	merchantConn, err := server.NewGRPCClient(merchantAddr, resilience.GuardDialOption("merchant", srv.Logger))
	if err != nil {
		categoryConn.Close()
		return nil, err
	}

	go func() {
		<-srv.Ctx.Done()
		srv.Logger.Info("Closing gRPC client connections in Product microservice")
		categoryConn.Close()
		merchantConn.Close()
	}()

	categoryClient := pbcategory.NewCategoryQueryServiceClient(categoryConn)
	merchantClient := pbmerchant.NewMerchantQueryServiceClient(merchantConn)

	guardCategory := resilience.NewDependencyGuard("category", 5, 30, 100, 3*time.Second, srv.Logger)
	guardMerchant := resilience.NewDependencyGuard("merchant", 5, 30, 100, 3*time.Second, srv.Logger)

	repos := repository.NewRepositories(srv.DB, categoryClient, merchantClient,
		repository.GuardOptions{
			Category: []adapter.GuardOption{adapter.WithDependencyGuard(guardCategory)},
			Merchant: []adapter.GuardOption{adapter.WithDependencyGuard(guardMerchant)},
		},
	)

	mencacheObj := mencache.NewMencache(srv.CacheStore)

	obs := observability.NewTraceLoggerObservability(srv.Logger)

	services := service.NewService(&service.Deps{
		Mencache:      mencacheObj,
		Ctx:           context.Background(),
		Repositories:  repos,
		Logger:        srv.Logger,
		Observability: obs,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbproduct.RegisterProductQueryServiceServer(gs, handlers.Product)
		pbproduct.RegisterProductCommandServiceServer(gs, handlers.ProductCommand)
	}

	return srv, nil
}
