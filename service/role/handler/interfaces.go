package handler

import (
	pb "github.com/MamangRust/monolith-pointofsale-pb/role"
	pbuserrole "github.com/MamangRust/monolith-pointofsale-pb/user_role"
)

type RoleQueryHandleGrpc interface {
	pb.RoleQueryServiceServer
}

type RoleCommandHandleGrpc interface {
	pb.RoleCommandServiceServer
}

type UserRoleHandleGrpc interface {
	pbuserrole.UserRoleServiceServer
}
