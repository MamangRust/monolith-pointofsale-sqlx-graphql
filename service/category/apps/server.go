package apps

import (
	"context"

	mencache "github.com/MamangRust/monolith-pointofsale-category/cache"
	"github.com/MamangRust/monolith-pointofsale-category/handler"
	"github.com/MamangRust/monolith-pointofsale-category/repository"
	"github.com/MamangRust/monolith-pointofsale-category/service"
	"github.com/MamangRust/monolith-pointofsale-pkg/server"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"

	pbcategory "github.com/MamangRust/monolith-pointofsale-pb/category"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepositories(srv.DB)
	traceLoggerObservability := observability.NewTraceLoggerObservability(srv.Logger)

	mencacheObj := mencache.NewMencache(srv.CacheStore)

	services := service.NewService(&service.Deps{
		Ctx:           context.Background(),
		Mencache:      mencacheObj,
		Repositories:  repos,
		Logger:        srv.Logger,
		Observability: traceLoggerObservability,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbcategory.RegisterCategoryQueryServiceServer(gs, handlers.Category)
		pbcategory.RegisterCategoryCommandServiceServer(gs, handlers.CategoryCommand)
		pbcategory.RegisterCategoryStatsServiceServer(gs, handlers.CategoryStats)
		pbcategory.RegisterCategoryStatsByIdServiceServer(gs, handlers.CategoryStatsById)
		pbcategory.RegisterCategoryStatsByMerchantServiceServer(gs, handlers.CategoryStatsByMerchant)
	}

	return srv, nil
}
