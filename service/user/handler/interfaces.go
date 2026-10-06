package handler

import (
	pb "github.com/MamangRust/monolith-pointofsale-pb/user"
)

type UserQueryHandleGrpc interface {
	pb.UserQueryServiceServer
}

type UserCommandHandleGrpc interface {
	pb.UserCommandServiceServer
}
