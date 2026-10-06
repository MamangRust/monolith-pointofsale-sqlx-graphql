package handler

import (
	"github.com/MamangRust/monolith-pointofsale-cashier/service"
	"github.com/MamangRust/monolith-pointofsale-pb/cashier"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Cashier                cashier.CashierQueryServiceServer
	CashierCommand         cashier.CashierCommandServiceServer
	CashierStats           cashier.CashierStatsServiceServer
	CashierStatsById       cashier.CashierStatsByIdServiceServer
	CashierStatsByMerchant cashier.CashierStatsByMerchantServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Cashier:                NewCashierQueryHandleGrpc(deps.Service, deps.Logger),
		CashierCommand:         NewCashierCommandHandleGrpc(deps.Service, deps.Logger),
		CashierStats:           NewCashierStatsHandleGrpc(deps.Service, deps.Logger),
		CashierStatsById:       NewCashierStatsByIdHandleGrpc(deps.Service, deps.Logger),
		CashierStatsByMerchant: NewCashierStatsByMerchantHandleGrpc(deps.Service, deps.Logger),
	}
}
