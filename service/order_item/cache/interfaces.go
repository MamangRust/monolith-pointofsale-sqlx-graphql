package mencache

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
)

type OrderItemQueryCache interface {
	GetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsRow, *int, bool)
	SetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems, data []*db.GetOrderItemsRow, total *int)

	GetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsActiveRow, *int, bool)
	SetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems, data []*db.GetOrderItemsActiveRow, total *int)

	GetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsTrashedRow, *int, bool)
	SetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems, data []*db.GetOrderItemsTrashedRow, total *int)

	GetCachedOrderItems(ctx context.Context, orderID int) ([]*db.OrderItem, bool)
	SetCachedOrderItems(ctx context.Context, data []*db.OrderItem)
}

type OrderItemCommandCache interface {
	// DeleteCachedOrderItems invalidates the cached item list for one order.
	DeleteCachedOrderItems(ctx context.Context, orderID int)
	// DeleteCachedOrderItemsAll invalidates every cached order-item list. Used
	// when the caller only knows the order_item id, not its order.
	DeleteCachedOrderItemsAll(ctx context.Context)
}
