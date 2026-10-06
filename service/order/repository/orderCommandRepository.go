package repository

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/order_errors"
	"github.com/jmoiron/sqlx"
)

const createOrder = `INSERT INTO orders (merchant_id, cashier_id, total_price)
VALUES ($1, $2, $3)
RETURNING *`

const updateOrder = `UPDATE orders
SET total_price = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE order_id = $1
  AND deleted_at IS NULL
  RETURNING *`

const trashedOrder = `UPDATE orders
SET
    deleted_at = current_timestamp
WHERE
    order_id = $1
    AND deleted_at IS NULL
    RETURNING *`

const restoreOrder = `UPDATE orders
SET
    deleted_at = NULL
WHERE
    order_id = $1
    AND deleted_at IS NOT NULL
  RETURNING *`

const deleteOrder = `DELETE FROM orders
WHERE order_id = $1
  AND deleted_at IS NULL`

const deleteOrderPermanently = `WITH target AS MATERIALIZED (
    SELECT
        o.order_id,
        o.merchant_id,
        o.cashier_id,
        o.total_price,
        o.created_at,
        o.updated_at,
        o.deleted_at
    FROM orders AS o
    WHERE o.order_id = $1
      AND o.deleted_at IS NOT NULL
), deleted_transactions AS (
    DELETE FROM transactions AS t
    WHERE t.order_id IN (SELECT target.order_id FROM target)
), deleted_order_items AS (
    DELETE FROM order_items AS oi
    WHERE oi.order_id IN (SELECT target.order_id FROM target)
), deleted_order AS (
    DELETE FROM orders AS o
    WHERE o.order_id IN (SELECT target.order_id FROM target)
    RETURNING o.order_id, o.merchant_id, o.cashier_id, o.total_price, o.created_at, o.updated_at, o.deleted_at
)
SELECT
    deleted_order.order_id,
    deleted_order.merchant_id,
    deleted_order.cashier_id,
    deleted_order.total_price,
    deleted_order.created_at,
    deleted_order.updated_at,
    deleted_order.deleted_at
FROM deleted_order`

const getAllOrdersTrashed = `SELECT *
FROM orders
WHERE deleted_at IS NOT NULL
ORDER BY order_id`

const deleteAllPermanentOrders = `WITH target AS MATERIALIZED (
    SELECT order_id FROM orders WHERE deleted_at IS NOT NULL
), deleted_transactions AS (
    DELETE FROM transactions WHERE order_id IN (SELECT order_id FROM target)
), deleted_order_items AS (
    DELETE FROM order_items WHERE order_id IN (SELECT order_id FROM target)
)
DELETE FROM orders
WHERE order_id IN (SELECT order_id FROM target)`

type orderCommandRepository struct {
	db *sqlx.DB
}

func NewOrderCommandRepository(db *sqlx.DB) OrderCommandRepository {
	return &orderCommandRepository{
		db: db,
	}
}

func (r *orderCommandRepository) DeleteOrder(ctx context.Context, orderID int) error {
	if _, err := r.db.ExecContext(ctx, deleteOrder, int32(orderID)); err != nil {
		return order_errors.ErrDeleteOrderPermanent
	}
	return nil
}

func (r *orderCommandRepository) CreateOrder(ctx context.Context, request *requests.CreateOrderRecordRequest) (*db.Order, error) {
	var order db.Order
	err := r.db.GetContext(ctx, &order, createOrder,
		int32(request.MerchantID),
		int32(request.CashierID),
		int64(request.TotalPrice),
	)
	if err != nil {
		return nil, order_errors.ErrCreateOrder
	}

	return &order, nil
}

func (r *orderCommandRepository) FindAllTrashed(ctx context.Context) ([]*db.Order, error) {
	var orders []*db.Order
	err := r.db.SelectContext(ctx, &orders, getAllOrdersTrashed)
	if err != nil {
		return nil, order_errors.ErrFindAllTrashed
	}
	return orders, nil
}

func (r *orderCommandRepository) UpdateOrder(ctx context.Context, request *requests.UpdateOrderRecordRequest) (*db.Order, error) {
	var order db.Order
	err := r.db.GetContext(ctx, &order, updateOrder,
		int32(request.OrderID),
		int64(request.TotalPrice),
	)
	if err != nil {
		return nil, order_errors.ErrUpdateOrder
	}

	return &order, nil
}

func (r *orderCommandRepository) TrashedOrder(ctx context.Context, orderID int) (*db.Order, error) {
	var order db.Order
	err := r.db.GetContext(ctx, &order, trashedOrder, int32(orderID))
	if err != nil {
		return nil, order_errors.ErrTrashedOrder.WithInternal(err)
	}

	return &order, nil
}

func (r *orderCommandRepository) RestoreOrder(ctx context.Context, orderID int) (*db.Order, error) {
	var order db.Order
	err := r.db.GetContext(ctx, &order, restoreOrder, int32(orderID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, order_errors.ErrRestoreOrderNotFound
		}
		return nil, order_errors.ErrRestoreOrder.WithInternal(err)
	}

	return &order, nil
}

func (r *orderCommandRepository) DeleteOrderPermanent(ctx context.Context, orderID int) (bool, error) {
	var deleted db.DeleteOrderPermanentlyRow
	err := r.db.GetContext(ctx, &deleted, deleteOrderPermanently, int32(orderID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, order_errors.ErrDeleteOrderPermanentNotFound
		}
		return false, order_errors.ErrDeleteOrderPermanent.WithInternal(err)
	}

	return true, nil
}

func (r *orderCommandRepository) DeleteAllOrderPermanent(ctx context.Context) (bool, error) {
	if _, err := r.db.ExecContext(ctx, deleteAllPermanentOrders); err != nil {
		return false, order_errors.ErrDeleteAllOrderPermanent.WithInternal(err)
	}
	return true, nil
}
