package handler

import (
	"github.com/MamangRust/monolith-pointofsale-category/service"
	pb "github.com/MamangRust/monolith-pointofsale-pb/category"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Category                pb.CategoryQueryServiceServer
	CategoryCommand         pb.CategoryCommandServiceServer
	CategoryStats           pb.CategoryStatsServiceServer
	CategoryStatsById       pb.CategoryStatsByIdServiceServer
	CategoryStatsByMerchant pb.CategoryStatsByMerchantServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Category:                NewCategoryQueryHandleGrpc(deps.Service, deps.Logger),
		CategoryCommand:         NewCategoryCommandHandleGrpc(deps.Service, deps.Logger),
		CategoryStats:           NewCategoryStatsHandleGrpc(deps.Service, deps.Logger),
		CategoryStatsById:       NewCategoryStatsByIdHandleGrpc(deps.Service, deps.Logger),
		CategoryStatsByMerchant: NewCategoryStatsByMerchantHandleGrpc(deps.Service, deps.Logger),
	}
}
