// Package product adapts the Product service gRPC API into the sqlc row types
// used by the consumer repositories. Query and command are exposed as separate
// interfaces so consumers can depend on the narrower one.
package product

import (
	"context"

	pbproduct "github.com/MamangRust/monolith-pointofsale-pb/product"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter/convert"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/resilience"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/product_errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// QueryRepository reads products from the Product service over gRPC.
type QueryRepository interface {
	FindById(ctx context.Context, productID int) (*db.Product, error)
}

// CommandRepository is the contract consumers depend on for product stock
// writes. Create/Update/Trash live entirely inside the owning service and are
// not exposed here.
type CommandRepository interface {
	AdjustProductStock(ctx context.Context, productID int, delta int) (*db.Product, error)
	DecrementProductCountStock(ctx context.Context, productID int, quantity int) (*db.Product, error)
	IncrementProductCountStock(ctx context.Context, productID int, quantity int) (*db.Product, error)
}

// Repository implements QueryRepository and CommandRepository over the
// generated product query/command clients, so consumers can share one instance
// across both roles.
type Repository struct {
	query   pbproduct.ProductQueryServiceClient
	command pbproduct.ProductCommandServiceClient
	guard   *resilience.DependencyGuard
}

// New builds a product adapter with both halves wired. Passing zero options
// leaves the guard nil, which makes DependencyGuard.Call a plain passthrough.
func New(query pbproduct.ProductQueryServiceClient, command pbproduct.ProductCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{query: query, command: command}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewQueryAdapter builds a read-only adapter for consumers that never write.
func NewQueryAdapter(query pbproduct.ProductQueryServiceClient, opts ...adapter.GuardOption) QueryRepository {
	return New(query, nil, opts...)
}

// NewCommandAdapter builds a write-only adapter for consumers that never read.
func NewCommandAdapter(command pbproduct.ProductCommandServiceClient, opts ...adapter.GuardOption) CommandRepository {
	return New(nil, command, opts...)
}

// SetGuard lets adapter.WithDependencyGuard attach a guard after construction.
func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

var (
	_ QueryRepository   = (*Repository)(nil)
	_ CommandRepository = (*Repository)(nil)
)

func (r *Repository) FindById(ctx context.Context, productID int) (*db.Product, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbproduct.ApiResponseProduct, error) {
		return r.query.FindById(callCtx, &pbproduct.FindByIdProductRequest{Id: int32(productID)})
	})
	if err != nil {
		return nil, product_errors.ErrFindById.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, product_errors.ErrFindById
	}

	return mapResponse(resp.Data), nil
}

// AdjustProductStock applies a signed delta to a product's stock through the
// owning service, which enforces the atomic `count_in_stock + delta >= 0` guard.
func (r *Repository) AdjustProductStock(ctx context.Context, productID int, delta int) (*db.Product, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbproduct.ApiResponseProduct, error) {
		return r.command.AdjustProductStock(callCtx, &pbproduct.AdjustProductStockRequest{
			ProductId: int32(productID),
			Delta:     int32(delta),
		})
	})
	if err != nil {
		// The owning service reports "would go negative" as InvalidArgument
		// (HTTP 400); surface it as the insufficient-stock domain error so the
		// order flow can roll back rather than log an internal failure.
		if status.Code(err) == codes.InvalidArgument {
			return nil, product_errors.ErrInsufficientProductStock.WithInternal(err)
		}
		return nil, product_errors.ErrUpdateProductCountStock.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, product_errors.ErrUpdateProductCountStock
	}

	return mapResponse(resp.Data), nil
}

// DecrementProductCountStock reduces stock by quantity; it is the negative-delta
// spelling used by the order flow when reserving items for a new order.
func (r *Repository) DecrementProductCountStock(ctx context.Context, productID int, quantity int) (*db.Product, error) {
	return r.AdjustProductStock(ctx, productID, -quantity)
}

// IncrementProductCountStock restores stock by quantity; it is the positive-delta
// spelling used by the order flow to compensate a failed order.
func (r *Repository) IncrementProductCountStock(ctx context.Context, productID int, quantity int) (*db.Product, error) {
	return r.AdjustProductStock(ctx, productID, quantity)
}

// mapResponse converts a proto product into the sqlc row type.
func mapResponse(p *pbproduct.ProductResponse) *db.Product {
	if p == nil {
		return nil
	}
	return &db.Product{
		ProductID:    p.Id,
		MerchantID:   p.MerchantId,
		CategoryID:   p.CategoryId,
		Name:         p.Name,
		Description:  convert.NullableString(p.Description),
		Price:        p.Price,
		CountInStock: p.CountInStock,
		Brand:        convert.NullableString(p.Brand),
		Weight:       convert.NullableInt32(p.Weight),
		SlugProduct:  convert.NullableString(p.SlugProduct),
		ImageProduct: convert.NullableString(p.ImageProduct),
		Barcode:      convert.NullableString(p.Barcode),
		CreatedAt:    convert.Timestamp(p.CreatedAt),
		UpdatedAt:    convert.Timestamp(p.UpdatedAt),
	}
}
