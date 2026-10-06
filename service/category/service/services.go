package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-pointofsale-category/cache"
	"github.com/MamangRust/monolith-pointofsale-category/repository"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
)

type Service struct {
	CategoryQuery           CategoryQueryService
	CategoryCommand         CategoryCommandService
	CategoryStats           CategoryStatsService
	CategoryStatsById       CategoryStatsByIdService
	CategoryStatsByMerchant CategoryStatsByMerchantService
}

type Deps struct {
	Ctx           context.Context
	Mencache      mencache.Mencache
	Repositories  *repository.Repositories
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	return &Service{
		CategoryQuery: NewCategoryQueryService(&categoryQueryDeps{
			Cache:         deps.Mencache,
			CategoryQuery: deps.Repositories.CategoryQuery,
			Logger:        deps.Logger,
			Observability: deps.Observability,
		}),
		CategoryCommand: NewCategoryCommandService(&categoryCommandDeps{
			Cache:           deps.Mencache,
			CategoryQuery:   deps.Repositories.CategoryQuery,
			CategoryCommand: deps.Repositories.CategoryCommand,
			Logger:          deps.Logger,
			Observability:   deps.Observability,
		}),
		CategoryStats: NewCategoryStatsService(&categoryStatsDeps{
			Cache:                   deps.Mencache,
			CategoryStatsRepository: deps.Repositories.CategoryStats,
			Logger:                  deps.Logger,
			Observability:           deps.Observability,
		}),
		CategoryStatsById: NewCategoryStatsByIdService(&categoryStatsByIdDeps{
			Cache:                       deps.Mencache,
			CategoryStatsByIdRepository: deps.Repositories.CategoryStatsById,
			Logger:                      deps.Logger,
			Observability:               deps.Observability,
		}),
		CategoryStatsByMerchant: NewCategoryStatsByMerchantService(&categoryStatsByMerchantDeps{
			Cache:                             deps.Mencache,
			CategoryStatsByMerchantRepository: deps.Repositories.CategoryStatsByMerchant,
			Logger:                            deps.Logger,
			Observability:                     deps.Observability,
		}),
	}
}
