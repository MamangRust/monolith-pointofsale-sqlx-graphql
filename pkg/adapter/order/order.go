// Package order adapts the Order service gRPC API into the sqlc row types used
// by the consumer repositories.
package order

import (
	"context"

	pborder "github.com/MamangRust/monolith-pointofsale-pb/order"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter/convert"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/resilience"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/order_errors"
)

// QueryRepository reads a single order from the Order service over gRPC.
type QueryRepository interface {
	FindById(ctx context.Context, id int) (*db.Order, error)
}

// Repository implements QueryRepository over a gRPC client.
type Repository struct {
	client pborder.OrderQueryServiceClient
	guard  *resilience.DependencyGuard
}

// New builds an order adapter. Passing zero options leaves the guard nil, which
// makes DependencyGuard.Call a plain passthrough.
func New(client pborder.OrderQueryServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{client: client}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// SetGuard lets adapter.WithDependencyGuard attach a guard after construction.
func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

var _ QueryRepository = (*Repository)(nil)

func (r *Repository) FindById(ctx context.Context, id int) (*db.Order, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pborder.ApiResponseOrder, error) {
		return r.client.FindById(callCtx, &pborder.FindByIdOrderRequest{Id: int32(id)})
	})
	if err != nil {
		return nil, order_errors.ErrFindById.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, order_errors.ErrFindById
	}

	return mapOrder(resp.Data), nil
}

// mapOrder converts a proto order into the sqlc row type.
func mapOrder(o *pborder.OrderResponse) *db.Order {
	if o == nil {
		return nil
	}
	return &db.Order{
		OrderID:    o.Id,
		MerchantID: o.MerchantId,
		CashierID:  o.CashierId,
		TotalPrice: int64(o.TotalPrice),
		CreatedAt:  convert.Timestamp(o.CreatedAt),
		UpdatedAt:  convert.Timestamp(o.UpdatedAt),
	}
}
