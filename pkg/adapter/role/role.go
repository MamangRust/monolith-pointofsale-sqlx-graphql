// Package role adapts the Role service gRPC query API into the sqlc row types
// used by the consumer repositories.
//
// QueryRepository is the read side. The user_roles write side lives in the
// sibling package pkg/adapter/user_role, which is backed by the role service's
// command API.
package role

import (
	"context"

	pbrole "github.com/MamangRust/monolith-pointofsale-pb/role"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter/convert"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/resilience"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/role_errors"
)

// QueryRepository is the contract consumers depend on for role reads.
type QueryRepository interface {
	FindById(ctx context.Context, roleID int) (*db.Role, error)
	FindByName(ctx context.Context, name string) (*db.Role, error)
}

// Repository implements QueryRepository on top of the generated role query client.
type Repository struct {
	client pbrole.RoleQueryServiceClient
	guard  *resilience.DependencyGuard
}

// New builds a role adapter. Passing zero options leaves the guard nil, which
// makes DependencyGuard.Call a plain passthrough.
func New(client pbrole.RoleQueryServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{client: client}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// SetGuard lets adapter.WithDependencyGuard attach a guard after construction.
func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

var _ QueryRepository = (*Repository)(nil)

func (r *Repository) FindById(ctx context.Context, roleID int) (*db.Role, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbrole.ApiResponseRole, error) {
		return r.client.FindByIdRole(callCtx, &pbrole.FindByIdRoleRequest{RoleId: int32(roleID)})
	})
	if err != nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}
	if resp == nil || resp.GetData() == nil {
		return nil, role_errors.ErrRoleNotFound
	}

	return roleToModel(resp.GetData()), nil
}

func (r *Repository) FindByName(ctx context.Context, name string) (*db.Role, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbrole.ApiResponseRole, error) {
		return r.client.FindByNameRole(callCtx, &pbrole.FindByNameRoleRequest{Name: name})
	})
	if err != nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}
	if resp == nil || resp.GetData() == nil {
		return nil, role_errors.ErrRoleNotFound
	}

	return roleToModel(resp.GetData()), nil
}

func roleToModel(role *pbrole.RoleResponse) *db.Role {
	if role == nil {
		return nil
	}

	return &db.Role{
		RoleID:    role.GetId(),
		RoleName:  role.GetName(),
		CreatedAt: convert.Timestamp(role.GetCreatedAt()),
		UpdatedAt: convert.Timestamp(role.GetUpdatedAt()),
	}
}
