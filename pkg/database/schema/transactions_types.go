// Code migrated from sqlc-generated sources; domain types kept for repository layer.
// source: transactions.sql

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type CreateTransactionParams struct {
	MerchantID    int32  `db:"merchant_id" json:"merchant_id"`
	PaymentMethod string `db:"payment_method" json:"payment_method"`
	Amount        int32  `db:"amount" json:"amount"`
	ChangeAmount  *int32 `db:"change_amount" json:"change_amount"`
	PaymentStatus string `db:"payment_status" json:"payment_status"`
	OrderID       int32  `db:"order_id" json:"order_id"`
}

type GetMonthlyAmountTransactionFailedParams struct {
	Column1 time.Time `db:"column1" json:"column_1"`
	Column2 time.Time `db:"column2" json:"column_2"`
	Column3 time.Time `db:"column3" json:"column_3"`
	Column4 time.Time `db:"column4" json:"column_4"`
}

type GetMonthlyAmountTransactionFailedRow struct {
	Year        string `db:"year" json:"year"`
	Month       string `db:"month" json:"month"`
	TotalFailed int64  `db:"total_failed" json:"total_failed"`
	TotalAmount int32  `db:"total_amount" json:"total_amount"`
}

type GetMonthlyAmountTransactionFailedByMerchantParams struct {
	Column1    time.Time `db:"column1" json:"column_1"`
	Column2    time.Time `db:"column2" json:"column_2"`
	Column3    time.Time `db:"column3" json:"column_3"`
	Column4    time.Time `db:"column4" json:"column_4"`
	MerchantID int32     `db:"merchant_id" json:"merchant_id"`
}

type GetMonthlyAmountTransactionFailedByMerchantRow struct {
	Year        string `db:"year" json:"year"`
	Month       string `db:"month" json:"month"`
	TotalFailed int64  `db:"total_failed" json:"total_failed"`
	TotalAmount int32  `db:"total_amount" json:"total_amount"`
}

type GetMonthlyAmountTransactionSuccessParams struct {
	Column1 time.Time `db:"column1" json:"column_1"`
	Column2 time.Time `db:"column2" json:"column_2"`
	Column3 time.Time `db:"column3" json:"column_3"`
	Column4 time.Time `db:"column4" json:"column_4"`
}

type GetMonthlyAmountTransactionSuccessRow struct {
	Year         string `db:"year" json:"year"`
	Month        string `db:"month" json:"month"`
	TotalSuccess int64  `db:"total_success" json:"total_success"`
	TotalAmount  int32  `db:"total_amount" json:"total_amount"`
}

type GetMonthlyAmountTransactionSuccessByMerchantParams struct {
	Column1    time.Time `db:"column1" json:"column_1"`
	Column2    time.Time `db:"column2" json:"column_2"`
	Column3    time.Time `db:"column3" json:"column_3"`
	Column4    time.Time `db:"column4" json:"column_4"`
	MerchantID int32     `db:"merchant_id" json:"merchant_id"`
}

type GetMonthlyAmountTransactionSuccessByMerchantRow struct {
	Year         string `db:"year" json:"year"`
	Month        string `db:"month" json:"month"`
	TotalSuccess int64  `db:"total_success" json:"total_success"`
	TotalAmount  int32  `db:"total_amount" json:"total_amount"`
}

type GetMonthlyTransactionMethodsByMerchantFailedParams struct {
	Column1    time.Time `db:"column1" json:"column_1"`
	Column2    time.Time `db:"column2" json:"column_2"`
	Column3    time.Time `db:"column3" json:"column_3"`
	Column4    time.Time `db:"column4" json:"column_4"`
	MerchantID int32     `db:"merchant_id" json:"merchant_id"`
}

type GetMonthlyTransactionMethodsByMerchantFailedRow struct {
	Month             string  `db:"month" json:"month"`
	PaymentMethod     string  `db:"payment_method" json:"payment_method"`
	TotalTransactions int64   `db:"total_transactions" json:"total_transactions"`
	TotalAmount       float64 `db:"total_amount" json:"total_amount"`
}

type GetMonthlyTransactionMethodsByMerchantSuccessParams struct {
	Column1    time.Time `db:"column1" json:"column_1"`
	Column2    time.Time `db:"column2" json:"column_2"`
	Column3    time.Time `db:"column3" json:"column_3"`
	Column4    time.Time `db:"column4" json:"column_4"`
	MerchantID int32     `db:"merchant_id" json:"merchant_id"`
}

type GetMonthlyTransactionMethodsByMerchantSuccessRow struct {
	Month             string  `db:"month" json:"month"`
	PaymentMethod     string  `db:"payment_method" json:"payment_method"`
	TotalTransactions int64   `db:"total_transactions" json:"total_transactions"`
	TotalAmount       float64 `db:"total_amount" json:"total_amount"`
}

