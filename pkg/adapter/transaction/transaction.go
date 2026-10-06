// Package transaction adapts the Transaction service gRPC API into the sqlc row
// types used by the consumer repositories.
package transaction

import (
	"context"

	pbtransaction "github.com/MamangRust/monolith-pointofsale-pb/transaction"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter"
	"github.com/MamangRust/monolith-pointofsale-pkg/adapter/convert"
	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-pkg/resilience"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/transaction_errors"
)

// BulkRepository enumerates transactions page by page. It is a narrow interface
// for offline consumers (the stats backfill) that need the whole collection.
type BulkRepository interface {
	FindAll(ctx context.Context, page, pageSize int) ([]*db.Transaction, int, error)
}

// Repository implements BulkRepository over a gRPC client.
type Repository struct {
	client pbtransaction.TransactionQueryServiceClient
	guard  *resilience.DependencyGuard
}

// New builds a transaction adapter. Passing zero options leaves the guard nil,
// which makes DependencyGuard.Call a plain passthrough.
func New(client pbtransaction.TransactionQueryServiceClient, opts ...adapter.GuardOption) *Repository {
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

var _ BulkRepository = (*Repository)(nil)

// FindAll returns one page of transactions plus the total number of rows, so
// callers can walk the collection page by page.
func (r *Repository) FindAll(ctx context.Context, page, pageSize int) ([]*db.Transaction, int, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbtransaction.ApiResponsePaginationTransaction, error) {
		return r.client.FindAll(callCtx, &pbtransaction.FindAllTransactionRequest{
			Page:     int32(page),
			PageSize: int32(pageSize),
		})
	})
	if err != nil {
		return nil, 0, transaction_errors.ErrFindAllTransactions.WithInternal(err)
	}
	if resp == nil {
		return nil, 0, transaction_errors.ErrFindAllTransactions
	}

	transactions := make([]*db.Transaction, 0, len(resp.Data))
	for _, t := range resp.Data {
		if t == nil {
			continue
		}
		transactions = append(transactions, mapTransaction(t))
	}

	total := 0
	if resp.Pagination != nil {
		total = int(resp.Pagination.TotalRecords)
	}
	return transactions, total, nil
}

// mapTransaction converts a proto transaction into the sqlc row type.
func mapTransaction(t *pbtransaction.TransactionResponse) *db.Transaction {
	if t == nil {
		return nil
	}
	return &db.Transaction{
		TransactionID: t.Id,
		OrderID:       t.OrderId,
		MerchantID:    t.MerchantId,
		PaymentMethod: t.PaymentMethod,
		Amount:        t.Amount,
		ChangeAmount:  convert.NullableInt32(t.ChangeAmount),
		PaymentStatus: t.PaymentStatus,
		CreatedAt:     convert.Timestamp(t.CreatedAt),
		UpdatedAt:     convert.Timestamp(t.UpdatedAt),
	}
}
