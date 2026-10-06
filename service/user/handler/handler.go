package handler

import (
	pb "github.com/MamangRust/monolith-pointofsale-pb/user"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-user/service"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	User        pb.UserQueryServiceServer
	UserCommand pb.UserCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		User:        NewUserQueryHandleGrpc(deps.Service, deps.Logger),
		UserCommand: NewUserCommandHandleGrpc(deps.Service, deps.Logger),
	}
}
