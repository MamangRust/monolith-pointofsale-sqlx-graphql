package repository

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/user_errors"
	"github.com/jmoiron/sqlx"
)

const createUser = `INSERT INTO
    users (
        firstname,
        lastname,
        email,
        password,
        verification_code,
        is_verified,
        created_at,
        updated_at
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        current_timestamp,
        current_timestamp
    ) RETURNING *`

const updateUser = `UPDATE users
SET
    firstname = $2,
    lastname = $3,
    email = $4,
    password = $5,
    updated_at = current_timestamp
WHERE
    user_id = $1
    AND deleted_at IS NULL
    RETURNING *`

const createUserRecord = `INSERT INTO
    users (
        firstname,
        lastname,
        email,
        password,
        verification_code,
        is_verified,
        created_at,
        updated_at
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        current_timestamp,
        current_timestamp
    ) RETURNING *`

const updateUserIsVerified = `UPDATE users
SET
    is_verified = $2,
    updated_at = current_timestamp
WHERE
    user_id = $1
    AND deleted_at IS NULL
    RETURNING *`

const updateUserPassword = `UPDATE users
SET
    password = $2,
    updated_at = current_timestamp
WHERE
    user_id = $1
    AND deleted_at IS NULL
    RETURNING *`

const trashUser = `UPDATE users
SET
    deleted_at = current_timestamp
WHERE
    user_id = $1
    AND deleted_at IS NULL
    RETURNING *`

const restoreUser = `UPDATE users
SET
    deleted_at = NULL
WHERE
    user_id = $1
    AND deleted_at IS NOT NULL
    RETURNING *`

const deleteUserPermanently = `DELETE FROM users WHERE user_id = $1 AND deleted_at IS NOT NULL`

const restoreAllUsers = `UPDATE users
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL`

const deleteAllPermanentUsers = `DELETE FROM users
WHERE
    deleted_at IS NOT NULL`

type userCommandRepository struct {
	db *sqlx.DB
}

func NewUserCommandRepository(db *sqlx.DB) *userCommandRepository {
	return &userCommandRepository{
		db: db,
	}
}

func (r *userCommandRepository) CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*db.User, error) {
	var user db.User
	err := r.db.GetContext(ctx, &user, createUser,
		request.FirstName,
		request.LastName,
		request.Email,
		request.Password,
		"",
		nil,
	)
	if err != nil {
		return nil, user_errors.ErrCreateUser
	}

	return &user, nil
}

func (r *userCommandRepository) UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*db.User, error) {
	var user db.User
	err := r.db.GetContext(ctx, &user, updateUser,
		int32(*request.UserID),
		request.FirstName,
		request.LastName,
		request.Email,
		request.Password,
	)
	if err != nil {
		return nil, user_errors.ErrUpdateUser
	}

	return &user, nil
}

// CreateUserRecord persists a registration coming from the auth service. The
// password is already hashed by the caller, and the verification fields are
// auth-managed, so they are written verbatim.
func (r *userCommandRepository) CreateUserRecord(ctx context.Context, request *requests.RegisterRequest) (*db.User, error) {
	var user db.User
	err := r.db.GetContext(ctx, &user, createUserRecord,
		request.FirstName,
		request.LastName,
		request.Email,
		request.Password,
		request.VerifiedCode,
		request.IsVerified,
	)
	if err != nil {
		return nil, user_errors.ErrCreateUser
	}

	return &user, nil
}

func (r *userCommandRepository) UpdateUserIsVerified(ctx context.Context, user_id int, is_verified bool) (*db.User, error) {
	var user db.User
	err := r.db.GetContext(ctx, &user, updateUserIsVerified, int32(user_id), is_verified)
	if err != nil {
		return nil, user_errors.ErrUpdateUserVerificationCode
	}

	return &user, nil
}

func (r *userCommandRepository) UpdateUserPassword(ctx context.Context, user_id int, password string) (*db.User, error) {
	var user db.User
	err := r.db.GetContext(ctx, &user, updateUserPassword, int32(user_id), password)
	if err != nil {
		return nil, user_errors.ErrUpdateUserPassword
	}

	return &user, nil
}

func (r *userCommandRepository) TrashedUser(ctx context.Context, user_id int) (*db.User, error) {
	var user db.User
	err := r.db.GetContext(ctx, &user, trashUser, int32(user_id))
	if err != nil {
		return nil, user_errors.ErrTrashedUser
	}

	return &user, nil
}

func (r *userCommandRepository) RestoreUser(ctx context.Context, user_id int) (*db.User, error) {
	var user db.User
	err := r.db.GetContext(ctx, &user, restoreUser, int32(user_id))
	if err != nil {
		return nil, user_errors.ErrRestoreUser
	}

	return &user, nil
}

func (r *userCommandRepository) DeleteUserPermanent(ctx context.Context, user_id int) (bool, error) {
	if _, err := r.db.ExecContext(ctx, deleteUserPermanently, int32(user_id)); err != nil {
		return false, user_errors.ErrDeleteUserPermanent
	}

	return true, nil
}

func (r *userCommandRepository) RestoreAllUser(ctx context.Context) (bool, error) {
	if _, err := r.db.ExecContext(ctx, restoreAllUsers); err != nil {
		return false, user_errors.ErrRestoreAllUsers
	}

	return true, nil
}

func (r *userCommandRepository) DeleteAllUserPermanent(ctx context.Context) (bool, error) {
	if _, err := r.db.ExecContext(ctx, deleteAllPermanentUsers); err != nil {
		return false, user_errors.ErrDeleteAllUsers
	}
	return true, nil
}
