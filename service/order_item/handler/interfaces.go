package handler

import (
	pb "github.com/MamangRust/monolith-pointofsale-pb"
)

type OrderItemHandlerGrpc interface {
	pb.OrderItemQueryServiceServer
}

type OrderItemCommandHandlerGrpc interface {
	pb.OrderItemCommandServiceServer
}
