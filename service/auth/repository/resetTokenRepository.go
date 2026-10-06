package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-pointofsale-shared/errors"
	"github.com/jmoiron/sqlx"
)

type resetTokenRepository struct {
	db *sqlx.DB
}

func NewResetTokenRepository(db *sqlx.DB) *resetTokenRepository {
	return &resetTokenRepository{
		db: db,
	}
}

const createResetToken = `INSERT INTO reset_tokens (user_id, token, expiry_date)
VALUES ($1, $2, $3)
RETURNING *`

const deleteResetToken = `DELETE FROM reset_tokens
WHERE user_id = $1`

const getResetToken = `SELECT * FROM reset_tokens
WHERE token = $1`

func (r *resetTokenRepository) FindByToken(ctx context.Context, code string) (*db.ResetToken, error) {
	var token db.ResetToken
	if err := r.db.GetContext(ctx, &token, getResetToken, code); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &token, nil
}

func (r *resetTokenRepository) CreateResetToken(ctx context.Context, req *requests.CreateResetTokenRequest) (*db.ResetToken, error) {
	return r.createResetToken(ctx, r.db, req)
}

// CreateResetTokenInTx persists the reset token inside the given database
// transaction so the caller can commit the token write and its outbox event
// atomically (Phase 6).
func (r *resetTokenRepository) CreateResetTokenInTx(ctx context.Context, tx *sqlx.Tx, req *requests.CreateResetTokenRequest) (*db.ResetToken, error) {
	return r.createResetToken(ctx, tx, req)
}

func (r *resetTokenRepository) createResetToken(ctx context.Context, q sqlx.QueryerContext, req *requests.CreateResetTokenRequest) (*db.ResetToken, error) {
	expiryDate, err := time.Parse("2006-01-02 15:04:05", req.ExpiredAt)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	var token db.ResetToken
	if err := sqlx.GetContext(ctx, q, &token, createResetToken,
		int64(req.UserID),
		req.ResetToken,
		expiryDate,
	); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &token, nil
}

func (r *resetTokenRepository) DeleteResetToken(ctx context.Context, user_id int) error {
	if _, err := r.db.ExecContext(ctx, deleteResetToken, int64(user_id)); err != nil {
		return sharedErrors.ErrInternal.WithInternal(err)
	}
	return nil
}