type GetMonthlyTransactionMethodsFailedParams struct {
	Column1 time.Time `db:"column1" json:"column_1"`
	Column2 time.Time `db:"column2" json:"column_2"`
	Column3 time.Time `db:"column3" json:"column_3"`
	Column4 time.Time `db:"column4" json:"column_4"`
}

type GetMonthlyTransactionMethodsFailedRow struct {
	Month             string  `db:"month" json:"month"`
	PaymentMethod     string  `db:"payment_method" json:"payment_method"`
	TotalTransactions int64   `db:"total_transactions" json:"total_transactions"`
	TotalAmount       float64 `db:"total_amount" json:"total_amount"`
}

type GetMonthlyTransactionMethodsSuccessParams struct {
	Column1 time.Time `db:"column1" json:"column_1"`
	Column2 time.Time `db:"column2" json:"column_2"`
	Column3 time.Time `db:"column3" json:"column_3"`
	Column4 time.Time `db:"column4" json:"column_4"`
}

type GetMonthlyTransactionMethodsSuccessRow struct {
	Month             string  `db:"month" json:"month"`
	PaymentMethod     string  `db:"payment_method" json:"payment_method"`
	TotalTransactions int64   `db:"total_transactions" json:"total_transactions"`
	TotalAmount       float64 `db:"total_amount" json:"total_amount"`
}

type GetTransactionByMerchantParams struct {
	Column1 string `db:"column1" json:"column_1"`
	Column2 int32  `db:"column2" json:"column_2"`
	Limit   int32  `db:"limit" json:"limit"`
	Offset  int32  `db:"offset" json:"offset"`
}

type GetTransactionByMerchantRow struct {
	TransactionID int32            `db:"transaction_id" json:"transaction_id"`
	OrderID       int32            `db:"order_id" json:"order_id"`
	MerchantID    int32            `db:"merchant_id" json:"merchant_id"`
	PaymentMethod string           `db:"payment_method" json:"payment_method"`
	Amount        int32            `db:"amount" json:"amount"`
	ChangeAmount  *int32           `db:"change_amount" json:"change_amount"`
	PaymentStatus string           `db:"payment_status" json:"payment_status"`
	CreatedAt     pgtype.Timestamp `db:"created_at" json:"created_at"`
	UpdatedAt     pgtype.Timestamp `db:"updated_at" json:"updated_at"`
	DeletedAt     pgtype.Timestamp `db:"deleted_at" json:"deleted_at"`
	TotalCount    int64            `db:"total_count" json:"total_count"`
}

type GetTransactionsParams struct {
	Column1 string `db:"column1" json:"column_1"`
	Limit   int32  `db:"limit" json:"limit"`
	Offset  int32  `db:"offset" json:"offset"`
}

type GetTransactionsRow struct {
	TransactionID int32            `db:"transaction_id" json:"transaction_id"`
	OrderID       int32            `db:"order_id" json:"order_id"`
	MerchantID    int32            `db:"merchant_id" json:"merchant_id"`
	PaymentMethod string           `db:"payment_method" json:"payment_method"`
	Amount        int32            `db:"amount" json:"amount"`
	ChangeAmount  *int32           `db:"change_amount" json:"change_amount"`
	PaymentStatus string           `db:"payment_status" json:"payment_status"`
	CreatedAt     pgtype.Timestamp `db:"created_at" json:"created_at"`
	UpdatedAt     pgtype.Timestamp `db:"updated_at" json:"updated_at"`
	DeletedAt     pgtype.Timestamp `db:"deleted_at" json:"deleted_at"`
	TotalCount    int64            `db:"total_count" json:"total_count"`
}

type GetTransactionsActiveParams struct {
	Column1 string `db:"column1" json:"column_1"`
	Limit   int32  `db:"limit" json:"limit"`
	Offset  int32  `db:"offset" json:"offset"`
}

type GetTransactionsActiveRow struct {
	TransactionID int32            `db:"transaction_id" json:"transaction_id"`
	OrderID       int32            `db:"order_id" json:"order_id"`
	MerchantID    int32            `db:"merchant_id" json:"merchant_id"`
	PaymentMethod string           `db:"payment_method" json:"payment_method"`
	Amount        int32            `db:"amount" json:"amount"`
	ChangeAmount  *int32           `db:"change_amount" json:"change_amount"`
	PaymentStatus string           `db:"payment_status" json:"payment_status"`
	CreatedAt     pgtype.Timestamp `db:"created_at" json:"created_at"`
	UpdatedAt     pgtype.Timestamp `db:"updated_at" json:"updated_at"`
	DeletedAt     pgtype.Timestamp `db:"deleted_at" json:"deleted_at"`
	TotalCount    int64            `db:"total_count" json:"total_count"`
}

type GetTransactionsTrashedParams struct {
	Column1 string `db:"column1" json:"column_1"`
	Limit   int32  `db:"limit" json:"limit"`
	Offset  int32  `db:"offset" json:"offset"`
}

