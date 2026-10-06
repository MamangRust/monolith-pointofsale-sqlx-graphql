package handler

import (
	pb "github.com/MamangRust/monolith-pointofsale-pb/category"
)

type CategoryQueryHandleGrpc interface {
	pb.CategoryQueryServiceServer
}

type CategoryCommandHandleGrpc interface {
	pb.CategoryCommandServiceServer
}

type CategoryStatsHandleGrpc interface {
	pb.CategoryStatsServiceServer
}
