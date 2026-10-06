package handler

import (
	"github.com/MamangRust/monolith-pointofsale-order/service"
	pb "github.com/MamangRust/monolith-pointofsale-pb/order"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Order                pb.OrderQueryServiceServer
	OrderCommand         pb.OrderCommandServiceServer
	OrderStats           pb.OrderStatsServiceServer
	OrderStatsByMerchant pb.OrderStatsByMerchantServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Order:                NewOrderQueryHandleGrpc(deps.Service, deps.Logger),
		OrderCommand:         NewOrderCommandHandleGrpc(deps.Service, deps.Logger),
		OrderStats:           NewOrderStatsHandleGrpc(deps.Service, deps.Logger),
		OrderStatsByMerchant: NewOrderStatsByMerchantHandleGrpc(deps.Service, deps.Logger),
	}
}
