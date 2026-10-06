package handler

import (
	pb "github.com/MamangRust/monolith-pointofsale-pb/transaction"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-transacton/service"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Transaction                pb.TransactionQueryServiceServer
	TransactionCommand         pb.TransactionCommandServiceServer
	TransactionStats           pb.TransactionStatsServiceServer
	TransactionStatsByMerchant pb.TransactionStatsByMerchantServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Transaction:                NewTransactionQueryHandleGrpc(deps.Service, deps.Logger),
		TransactionCommand:         NewTransactionCommandHandleGrpc(deps.Service, deps.Logger),
		TransactionStats:           NewTransactionStatsHandleGrpc(deps.Service, deps.Logger),
		TransactionStatsByMerchant: NewTransactionStatsByMerchantHandleGrpc(deps.Service, deps.Logger),
	}
}
