package orderitem_cache

import (
	"context"

	"github.com/MamangRust/monolith-pointofsale-apigateway/internal/model"
)

type OrderItemQueryCache interface {
	GetCachedOrderItemsAll(ctx context.Context, req *model.FindAllOrderItemInput) (*model.APIResponsePaginationOrderItem, bool)
	SetCachedOrderItemsAll(ctx context.Context, req *model.FindAllOrderItemInput, res *model.APIResponsePaginationOrderItem)

	GetCachedOrderItemActive(ctx context.Context, req *model.FindAllOrderItemInput) (*model.APIResponsePaginationOrderItemDeleteAt, bool)
	SetCachedOrderItemActive(ctx context.Context, req *model.FindAllOrderItemInput, res *model.APIResponsePaginationOrderItemDeleteAt)

	GetCachedOrderItemTrashed(ctx context.Context, req *model.FindAllOrderItemInput) (*model.APIResponsePaginationOrderItemDeleteAt, bool)
	SetCachedOrderItemTrashed(ctx context.Context, req *model.FindAllOrderItemInput, res *model.APIResponsePaginationOrderItemDeleteAt)

	GetCachedOrderItems(ctx context.Context, orderID int) (*model.APIResponsesOrderItem, bool)
	SetCachedOrderItems(ctx context.Context, res *model.APIResponsesOrderItem)
}
