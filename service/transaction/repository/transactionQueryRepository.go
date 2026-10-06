package repository

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/transaction_errors"
	"github.com/jmoiron/sqlx"
)

const getTransactions = `SELECT *, COUNT(*) OVER () AS total_count
FROM transactions
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR payment_method ILIKE '%' || $1 || '%'
        OR payment_status ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3`

const getTransactionsActive = `SELECT *, COUNT(*) OVER () AS total_count
FROM transactions
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR payment_method ILIKE '%' || $1 || '%'
        OR payment_status ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3`

const getTransactionsTrashed = `SELECT *, COUNT(*) OVER () AS total_count
FROM transactions
WHERE
    deleted_at IS NOT NULL
    AND (
        $1::TEXT IS NULL
        OR payment_method ILIKE '%' || $1 || '%'
        OR payment_status ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3`

const getTransactionByMerchant = `SELECT *, COUNT(*) OVER () AS total_count
FROM transactions
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR payment_method ILIKE '%' || $1 || '%'
        OR payment_status ILIKE '%' || $1 || '%'
    )
    AND (
        $2::INT IS NULL
        OR merchant_id = $2
    )
ORDER BY created_at DESC
LIMIT $3
OFFSET
    $4`

const getTransactionByID = `SELECT *
FROM transactions
WHERE
    transaction_id = $1
    AND deleted_at IS NULL`

const getTransactionByOrderID = `SELECT *
FROM transactions
WHERE
    order_id = $1
    AND deleted_at IS NULL`

type transactionQueryRepository struct {
	db *sqlx.DB
}

func NewTransactionQueryRepository(db *sqlx.DB) *transactionQueryRepository {
	return &transactionQueryRepository{
		db: db,
	}
}

func (r *transactionQueryRepository) FindAllTransactions(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var transactions []*db.GetTransactionsRow
	err := r.db.SelectContext(ctx, &transactions, getTransactions,
		req.Search,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, transaction_errors.ErrFindAllTransactions
	}

	var totalCount int
	if len(transactions) > 0 {
		totalCount = int(transactions[0].TotalCount)
	} else {
		totalCount = 0
	}

	return transactions, &totalCount, nil
}

func (r *transactionQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsActiveRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var transactions []*db.GetTransactionsActiveRow
	err := r.db.SelectContext(ctx, &transactions, getTransactionsActive,
		req.Search,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, transaction_errors.ErrFindByActive
	}

	var totalCount int
	if len(transactions) > 0 {
		totalCount = int(transactions[0].TotalCount)
	} else {
		totalCount = 0
	}

	return transactions, &totalCount, nil
}

func (r *transactionQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsTrashedRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var transactions []*db.GetTransactionsTrashedRow
	err := r.db.SelectContext(ctx, &transactions, getTransactionsTrashed,
		req.Search,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, transaction_errors.ErrFindByTrashed
	}

	var totalCount int
	if len(transactions) > 0 {
		totalCount = int(transactions[0].TotalCount)
	} else {
		totalCount = 0
	}

	return transactions, &totalCount, nil
}

func (r *transactionQueryRepository) FindByMerchant(
	ctx context.Context,
	req *requests.FindAllTransactionByMerchant,
) ([]*db.GetTransactionByMerchantRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var transactions []*db.GetTransactionByMerchantRow
	err := r.db.SelectContext(ctx, &transactions, getTransactionByMerchant,
		req.Search,
		int32(req.MerchantID),
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, transaction_errors.ErrFindByMerchant
	}

	var totalCount int
	if len(transactions) > 0 {
		totalCount = int(transactions[0].TotalCount)
	} else {
		totalCount = 0
	}

	return transactions, &totalCount, nil
}

func (r *transactionQueryRepository) FindById(ctx context.Context, transaction_id int) (*db.Transaction, error) {
	var transaction db.Transaction
	err := r.db.GetContext(ctx, &transaction, getTransactionByID, int32(transaction_id))
	if err != nil {
		return nil, transaction_errors.ErrFindById
	}

	return &transaction, nil
}

func (r *transactionQueryRepository) FindByOrderId(ctx context.Context, order_id int) (*db.Transaction, error) {
	var transaction db.Transaction
	err := r.db.GetContext(ctx, &transaction, getTransactionByOrderID, int32(order_id))
	if err != nil {
		return nil, transaction_errors.ErrFindByOrderId
	}

	return &transaction, nil
}
