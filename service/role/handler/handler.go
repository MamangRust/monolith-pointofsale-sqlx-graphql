package handler

import (
	pb "github.com/MamangRust/monolith-pointofsale-pb/role"
	pbuserrole "github.com/MamangRust/monolith-pointofsale-pb/user_role"
	"github.com/MamangRust/monolith-pointofsale-pkg/logger"
	"github.com/MamangRust/monolith-pointofsale-role/service"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Role        pb.RoleQueryServiceServer
	RoleCommand pb.RoleCommandServiceServer
	UserRole    pbuserrole.UserRoleServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Role:        NewRoleQueryHandleGrpc(deps.Service, deps.Logger),
		RoleCommand: NewRoleCommandHandleGrpc(deps.Service, deps.Logger),
		UserRole:    NewUserRoleHandleGrpc(deps.Service, deps.Logger),
	}
}
