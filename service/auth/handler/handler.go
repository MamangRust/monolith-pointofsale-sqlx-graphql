package handler

import (
	"github.com/MamangRust/monolith-pointofsale-auth/service"
	pb "github.com/MamangRust/monolith-pointofsale-pb"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Auth pb.AuthServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Auth: NewAuthHandleGrpc(
			deps.Service,
			deps.Logger,
		),
	}
}
