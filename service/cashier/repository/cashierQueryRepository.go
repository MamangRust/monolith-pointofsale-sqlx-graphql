package repository

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/cashier_errors"
	"github.com/jmoiron/sqlx"
)

const getCashiers = `SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM cashiers
WHERE deleted_at IS NULL
  AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3`

const getCashiersActive = `SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM cashiers
WHERE deleted_at IS NULL
  AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3`

const getCashiersTrashed = `SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM cashiers
WHERE deleted_at IS NOT NULL
  AND ($1::TEXT IS NULL OR name ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3`

const getCashiersByMerchant = `SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM cashiers
WHERE merchant_id = $1
  AND deleted_at IS NULL
  AND ($2::TEXT IS NULL OR name ILIKE '%' || $2 || '%')
ORDER BY created_at DESC
LIMIT $3 OFFSET $4`

const getCashierById = `SELECT *
FROM cashiers
WHERE cashier_id = $1
  AND deleted_at IS NULL`

type cashierQueryRepository struct {
	db *sqlx.DB
}

func NewCashierQueryRepository(db *sqlx.DB) CashierQueryRepository {
	return &cashierQueryRepository{
		db: db,
	}
}

func (r *cashierQueryRepository) FindAllCashiers(ctx context.Context, req *requests.FindAllCashiers) ([]*db.GetCashiersRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var cashiers []*db.GetCashiersRow
	err := r.db.SelectContext(ctx, &cashiers, getCashiers,
		req.Search,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, cashier_errors.ErrFindAllCashiers
	}

	totalCount := 0
	if len(cashiers) > 0 {
		totalCount = int(cashiers[0].TotalCount)
	}

	return cashiers, &totalCount, nil
}

func (r *cashierQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllCashiers) ([]*db.GetCashiersActiveRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var cashiers []*db.GetCashiersActiveRow
	err := r.db.SelectContext(ctx, &cashiers, getCashiersActive,
		req.Search,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, cashier_errors.ErrFindActiveCashiers
	}

	totalCount := 0
	if len(cashiers) > 0 {
		totalCount = int(cashiers[0].TotalCount)
	}

	return cashiers, &totalCount, nil
}

func (r *cashierQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllCashiers) ([]*db.GetCashiersTrashedRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var cashiers []*db.GetCashiersTrashedRow
	err := r.db.SelectContext(ctx, &cashiers, getCashiersTrashed,
		req.Search,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, cashier_errors.ErrFindTrashedCashiers
	}

	totalCount := 0
	if len(cashiers) > 0 {
		totalCount = int(cashiers[0].TotalCount)
	}

	return cashiers, &totalCount, nil
}

func (r *cashierQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant) ([]*db.GetCashiersByMerchantRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var cashiers []*db.GetCashiersByMerchantRow
	err := r.db.SelectContext(ctx, &cashiers, getCashiersByMerchant,
		int32(req.MerchantID),
		req.Search,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, cashier_errors.ErrFindCashiersByMerchant
	}

	totalCount := 0
	if len(cashiers) > 0 {
		totalCount = int(cashiers[0].TotalCount)
	}

	return cashiers, &totalCount, nil
}

func (r *cashierQueryRepository) FindById(ctx context.Context, cashier_id int) (*db.Cashier, error) {
	var cashier db.Cashier
	err := r.db.GetContext(ctx, &cashier, getCashierById, int32(cashier_id))
	if err != nil {
		return nil, cashier_errors.ErrFindCashierById
	}

	return &cashier, nil
}
