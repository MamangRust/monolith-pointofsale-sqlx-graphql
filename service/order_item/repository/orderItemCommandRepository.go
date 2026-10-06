package repository

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	orderitem_errors "github.com/MamangRust/monolith-pointofsale-shared/errors/order_item_errors"
	"github.com/jmoiron/sqlx"
)

const createOrderItem = `INSERT INTO order_items (order_id, product_id, quantity, price)
VALUES ($1, $2, $3, $4)
RETURNING *`

const updateOrderItem = `UPDATE order_items
SET product_id = $2,
    quantity = $3,
    price = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE order_item_id = $1
  AND deleted_at IS NULL
  RETURNING *`

const deleteOrderItem = `DELETE FROM order_items
WHERE order_item_id = $1
  AND deleted_at IS NULL`

type orderItemCommandRepository struct {
	db *sqlx.DB
}

func NewOrderItemCommandRepository(db *sqlx.DB) OrderItemCommandRepository {
	return &orderItemCommandRepository{
		db: db,
	}
}

func (r *orderItemCommandRepository) CreateOrderItem(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*db.OrderItem, error) {
	var item db.OrderItem
	err := r.db.GetContext(ctx, &item, createOrderItem,
		int32(req.OrderID),
		int32(req.ProductID),
		int32(req.Quantity),
		int32(req.Price),
	)
	if err != nil {
		return nil, orderitem_errors.ErrCreateOrderItem
	}

	return &item, nil
}

func (r *orderItemCommandRepository) UpdateOrderItem(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*db.OrderItem, error) {
	var item db.OrderItem
	err := r.db.GetContext(ctx, &item, updateOrderItem,
		int32(req.OrderItemID),
		int32(req.ProductID),
		int32(req.Quantity),
		int32(req.Price),
	)
	if err != nil {
		return nil, orderitem_errors.ErrUpdateOrderItem
	}

	return &item, nil
}

func (r *orderItemCommandRepository) DeleteOrderItem(ctx context.Context, orderItemID int) error {
	if _, err := r.db.ExecContext(ctx, deleteOrderItem, int32(orderItemID)); err != nil {
		return orderitem_errors.ErrDeleteOrderItemPermanent
	}
	return nil
}
