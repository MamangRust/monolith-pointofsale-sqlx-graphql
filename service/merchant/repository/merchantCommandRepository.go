package repository

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-pointofsale-shared/errors"
	"github.com/jmoiron/sqlx"
)

type merchantCommandRepository struct {
	db *sqlx.DB
}

func NewMerchantCommandRepository(db *sqlx.DB) MerchantCommandRepository {
	return &merchantCommandRepository{
		db: db,
	}
}

const createMerchant = `INSERT INTO merchants (user_id, name, description, address, contact_email, contact_phone, status)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *`

const updateMerchant = `UPDATE merchants
SET name = $2,
    description = $3,
    address = $4,
    contact_email = $5,
    contact_phone = $6,
    status = $7,
    updated_at = CURRENT_TIMESTAMP
WHERE merchant_id = $1
  AND deleted_at IS NULL
  RETURNING *`

const updateMerchantStatus = `UPDATE merchants
SET status = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE merchant_id = $1
  AND deleted_at IS NULL
  RETURNING *`

const trashMerchant = `UPDATE merchants
SET
    deleted_at = current_timestamp
WHERE
    merchant_id = $1
    AND deleted_at IS NULL
    RETURNING *`

const restoreMerchant = `UPDATE merchants
SET
    deleted_at = NULL
WHERE
    merchant_id = $1
    AND deleted_at IS NOT NULL
  RETURNING *`

const deleteMerchantPermanently = `DELETE FROM merchants WHERE merchant_id = $1 AND deleted_at IS NOT NULL`

const restoreAllMerchants = `UPDATE merchants
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL`

const deleteAllPermanentMerchants = `DELETE FROM merchants
WHERE
    deleted_at IS NOT NULL`

func (r *merchantCommandRepository) CreateMerchant(ctx context.Context, request *requests.CreateMerchantRequest) (*db.Merchant, error) {
	return r.createMerchant(ctx, r.db, request)
}

// CreateMerchantInTx persists the merchant inside the given database transaction
// so the caller can commit the business write and its outbox event atomically
// (Phase 6 — transactional outbox).
func (r *merchantCommandRepository) CreateMerchantInTx(ctx context.Context, tx *sqlx.Tx, request *requests.CreateMerchantRequest) (*db.Merchant, error) {
	return r.createMerchant(ctx, tx, request)
}

func (r *merchantCommandRepository) createMerchant(ctx context.Context, q sqlx.QueryerContext, request *requests.CreateMerchantRequest) (*db.Merchant, error) {
	var merchant db.Merchant
	if err := sqlx.GetContext(ctx, q, &merchant, createMerchant,
		int32(request.UserID),
		request.Name,
		&request.Description,
		&request.Address,
		&request.ContactEmail,
		&request.ContactPhone,
		"inactive",
	); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return &merchant, nil
}

// UpdateMerchantStatusInTx updates the merchant status inside the given database
// transaction so the caller can commit the business write and its outbox event
// atomically (Phase 6 — transactional outbox).
func (r *merchantCommandRepository) UpdateMerchantStatusInTx(ctx context.Context, tx *sqlx.Tx, request *requests.UpdateMerchantStatusRequest) (*db.Merchant, error) {
	return r.updateMerchantStatus(ctx, tx, request)
}

func (r *merchantCommandRepository) UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*db.Merchant, error) {
	return r.updateMerchantStatus(ctx, r.db, request)
}

func (r *merchantCommandRepository) updateMerchantStatus(ctx context.Context, q sqlx.QueryerContext, request *requests.UpdateMerchantStatusRequest) (*db.Merchant, error) {
	var merchant db.Merchant
	if err := sqlx.GetContext(ctx, q, &merchant, updateMerchantStatus,
		int32(*request.MerchantID),
		request.Status,
	); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return &merchant, nil
}

func (r *merchantCommandRepository) UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantRequest) (*db.Merchant, error) {
	var merchant db.Merchant
	if err := r.db.GetContext(ctx, &merchant, updateMerchant,
		int32(*request.MerchantID),
		request.Name,
		&request.Description,
		&request.Address,
		&request.ContactEmail,
		&request.ContactPhone,
		request.Status,
	); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return &merchant, nil
}

func (r *merchantCommandRepository) TrashedMerchant(ctx context.Context, merchant_id int) (*db.Merchant, error) {
	var merchant db.Merchant
	if err := r.db.GetContext(ctx, &merchant, trashMerchant, int32(merchant_id)); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return &merchant, nil
}

func (r *merchantCommandRepository) RestoreMerchant(ctx context.Context, merchant_id int) (*db.Merchant, error) {
	var merchant db.Merchant
	if err := r.db.GetContext(ctx, &merchant, restoreMerchant, int32(merchant_id)); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return &merchant, nil
}

func (r *merchantCommandRepository) DeleteMerchantPermanent(ctx context.Context, merchant_id int) (bool, error) {
	if _, err := r.db.ExecContext(ctx, deleteMerchantPermanently, int32(merchant_id)); err != nil {
		return false, sharedErrors.ErrInternal.WithInternal(err)
	}

	return true, nil
}

func (r *merchantCommandRepository) RestoreAllMerchant(ctx context.Context) (bool, error) {
	if _, err := r.db.ExecContext(ctx, restoreAllMerchants); err != nil {
		return false, sharedErrors.ErrInternal.WithInternal(err)
	}
	return true, nil
}

func (r *merchantCommandRepository) DeleteAllMerchantPermanent(ctx context.Context) (bool, error) {
	if _, err := r.db.ExecContext(ctx, deleteAllPermanentMerchants); err != nil {
		return false, sharedErrors.ErrInternal.WithInternal(err)
	}
	return true, nil
}
