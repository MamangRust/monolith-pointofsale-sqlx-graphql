package mencache

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
)

type RoleCommandCache interface {
	DeleteCachedRole(ctx context.Context, id int)
	DeleteCachedRoleAllCache(ctx context.Context)
}

type RoleQueryCache interface {
	SetCachedRoles(ctx context.Context, req *requests.FindAllRoles, data []*db.GetRolesRow, total *int)
	SetCachedRoleById(ctx context.Context, data *db.Role)
	SetCachedRoleByUserId(ctx context.Context, userId int, data []*db.Role)
	SetCachedRoleActive(ctx context.Context, req *requests.FindAllRoles, data []*db.GetActiveRolesRow, total *int)
	SetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRoles, data []*db.GetTrashedRolesRow, total *int)

	GetCachedRoles(ctx context.Context, req *requests.FindAllRoles) ([]*db.GetRolesRow, *int, bool)
	GetCachedRoleByUserId(ctx context.Context, userId int) ([]*db.Role, bool)
	GetCachedRoleById(ctx context.Context, id int) (*db.Role, bool)
	GetCachedRoleActive(ctx context.Context, req *requests.FindAllRoles) ([]*db.GetActiveRolesRow, *int, bool)
	GetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRoles) ([]*db.GetTrashedRolesRow, *int, bool)
}
