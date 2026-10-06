package repository

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-pointofsale-shared/errors"
	"github.com/jmoiron/sqlx"
)

type merchantQueryRepository struct {
	db *sqlx.DB
}

func NewMerchantQueryRepository(db *sqlx.DB) MerchantQueryRepository {
	return &merchantQueryRepository{
		db: db,
	}
}

const getMerchants = `SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM merchants
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%' OR contact_email ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3`

const getMerchantsActive = `SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM merchants
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%' OR contact_email ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3`

const getMerchantsTrashed = `SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM merchants
WHERE deleted_at IS NOT NULL
AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%' OR contact_email ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3`

const getMerchantByID = `SELECT *
FROM merchants
WHERE merchant_id = $1
  AND deleted_at IS NULL`

func (r *merchantQueryRepository) FindAllMerchants(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetMerchantsRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var merchants []*db.GetMerchantsRow
	err := r.db.SelectContext(ctx, &merchants, getMerchants,
		req.Search,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	totalCount := 0
	if len(merchants) > 0 {
		totalCount = int(merchants[0].TotalCount)
	}

	return merchants, &totalCount, nil
}

func (r *merchantQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetMerchantsActiveRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var merchants []*db.GetMerchantsActiveRow
	err := r.db.SelectContext(ctx, &merchants, getMerchantsActive,
		req.Search,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	totalCount := 0
	if len(merchants) > 0 {
		totalCount = int(merchants[0].TotalCount)
	}

	return merchants, &totalCount, nil
}

func (r *merchantQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetMerchantsTrashedRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var merchants []*db.GetMerchantsTrashedRow
	err := r.db.SelectContext(ctx, &merchants, getMerchantsTrashed,
		req.Search,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	totalCount := 0
	if len(merchants) > 0 {
		totalCount = int(merchants[0].TotalCount)
	}

	return merchants, &totalCount, nil
}

func (r *merchantQueryRepository) FindById(ctx context.Context, userID int) (*db.Merchant, error) {
	var merchant db.Merchant
	if err := r.db.GetContext(ctx, &merchant, getMerchantByID, int32(userID)); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return &merchant, nil
}
