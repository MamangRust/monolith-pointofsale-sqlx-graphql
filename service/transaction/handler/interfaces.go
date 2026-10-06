package handler

import (
	pb "github.com/MamangRust/monolith-pointofsale-pb/transaction"
)

type TransactionQueryHandleGrpc interface {
	pb.TransactionQueryServiceServer
}

type TransactionCommandHandleGrpc interface {
	pb.TransactionCommandServiceServer
}

type TransactionStatsHandleGrpc interface {
	pb.TransactionStatsServiceServer
}
