// Package user adapts the User service gRPC API into the sqlc row types used by
// the consumer repositories. Query and command are exposed as separate
// interfaces so consumers can depend on the narrower one.
package user

import (
	"context"

	pbuser "github.com/MamangRust/monolith-pointofsale-pb/user"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter/convert"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/resilience"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/user_errors"
)

// QueryRepository is the contract consumers depend on for user reads.
type QueryRepository interface {
	FindById(ctx context.Context, id int) (*db.User, error)
	FindByEmail(ctx context.Context, email string) (*db.User, error)
	FindByEmailAndVerify(ctx context.Context, email string) (*db.User, error)
	FindByVerificationCode(ctx context.Context, code string) (*db.User, error)
}

// CommandRepository is the contract consumers depend on for user writes on
// behalf of the auth service: registration (password already hashed) plus the
// two mutations behind VerifyCode and ResetPassword.
type CommandRepository interface {
	CreateUser(ctx context.Context, request *requests.RegisterRequest) (*db.User, error)
	UpdateUserIsVerified(ctx context.Context, userID int, isVerified bool) (*db.User, error)
	UpdateUserPassword(ctx context.Context, userID int, password string) (*db.User, error)
}

// Repository implements QueryRepository and CommandRepository over the
// generated user query/command clients, so consumers can share one instance
// across both roles.
type Repository struct {
	query   pbuser.UserQueryServiceClient
	command pbuser.UserCommandServiceClient
	guard   *resilience.DependencyGuard
}

// New builds a user adapter with both halves wired. Passing zero options leaves
// the guard nil, which makes DependencyGuard.Call a plain passthrough.
func New(query pbuser.UserQueryServiceClient, command pbuser.UserCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{query: query, command: command}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewQueryAdapter builds a read-only adapter for consumers that never write.
func NewQueryAdapter(query pbuser.UserQueryServiceClient, opts ...adapter.GuardOption) QueryRepository {
	return New(query, nil, opts...)
}

// NewCommandAdapter builds a write-only adapter for consumers that never read.
func NewCommandAdapter(command pbuser.UserCommandServiceClient, opts ...adapter.GuardOption) CommandRepository {
	return New(nil, command, opts...)
}

// SetGuard lets adapter.WithDependencyGuard attach a guard after construction.
func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

var (
	_ QueryRepository   = (*Repository)(nil)
	_ CommandRepository = (*Repository)(nil)
)

func (r *Repository) FindById(ctx context.Context, id int) (*db.User, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbuser.ApiResponseUser, error) {
		return r.query.FindById(callCtx, &pbuser.FindByIdUserRequest{Id: int32(id)})
	})
	if err != nil {
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, user_errors.ErrUserNotFound
	}

	return userToModel(resp.Data), nil
}

// FindByEmail returns the account with its password hash so login can verify
// credentials. A missing account is reported as ErrUserNotFound.
func (r *Repository) FindByEmail(ctx context.Context, email string) (*db.User, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbuser.ApiResponseUserWithPassword, error) {
		return r.query.FindByEmail(callCtx, &pbuser.FindByEmailRequest{Email: email})
	})
	if err != nil {
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, user_errors.ErrUserNotFound
	}

	return userWithPasswordToModel(resp.Data), nil
}

// FindByEmailAndVerify is the login-time lookup restricted to verified
// accounts; unverified users must fail credential verification.
func (r *Repository) FindByEmailAndVerify(ctx context.Context, email string) (*db.User, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbuser.ApiResponseUserWithPassword, error) {
		return r.query.FindByEmailAndVerify(callCtx, &pbuser.FindByEmailRequest{Email: email})
	})
	if err != nil {
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, user_errors.ErrUserNotFound
	}

	return userWithPasswordToModel(resp.Data), nil
}

// FindByVerificationCode exchanges the emailed code for the account it belongs
// to, backing the VerifyCode flow.
func (r *Repository) FindByVerificationCode(ctx context.Context, code string) (*db.User, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbuser.ApiResponseUser, error) {
		return r.query.FindByVerificationCode(callCtx, &pbuser.FindByVerificationCodeRequest{VerificationCode: code})
	})
	if err != nil {
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, user_errors.ErrUserNotFound
	}

	return userToModel(resp.Data), nil
}

// CreateUser persists a registration. The password must already be hashed by
// the caller; the verification fields pass through as auth-managed values.
func (r *Repository) CreateUser(ctx context.Context, request *requests.RegisterRequest) (*db.User, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbuser.ApiResponseUser, error) {
		return r.command.CreateUserRecord(callCtx, &pbuser.CreateUserRecordRequest{
			Firstname:        request.FirstName,
			Lastname:         request.LastName,
			Email:            request.Email,
			Password:         request.Password,
			VerificationCode: request.VerifiedCode,
			IsVerified:       request.IsVerified,
		})
	})
	if err != nil {
		return nil, user_errors.ErrCreateUser.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, user_errors.ErrCreateUser
	}

	return userToModel(resp.Data), nil
}

func (r *Repository) UpdateUserIsVerified(ctx context.Context, userID int, isVerified bool) (*db.User, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbuser.ApiResponseUser, error) {
		return r.command.UpdateUserIsVerified(callCtx, &pbuser.UpdateUserIsVerifiedUserRequest{
			Id:         int32(userID),
			IsVerified: isVerified,
		})
	})
	if err != nil {
		return nil, user_errors.ErrUpdateUserVerificationCode.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, user_errors.ErrUpdateUserVerificationCode
	}

	return userToModel(resp.Data), nil
}

func (r *Repository) UpdateUserPassword(ctx context.Context, userID int, password string) (*db.User, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbuser.ApiResponseUser, error) {
		return r.command.UpdateUserPassword(callCtx, &pbuser.UpdateUserPasswordUserRequest{
			Id:       int32(userID),
			Password: password,
		})
	})
	if err != nil {
		return nil, user_errors.ErrUpdateUserPassword.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, user_errors.ErrUpdateUserPassword
	}

	return userToModel(resp.Data), nil
}

func userToModel(u *pbuser.UserResponse) *db.User {
	if u == nil {
		return nil
	}
	return &db.User{
		UserID:    u.Id,
		Firstname: u.Firstname,
		Lastname:  u.Lastname,
		Email:     u.Email,
		CreatedAt: convert.Timestamp(u.CreatedAt),
		UpdatedAt: convert.Timestamp(u.UpdatedAt),
	}
}

func userWithPasswordToModel(u *pbuser.UserResponseWithPassword) *db.User {
	if u == nil {
		return nil
	}
	return &db.User{
		UserID:    u.Id,
		Firstname: u.Firstname,
		Lastname:  u.Lastname,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: convert.Timestamp(u.CreatedAt),
		UpdatedAt: convert.Timestamp(u.UpdatedAt),
	}
}
