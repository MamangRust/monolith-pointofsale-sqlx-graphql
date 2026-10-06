package repository

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/jmoiron/sqlx"
)

const getOrderItems = `SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM order_items
WHERE deleted_at IS NULL
  AND ($1::TEXT IS NULL OR order_id::TEXT ILIKE '%' || $1 || '%' OR product_id::TEXT ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3`

const getOrderItemsActive = `SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM order_items
WHERE deleted_at IS NULL
  AND ($1::TEXT IS NULL OR order_id::TEXT ILIKE '%' || $1 || '%' OR product_id::TEXT ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3`

const getOrderItemsTrashed = `SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM order_items
WHERE deleted_at IS NOT NULL
  AND ($1::TEXT IS NULL OR order_id::TEXT ILIKE '%' || $1 || '%' OR product_id::TEXT ILIKE '%' || $1 || '%')
ORDER BY deleted_at DESC
LIMIT $2 OFFSET $3`

const getOrderItemsByOrder = `SELECT *
FROM order_items
WHERE order_id = $1
  AND deleted_at IS NULL`

type orderItemQueryRepository struct {
	db *sqlx.DB
}

func NewOrderItemQueryRepository(db *sqlx.DB) *orderItemQueryRepository {
	return &orderItemQueryRepository{
		db: db,
	}
}

func (r *orderItemQueryRepository) FindAllOrderItems(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var orderItems []*db.GetOrderItemsRow
	err := r.db.SelectContext(ctx, &orderItems, getOrderItems,
		req.Search,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, err
	}

	totalCount := 0
	if len(orderItems) > 0 {
		totalCount = int(orderItems[0].TotalCount)
	}

	return orderItems, &totalCount, nil
}

func (r *orderItemQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsActiveRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var orderItems []*db.GetOrderItemsActiveRow
	err := r.db.SelectContext(ctx, &orderItems, getOrderItemsActive,
		req.Search,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, err
	}

	totalCount := 0
	if len(orderItems) > 0 {
		totalCount = int(orderItems[0].TotalCount)
	}

	return orderItems, &totalCount, nil
}

func (r *orderItemQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsTrashedRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var orderItems []*db.GetOrderItemsTrashedRow
	err := r.db.SelectContext(ctx, &orderItems, getOrderItemsTrashed,
		req.Search,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, err
	}

	totalCount := 0
	if len(orderItems) > 0 {
		totalCount = int(orderItems[0].TotalCount)
	}

	return orderItems, &totalCount, nil
}

func (r *orderItemQueryRepository) FindOrderItemByOrder(ctx context.Context, order_id int) ([]*db.OrderItem, error) {
	var orderItems []*db.OrderItem
	err := r.db.SelectContext(ctx, &orderItems, getOrderItemsByOrder, int32(order_id))
	if err != nil {
		return nil, err
	}

	return orderItems, nil
}
