// Package orderitem adapts the OrderItem service gRPC API into the sqlc row
// types used by the consumer repositories. Query and command are exposed as
// separate interfaces so consumers can depend on the narrower one.
package orderitem

import (
	"context"

	pb "github.com/MamangRust/monolith-pointofsale-pb"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter/convert"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/resilience"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	orderitem_errors "github.com/MamangRust/monolith-pointofsale-shared/errors/order_item_errors"
)

// QueryRepository reads order items from the OrderItem service over gRPC.
type QueryRepository interface {
	FindOrderItemByOrder(ctx context.Context, orderID int) ([]*db.OrderItem, error)
	CalculateTotalPrice(ctx context.Context, orderID int) (*int32, error)
}

// CommandRepository is the contract consumers depend on for order-item writes.
type CommandRepository interface {
	CreateOrderItem(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*db.OrderItem, error)
	UpdateOrderItem(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*db.OrderItem, error)
	DeleteOrderItem(ctx context.Context, orderItemID int) error
}

// BulkRepository enumerates order items page by page. It is deliberately a
// separate narrow interface: offline consumers (the stats backfill) need the
// whole collection, while ordinary consumers only ever read one order at a
// time and should not depend on the extra method.
type BulkRepository interface {
	FindAll(ctx context.Context, page, pageSize int) ([]*db.OrderItem, int, error)
}

// Repository implements QueryRepository, CommandRepository and BulkRepository
// over the generated order-item query/command clients, so consumers can share
// one instance across the read and write roles.
type Repository struct {
	query   pb.OrderItemQueryServiceClient
	command pb.OrderItemCommandServiceClient
	guard   *resilience.DependencyGuard
}

// New builds an order-item adapter with both halves wired. Passing zero options
// leaves the guard nil, which makes DependencyGuard.Call a plain passthrough.
func New(query pb.OrderItemQueryServiceClient, command pb.OrderItemCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{query: query, command: command}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewQueryAdapter builds a read-only adapter for consumers that never write.
func NewQueryAdapter(query pb.OrderItemQueryServiceClient, opts ...adapter.GuardOption) QueryRepository {
	return New(query, nil, opts...)
}

// NewCommandAdapter builds a write-only adapter for consumers that never read.
func NewCommandAdapter(command pb.OrderItemCommandServiceClient, opts ...adapter.GuardOption) CommandRepository {
	return New(nil, command, opts...)
}

// SetGuard lets adapter.WithDependencyGuard attach a guard after construction.
func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

var (
	_ QueryRepository   = (*Repository)(nil)
	_ CommandRepository = (*Repository)(nil)
	_ BulkRepository    = (*Repository)(nil)
)

func (r *Repository) FindOrderItemByOrder(ctx context.Context, orderID int) ([]*db.OrderItem, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pb.ApiResponsesOrderItem, error) {
		return r.query.FindOrderItemByOrder(callCtx, &pb.FindByIdOrderItemRequest{Id: int32(orderID)})
	})
	if err != nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder
	}

	res := make([]*db.OrderItem, 0, len(resp.Data))
	for _, item := range resp.Data {
		if item == nil {
			continue
		}
		res = append(res, mapOrderItem(item))
	}

	return res, nil
}

// CalculateTotalPrice derives the order total from its items. It is served by
// the query RPC rather than a dedicated endpoint, so callers treat it as a read.
func (r *Repository) CalculateTotalPrice(ctx context.Context, orderID int) (*int32, error) {
	items, err := r.FindOrderItemByOrder(ctx, orderID)
	if err != nil {
		return nil, orderitem_errors.ErrCalculateTotalPrice
	}

	var total int32
	for _, item := range items {
		if item != nil {
			total += item.Quantity * item.Price
		}
	}

	return &total, nil
}

// FindAll returns one page of active order items plus the total number of rows,
// so callers can walk the collection page by page.
func (r *Repository) FindAll(ctx context.Context, page, pageSize int) ([]*db.OrderItem, int, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pb.ApiResponsePaginationOrderItem, error) {
		return r.query.FindAll(callCtx, &pb.FindAllOrderItemRequest{
			Page:     int32(page),
			PageSize: int32(pageSize),
		})
	})
	if err != nil {
		return nil, 0, orderitem_errors.ErrFindAllOrderItems.WithInternal(err)
	}
	if resp == nil {
		return nil, 0, orderitem_errors.ErrFindAllOrderItems
	}

	items := make([]*db.OrderItem, 0, len(resp.Data))
	for _, item := range resp.Data {
		if item == nil {
			continue
		}
		items = append(items, mapOrderItem(item))
	}

	total := 0
	if resp.Pagination != nil {
		total = int(resp.Pagination.TotalRecords)
	}
	return items, total, nil
}

func (r *Repository) CreateOrderItem(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*db.OrderItem, error) {
	if req == nil {
		return nil, orderitem_errors.ErrCreateOrderItem
	}

	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pb.ApiResponseOrderItem, error) {
		return r.command.CreateOrderItem(callCtx, &pb.CreateOrderItemRecordRequest{
			OrderId:   int32(req.OrderID),
			ProductId: int32(req.ProductID),
			Quantity:  int32(req.Quantity),
			Price:     int32(req.Price),
		})
	})
	if err != nil {
		return nil, orderitem_errors.ErrCreateOrderItem.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, orderitem_errors.ErrCreateOrderItem
	}

	return mapOrderItem(resp.Data), nil
}

func (r *Repository) UpdateOrderItem(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*db.OrderItem, error) {
	if req == nil {
		return nil, orderitem_errors.ErrUpdateOrderItem
	}

	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pb.ApiResponseOrderItem, error) {
		return r.command.UpdateOrderItem(callCtx, &pb.UpdateOrderItemRecordRequest{
			OrderItemId: int32(req.OrderItemID),
			OrderId:     int32(req.OrderID),
			ProductId:   int32(req.ProductID),
			Quantity:    int32(req.Quantity),
			Price:       int32(req.Price),
		})
	})
	if err != nil {
		return nil, orderitem_errors.ErrUpdateOrderItem.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, orderitem_errors.ErrUpdateOrderItem
	}

	return mapOrderItem(resp.Data), nil
}

func (r *Repository) DeleteOrderItem(ctx context.Context, orderItemID int) error {
	_, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pb.ApiResponseOrderItemDelete, error) {
		return r.command.DeleteOrderItem(callCtx, &pb.DeleteOrderItemRecordRequest{OrderItemId: int32(orderItemID)})
	})
	if err != nil {
		return orderitem_errors.ErrDeleteOrderItemPermanent.WithInternal(err)
	}

	return nil
}

// mapOrderItem converts a proto order item into the sqlc row type.
func mapOrderItem(item *pb.OrderItemResponse) *db.OrderItem {
	if item == nil {
		return nil
	}
	return &db.OrderItem{
		OrderItemID: item.Id,
		OrderID:     item.OrderId,
		ProductID:   item.ProductId,
		Quantity:    item.Quantity,
		Price:       item.Price,
		CreatedAt:   convert.Timestamp(item.CreatedAt),
		UpdatedAt:   convert.Timestamp(item.UpdatedAt),
	}
}
