package handler

import (
	"github.com/MamangRust/monolith-pointofsale-order-item/service"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	OrderItem        OrderItemHandlerGrpc
	OrderItemCommand OrderItemCommandHandlerGrpc
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		OrderItem: NewOrderItemHandlerGrpc(
			deps.Service.OrderItemQuery,
			deps.Logger,
		),
		OrderItemCommand: NewOrderItemCommandHandlerGrpc(
			deps.Service.OrderItemCommand,
			deps.Logger,
		),
	}
}
