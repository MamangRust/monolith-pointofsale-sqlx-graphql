package handler

import (
	"github.com/MamangRust/monolith-pointofsale-pb/cashier"
)

type CashierQueryHandleGrpc interface {
	cashier.CashierQueryServiceServer
}

type CashierCommandHandleGrpc interface {
	cashier.CashierCommandServiceServer
}

type CashierStatsHandleGrpc interface {
	cashier.CashierStatsServiceServer
}
