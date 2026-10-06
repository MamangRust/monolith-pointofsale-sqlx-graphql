package repository

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/order_errors"
	"github.com/jmoiron/sqlx"
)

const getOrders = `SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM orders
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL OR order_id::TEXT ILIKE '%' || $1 || '%' OR total_price::TEXT ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3`

const getOrdersActive = `SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM orders
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL OR order_id::TEXT ILIKE '%' || $1 || '%' OR total_price::TEXT ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3`

const getOrdersTrashed = `SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM orders
WHERE deleted_at IS NOT NULL
AND ($1::TEXT IS NULL OR order_id::TEXT ILIKE '%' || $1 || '%' OR total_price::TEXT ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3`

const getOrdersByMerchant = `SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM orders
WHERE
    deleted_at IS NULL
    AND ($1::TEXT IS NULL OR order_id::TEXT ILIKE '%' || $1 || '%' OR total_price::TEXT ILIKE '%' || $1 || '%')
    AND ($4::INT IS NULL OR merchant_id = $4)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3`

const getOrderByID = `SELECT *
FROM orders
WHERE order_id = $1
  AND deleted_at IS NULL`

const getOrderByIDTrashed = `SELECT *
FROM orders
WHERE order_id = $1`

type orderQueryRepository struct {
	db *sqlx.DB
}

func NewOrderQueryRepository(db *sqlx.DB) OrderQueryRepository {
	return &orderQueryRepository{
		db: db,
	}
}

func (r *orderQueryRepository) FindAllOrders(ctx context.Context, req *requests.FindAllOrders) ([]*db.GetOrdersRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var res []*db.GetOrdersRow
	err := r.db.SelectContext(ctx, &res, getOrders,
		req.Search,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, order_errors.ErrFindAllOrders
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	}

	return res, &totalCount, nil
}

func (r *orderQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllOrders) ([]*db.GetOrdersActiveRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var res []*db.GetOrdersActiveRow
	err := r.db.SelectContext(ctx, &res, getOrdersActive,
		req.Search,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, order_errors.ErrFindByActive
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	}

	return res, &totalCount, nil
}

func (r *orderQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllOrders) ([]*db.GetOrdersTrashedRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var res []*db.GetOrdersTrashedRow
	err := r.db.SelectContext(ctx, &res, getOrdersTrashed,
		req.Search,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, order_errors.ErrFindByTrashed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	}

	return res, &totalCount, nil
}

func (r *orderQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllOrderMerchant) ([]*db.GetOrdersByMerchantRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var res []*db.GetOrdersByMerchantRow
	err := r.db.SelectContext(ctx, &res, getOrdersByMerchant,
		req.Search,
		int32(req.PageSize),
		int32(offset),
		int32(req.MerchantID),
	)
	if err != nil {
		return nil, nil, order_errors.ErrFindByMerchant
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	}

	return res, &totalCount, nil
}

func (r *orderQueryRepository) FindById(ctx context.Context, orderID int) (*db.Order, error) {
	var res db.Order
	err := r.db.GetContext(ctx, &res, getOrderByID, int32(orderID))
	if err != nil {
		return nil, order_errors.ErrFindById
	}

	return &res, nil
}

func (r *orderQueryRepository) FindByTrashedId(ctx context.Context, orderID int) (*db.Order, error) {
	var res db.Order
	err := r.db.GetContext(ctx, &res, getOrderByIDTrashed, int32(orderID))
	if err != nil {
		return nil, order_errors.ErrFindByTrashedId
	}

	return &res, nil
}
