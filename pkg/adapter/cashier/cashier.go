// Package cashier adapts the Cashier service gRPC API into the sqlc row types
// used by the consumer repositories.
package cashier

import (
	"context"

	pbcashier "github.com/MamangRust/monolith-pointofsale-pb/cashier"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter/convert"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/resilience"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/cashier_errors"
)

// QueryRepository reads a single cashier from the Cashier service over gRPC.
type QueryRepository interface {
	FindById(ctx context.Context, id int) (*db.Cashier, error)
}

// Repository implements QueryRepository over a gRPC client.
type Repository struct {
	client pbcashier.CashierQueryServiceClient
	guard  *resilience.DependencyGuard
}

// New builds a cashier adapter. Passing zero options leaves the guard nil,
// which makes DependencyGuard.Call a plain passthrough.
func New(client pbcashier.CashierQueryServiceClient, opts ...adapter.GuardOption) *Repository {
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

func (r *Repository) FindById(ctx context.Context, id int) (*db.Cashier, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbcashier.ApiResponseCashier, error) {
		return r.client.FindById(callCtx, &pbcashier.FindByIdCashierRequest{Id: int32(id)})
	})
	if err != nil {
		return nil, cashier_errors.ErrFindCashierById.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, cashier_errors.ErrFindCashierById
	}

	return mapCashier(resp.Data), nil
}

// mapCashier converts a proto cashier into the sqlc row type. The proto response
// does not carry user_id, so that column keeps its zero value.
func mapCashier(c *pbcashier.CashierResponse) *db.Cashier {
	if c == nil {
		return nil
	}
	return &db.Cashier{
		CashierID:  c.Id,
		MerchantID: c.MerchantId,
		Name:       c.Name,
		CreatedAt:  convert.Timestamp(c.CreatedAt),
		UpdatedAt:  convert.Timestamp(c.UpdatedAt),
	}
}
