package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-pointofsale-cashier/cache"
	"github.com/MamangRust/monolith-pointofsale-cashier/repository"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
)

type Service struct {
	CashierQuery           CashierQueryService
	CashierCommand         CashierCommandService
	CashierStats           CashierStatsService
	CashierStatsById       CashierStatsByIdService
	CashierStatsByMerchant CashierStatsByMerchant
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
		CashierQuery: NewCashierQueryService(&cashierQueryDeps{
			Cache:         deps.Mencache,
			CashierQuery:  deps.Repositories.CashierQuery,
			Logger:        deps.Logger,
			Observability: deps.Observability,
		}),
		CashierCommand: NewCashierCommandService(&cashierCommandDeps{
			Cache:          deps.Mencache,
			MerchantQuery:  deps.Repositories.MerchantQuery,
			UserQuery:      deps.Repositories.UserQuery,
			CashierCommand: deps.Repositories.CashierCommand,
			Logger:         deps.Logger,
			Observability:  deps.Observability,
		}),
		CashierStats: NewCashierStatsService(&cashierStatsDeps{
			Cache:         deps.Mencache,
			CashierStats:  deps.Repositories.CashierStats,
			Logger:        deps.Logger,
			Observability: deps.Observability,
		}),
		CashierStatsById: NewCashierStatsByIdService(&cashierStatsByIdDeps{
			Cache:         deps.Mencache,
			CashierStats:  deps.Repositories.CashierStatsById,
			Logger:        deps.Logger,
			Observability: deps.Observability,
		}),
		CashierStatsByMerchant: NewCashierStatsByMerchantService(&cashierStatsByMerchantDeps{
			Cache:         deps.Mencache,
			CashierStats:  deps.Repositories.CashierStatsByMerchant,
			Logger:        deps.Logger,
			Observability: deps.Observability,
		}),
	}
}
