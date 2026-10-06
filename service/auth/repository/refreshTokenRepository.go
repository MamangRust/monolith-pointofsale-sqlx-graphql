package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-pointofsale-shared/errors"
	"github.com/jmoiron/sqlx"
)

type refreshTokenRepository struct {
	db *sqlx.DB
}

func NewRefreshTokenRepository(db *sqlx.DB) *refreshTokenRepository {
	return &refreshTokenRepository{
		db: db,
	}
}

const createRefreshToken = `INSERT INTO refresh_tokens (user_id, token, expiration, created_at, updated_at)
VALUES ($1, $2, $3, current_timestamp, current_timestamp)
RETURNING refresh_token_id, user_id, token, expiration, created_at, updated_at, deleted_at`

const findRefreshTokenByToken = `SELECT refresh_token_id, user_id, token, expiration, created_at, updated_at, deleted_at
FROM refresh_tokens
WHERE token = $1 AND deleted_at IS NULL`

const findRefreshTokenByUserId = `SELECT
    refresh_token_id,
    user_id,
    token,
    expiration,
    created_at,
    updated_at,
    deleted_at
FROM
    refresh_tokens
WHERE
    user_id = $1 AND deleted_at IS NULL
ORDER BY
    created_at DESC
LIMIT 1`

const updateRefreshTokenByUserId = `UPDATE refresh_tokens
SET token = $2, expiration = $3, updated_at = current_timestamp
WHERE user_id = $1 AND deleted_at IS NULL
RETURNING *`

const deleteRefreshToken = `DELETE FROM refresh_tokens
WHERE token = $1`

const deleteRefreshTokenByUserId = `DELETE FROM refresh_tokens
WHERE user_id = $1`

func (r *refreshTokenRepository) FindByToken(ctx context.Context, token string) (*db.RefreshToken, error) {
	var refreshToken db.RefreshToken
	if err := r.db.GetContext(ctx, &refreshToken, findRefreshTokenByToken, token); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &refreshToken, nil
}

func (r *refreshTokenRepository) FindByUserId(ctx context.Context, user_id int) (*db.RefreshToken, error) {
	var refreshToken db.RefreshToken
	if err := r.db.GetContext(ctx, &refreshToken, findRefreshTokenByUserId, int32(user_id)); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &refreshToken, nil
}

func (r *refreshTokenRepository) CreateRefreshToken(ctx context.Context, req *requests.CreateRefreshToken) (*db.RefreshToken, error) {
	layout := "2006-01-02 15:04:05"
	expirationTime, err := time.Parse(layout, req.ExpiresAt)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	var refreshToken db.RefreshToken
	if err := r.db.GetContext(ctx, &refreshToken, createRefreshToken,
		int32(req.UserId),
		req.Token,
		expirationTime,
	); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &refreshToken, nil
}

func (r *refreshTokenRepository) UpdateRefreshToken(ctx context.Context, req *requests.UpdateRefreshToken) (*db.RefreshToken, error) {
	layout := "2006-01-02 15:04:05"
	expirationTime, err := time.Parse(layout, req.ExpiresAt)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	var refreshToken db.RefreshToken
	if err := r.db.GetContext(ctx, &refreshToken, updateRefreshTokenByUserId,
		int32(req.UserId),
		req.Token,
		expirationTime,
	); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &refreshToken, nil
}

func (r *refreshTokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	if _, err := r.db.ExecContext(ctx, deleteRefreshToken, token); err != nil {
		return sharedErrors.ErrInternal.WithInternal(err)
	}
	return nil
}

func (r *refreshTokenRepository) DeleteRefreshTokenByUserId(ctx context.Context, user_id int) error {
	if _, err := r.db.ExecContext(ctx, deleteRefreshTokenByUserId, int32(user_id)); err != nil {
		return sharedErrors.ErrInternal.WithInternal(err)
	}
	return nil
}
