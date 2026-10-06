// Package merchant adapts the Merchant service gRPC API into the sqlc row types
// used by the consumers of merchants (cashier, order, transaction, product).
package merchant

import (
	"context"

	pbmerchant "github.com/MamangRust/monolith-pointofsale-pb/merchant"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter/convert"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/resilience"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/merchant_errors"
)

// Repository reads merchants from the Merchant service over gRPC.
type Repository interface {
	FindById(ctx context.Context, id int) (*db.Merchant, error)
}

type repository struct {
	client pbmerchant.MerchantQueryServiceClient
	guard  *resilience.DependencyGuard
}

// New builds a merchant adapter. Passing zero options leaves the guard nil,
// which makes DependencyGuard.Call a plain passthrough.
func New(client pbmerchant.MerchantQueryServiceClient, opts ...adapter.GuardOption) Repository {
	r := &repository{client: client}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// SetGuard lets adapter.WithDependencyGuard attach a guard after construction.
func (r *repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

var _ Repository = (*repository)(nil)

func (r *repository) FindById(ctx context.Context, id int) (*db.Merchant, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbmerchant.ApiResponseMerchant, error) {
		return r.client.FindById(callCtx, &pbmerchant.FindByIdMerchantRequest{Id: int32(id)})
	})
	if err != nil {
		return nil, merchant_errors.ErrFindById.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, merchant_errors.ErrFindById
	}

	m := resp.Data
	return &db.Merchant{
		MerchantID:   m.Id,
		UserID:       m.UserId,
		Name:         m.Name,
		Description:  convert.NullableString(m.Description),
		Address:      convert.NullableString(m.Address),
		ContactEmail: convert.NullableString(m.ContactEmail),
		ContactPhone: convert.NullableString(m.ContactPhone),
		Status:       m.Status,
		CreatedAt:    convert.Timestamp(m.CreatedAt),
		UpdatedAt:    convert.Timestamp(m.UpdatedAt),
	}, nil
}
