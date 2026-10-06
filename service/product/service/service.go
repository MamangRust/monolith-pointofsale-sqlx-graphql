package service

import (
	"context"

	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	mencache "github.com/MamangRust/monolith-pointofsale-product/cache"
	"github.com/MamangRust/monolith-pointofsale-product/repository"
	"github.com/MamangRust/monolith-pointofsale-shared/observability"
)

type Service struct {
	ProductQuery   ProductQueryService
	ProductCommand ProductCommandService
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
		ProductQuery:   NewProductQueryService(deps.Mencache, deps.Repositories.ProductQuery, deps.Logger, deps.Observability),
		ProductCommand: NewProductCommandService(deps.Mencache, deps.Repositories.CategoryQuery, deps.Repositories.MerchantQuery, deps.Repositories.ProductQuery, deps.Repositories.ProductCommand, deps.Logger, deps.Observability),
	}
}
