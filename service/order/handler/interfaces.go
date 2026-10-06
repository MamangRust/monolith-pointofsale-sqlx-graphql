package handler

import (
	pb "github.com/MamangRust/monolith-pointofsale-pb/order"
)

type OrderQueryHandleGrpc interface {
	pb.OrderQueryServiceServer
}

type OrderCommandHandleGrpc interface {
	pb.OrderCommandServiceServer
}

type OrderStatsHandleGrpc interface {
	pb.OrderStatsServiceServer
}
