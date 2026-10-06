// Package user_role adapts the UserRole service gRPC API for services that must
// not write the user_roles table directly.
//
// The role service stays the owner of user_roles; consumers reach it through
// this adapter instead of running SQL against the table themselves.
package user_role

import (
	"context"

	pbuserrole "github.com/MamangRust/monolith-pointofsale-pb/user_role"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/resilience"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	userrole_errors "github.com/MamangRust/monolith-pointofsale-shared/errors/user_role_errors"
)

// CommandRepository is the contract consumers depend on for user_roles writes.
type CommandRepository interface {
	AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*db.UserRole, error)
	RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error
}

// Repository implements CommandRepository on top of the generated user-role
// service client.
type Repository struct {
	client pbuserrole.UserRoleServiceClient
	guard  *resilience.DependencyGuard
}

// New builds a user_roles adapter. Passing zero options leaves the guard nil,
// which makes DependencyGuard.Call a plain passthrough.
func New(client pbuserrole.UserRoleServiceClient, opts ...adapter.GuardOption) *Repository {
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

var _ CommandRepository = (*Repository)(nil)

func (r *Repository) AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*db.UserRole, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbuserrole.ApiResponseUserRole, error) {
		return r.client.AssignRoleToUser(callCtx, &pbuserrole.AssignRoleToUserRequest{
			UserId: int32(req.UserId),
			RoleId: int32(req.RoleId),
		})
	})
	if err != nil {
		return nil, userrole_errors.ErrAssignRoleToUser.WithInternal(err)
	}
	if resp == nil || resp.GetData() == nil {
		return nil, userrole_errors.ErrAssignRoleToUser
	}

	return &db.UserRole{
		UserRoleID: resp.GetData().GetUserRoleId(),
		UserID:     resp.GetData().GetUserId(),
		RoleID:     resp.GetData().GetRoleId(),
	}, nil
}

func (r *Repository) RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error {
	_, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (struct{}, error) {
		_, callErr := r.client.RemoveRoleFromUser(callCtx, &pbuserrole.RemoveRoleFromUserRequest{
			UserId: int32(req.UserId),
			RoleId: int32(req.RoleId),
		})
		return struct{}{}, callErr
	})
	if err != nil {
		return userrole_errors.ErrRemoveRole.WithInternal(err)
	}

	return nil
}
