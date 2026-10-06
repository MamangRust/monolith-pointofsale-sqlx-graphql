package apps

import (
	"context"
	pb "github.com/MamangRust/monolith-pointofsale-pb"

	mencache "github.com/MamangRust/monolith-pointofsale-order-item/cache"
	"github.com/MamangRust/monolith-pointofsale-order-item/handler"
	"github.com/MamangRust/monolith-pointofsale-order-item/repository"
	"github.com/MamangRust/monolith-pointofsale-order-item/service"
	"github.com/MamangRust/monolith-pointofsale-pkg/server"
	"github.com/MamangRust/monolith-pointofsale-shared/cache"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"

	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepositories(srv.DB)
	traceLoggerObservability := observability.NewTraceLoggerObservability(srv.Logger)

	cacheMetrics, err := observability.NewCacheMetrics("order-item")
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
		pb.RegisterOrderItemQueryServiceServer(gs, handlers.OrderItem)
		pb.RegisterOrderItemCommandServiceServer(gs, handlers.OrderItemCommand)
	}

	return srv, nil
}