type GetTransactionsTrashedRow struct {
	TransactionID int32            `db:"transaction_id" json:"transaction_id"`
	OrderID       int32            `db:"order_id" json:"order_id"`
	MerchantID    int32            `db:"merchant_id" json:"merchant_id"`
	PaymentMethod string           `db:"payment_method" json:"payment_method"`
	Amount        int32            `db:"amount" json:"amount"`
	ChangeAmount  *int32           `db:"change_amount" json:"change_amount"`
	PaymentStatus string           `db:"payment_status" json:"payment_status"`
	CreatedAt     pgtype.Timestamp `db:"created_at" json:"created_at"`
	UpdatedAt     pgtype.Timestamp `db:"updated_at" json:"updated_at"`
	DeletedAt     pgtype.Timestamp `db:"deleted_at" json:"deleted_at"`
	TotalCount    int64            `db:"total_count" json:"total_count"`
}

type GetYearlyAmountTransactionFailedRow struct {
	Year        string `db:"year" json:"year"`
	TotalFailed int32  `db:"total_failed" json:"total_failed"`
	TotalAmount int32  `db:"total_amount" json:"total_amount"`
}

type GetYearlyAmountTransactionFailedByMerchantParams struct {
	Column1    int32 `db:"column1" json:"column_1"`
	MerchantID int32 `db:"merchant_id" json:"merchant_id"`
}

type GetYearlyAmountTransactionFailedByMerchantRow struct {
	Year        string `db:"year" json:"year"`
	TotalFailed int32  `db:"total_failed" json:"total_failed"`
	TotalAmount int32  `db:"total_amount" json:"total_amount"`
}

type GetYearlyAmountTransactionSuccessRow struct {
	Year         string `db:"year" json:"year"`
	TotalSuccess int32  `db:"total_success" json:"total_success"`
	TotalAmount  int32  `db:"total_amount" json:"total_amount"`
}

type GetYearlyAmountTransactionSuccessByMerchantParams struct {
	Column1    int32 `db:"column1" json:"column_1"`
	MerchantID int32 `db:"merchant_id" json:"merchant_id"`
}

type GetYearlyAmountTransactionSuccessByMerchantRow struct {
	Year         string `db:"year" json:"year"`
	TotalSuccess int32  `db:"total_success" json:"total_success"`
	TotalAmount  int32  `db:"total_amount" json:"total_amount"`
}

type GetYearlyTransactionMethodsByMerchantFailedParams struct {
	Column1    time.Time `db:"column1" json:"column_1"`
	MerchantID int32     `db:"merchant_id" json:"merchant_id"`
}

type GetYearlyTransactionMethodsByMerchantFailedRow struct {
	Year              string  `db:"year" json:"year"`
	PaymentMethod     string  `db:"payment_method" json:"payment_method"`
	TotalTransactions int64   `db:"total_transactions" json:"total_transactions"`
	TotalAmount       float64 `db:"total_amount" json:"total_amount"`
}

type GetYearlyTransactionMethodsByMerchantSuccessParams struct {
	Column1    time.Time `db:"column1" json:"column_1"`
	MerchantID int32     `db:"merchant_id" json:"merchant_id"`
}

type GetYearlyTransactionMethodsByMerchantSuccessRow struct {
	Year              string  `db:"year" json:"year"`
	PaymentMethod     string  `db:"payment_method" json:"payment_method"`
	TotalTransactions int64   `db:"total_transactions" json:"total_transactions"`
	TotalAmount       float64 `db:"total_amount" json:"total_amount"`
}

type GetYearlyTransactionMethodsFailedRow struct {
	Year              string  `db:"year" json:"year"`
	PaymentMethod     string  `db:"payment_method" json:"payment_method"`
	TotalTransactions int64   `db:"total_transactions" json:"total_transactions"`
	TotalAmount       float64 `db:"total_amount" json:"total_amount"`
}

type GetYearlyTransactionMethodsSuccessRow struct {
	Year              string  `db:"year" json:"year"`
	PaymentMethod     string  `db:"payment_method" json:"payment_method"`
	TotalTransactions int64   `db:"total_transactions" json:"total_transactions"`
	TotalAmount       float64 `db:"total_amount" json:"total_amount"`
}

type UpdateTransactionParams struct {
	TransactionID int32  `db:"transaction_id" json:"transaction_id"`
	MerchantID    int32  `db:"merchant_id" json:"merchant_id"`
	PaymentMethod string `db:"payment_method" json:"payment_method"`
	Amount        int32  `db:"amount" json:"amount"`
	ChangeAmount  *int32 `db:"change_amount" json:"change_amount"`
	PaymentStatus string `db:"payment_status" json:"payment_status"`
	OrderID       int32  `db:"order_id" json:"order_id"`
}
