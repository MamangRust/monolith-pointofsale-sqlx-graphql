package handler

import (
	pb "github.com/MamangRust/monolith-pointofsale-pb/product"
)

type ProductQueryHandleGrpc interface {
	pb.ProductQueryServiceServer
}

type ProductCommandHandleGrpc interface {
	pb.ProductCommandServiceServer
}
