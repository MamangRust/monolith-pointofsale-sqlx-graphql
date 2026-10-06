// Code migrated from sqlc-generated sources; domain types kept for repository layer.
// source: cashiers.sql

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type CreateCashierParams struct {
	MerchantID int32  `db:"merchant_id" json:"merchant_id"`
	UserID     int32  `db:"user_id" json:"user_id"`
	Name       string `db:"name" json:"name"`
}

type GetCashiersParams struct {
	Column1 string `db:"column1" json:"column_1"`
	Limit   int32  `db:"limit" json:"limit"`
	Offset  int32  `db:"offset" json:"offset"`
}

type GetCashiersRow struct {
	CashierID  int32            `db:"cashier_id" json:"cashier_id"`
	MerchantID int32            `db:"merchant_id" json:"merchant_id"`
	UserID     int32            `db:"user_id" json:"user_id"`
	Name       string           `db:"name" json:"name"`
	CreatedAt  pgtype.Timestamp `db:"created_at" json:"created_at"`
	UpdatedAt  pgtype.Timestamp `db:"updated_at" json:"updated_at"`
	DeletedAt  pgtype.Timestamp `db:"deleted_at" json:"deleted_at"`
	TotalCount int64            `db:"total_count" json:"total_count"`
}

type GetCashiersActiveParams struct {
	Column1 string `db:"column1" json:"column_1"`
	Limit   int32  `db:"limit" json:"limit"`
	Offset  int32  `db:"offset" json:"offset"`
}

type GetCashiersActiveRow struct {
	CashierID  int32            `db:"cashier_id" json:"cashier_id"`
	MerchantID int32            `db:"merchant_id" json:"merchant_id"`
	UserID     int32            `db:"user_id" json:"user_id"`
	Name       string           `db:"name" json:"name"`
	CreatedAt  pgtype.Timestamp `db:"created_at" json:"created_at"`
	UpdatedAt  pgtype.Timestamp `db:"updated_at" json:"updated_at"`
	DeletedAt  pgtype.Timestamp `db:"deleted_at" json:"deleted_at"`
	TotalCount int64            `db:"total_count" json:"total_count"`
}

type GetCashiersByMerchantParams struct {
	MerchantID int32  `db:"merchant_id" json:"merchant_id"`
	Column2    string `db:"column2" json:"column_2"`
	Limit      int32  `db:"limit" json:"limit"`
	Offset     int32  `db:"offset" json:"offset"`
}

type GetCashiersByMerchantRow struct {
	CashierID  int32            `db:"cashier_id" json:"cashier_id"`
	MerchantID int32            `db:"merchant_id" json:"merchant_id"`
	UserID     int32            `db:"user_id" json:"user_id"`
	Name       string           `db:"name" json:"name"`
	CreatedAt  pgtype.Timestamp `db:"created_at" json:"created_at"`
	UpdatedAt  pgtype.Timestamp `db:"updated_at" json:"updated_at"`
	DeletedAt  pgtype.Timestamp `db:"deleted_at" json:"deleted_at"`
	TotalCount int64            `db:"total_count" json:"total_count"`
}

type GetCashiersTrashedParams struct {
	Column1 string `db:"column1" json:"column_1"`
	Limit   int32  `db:"limit" json:"limit"`
	Offset  int32  `db:"offset" json:"offset"`
}

type GetCashiersTrashedRow struct {
	CashierID  int32            `db:"cashier_id" json:"cashier_id"`
	MerchantID int32            `db:"merchant_id" json:"merchant_id"`
	UserID     int32            `db:"user_id" json:"user_id"`
	Name       string           `db:"name" json:"name"`
	CreatedAt  pgtype.Timestamp `db:"created_at" json:"created_at"`
	UpdatedAt  pgtype.Timestamp `db:"updated_at" json:"updated_at"`
	DeletedAt  pgtype.Timestamp `db:"deleted_at" json:"deleted_at"`
	TotalCount int64            `db:"total_count" json:"total_count"`
}

type GetMonthlyCashierRow struct {
	CashierID   int32   `db:"cashier_id" json:"cashier_id"`
	CashierName string  `db:"cashier_name" json:"cashier_name"`
	Month       string  `db:"month" json:"month"`
	OrderCount  int64   `db:"order_count" json:"order_count"`
	TotalSales  float64 `db:"total_sales" json:"total_sales"`
}

type GetMonthlyCashierByCashierIdParams struct {
	Column1   time.Time `db:"column1" json:"column_1"`
	CashierID int32     `db:"cashier_id" json:"cashier_id"`
}

type GetMonthlyCashierByCashierIdRow struct {
	CashierID   int32  `db:"cashier_id" json:"cashier_id"`
	CashierName string `db:"cashier_name" json:"cashier_name"`
	Month       string `db:"month" json:"month"`
	OrderCount  int64  `db:"order_count" json:"order_count"`
	TotalSales  int64  `db:"total_sales" json:"total_sales"`
}

type GetMonthlyCashierByMerchantParams struct {
	Column1    time.Time `db:"column1" json:"column_1"`
	MerchantID int32     `db:"merchant_id" json:"merchant_id"`
}

type GetMonthlyCashierByMerchantRow struct {
	CashierID   int32  `db:"cashier_id" json:"cashier_id"`
	CashierName string `db:"cashier_name" json:"cashier_name"`
	Month       string `db:"month" json:"month"`
	OrderCount  int64  `db:"order_count" json:"order_count"`
	TotalSales  int64  `db:"total_sales" json:"total_sales"`
}

type GetMonthlyTotalSalesByIdParams struct {
	Extract     pgtype.Date      `db:"extract" json:"extract"`
	CreatedAt   pgtype.Timestamp `db:"created_at" json:"created_at"`
	CreatedAt_2 pgtype.Timestamp `db:"created_at_2" json:"created_at_2"`
	CreatedAt_3 pgtype.Timestamp `db:"created_at_3" json:"created_at_3"`
	CashierID   int32            `db:"cashier_id" json:"cashier_id"`
}

type GetMonthlyTotalSalesByIdRow struct {
	Year       string `db:"year" json:"year"`
	Month      string `db:"month" json:"month"`
	TotalSales int32  `db:"total_sales" json:"total_sales"`
}

type GetMonthlyTotalSalesByMerchantParams struct {
	Extract     pgtype.Date      `db:"extract" json:"extract"`
	CreatedAt   pgtype.Timestamp `db:"created_at" json:"created_at"`
	CreatedAt_2 pgtype.Timestamp `db:"created_at_2" json:"created_at_2"`
	CreatedAt_3 pgtype.Timestamp `db:"created_at_3" json:"created_at_3"`
	MerchantID  int32            `db:"merchant_id" json:"merchant_id"`
}

type GetMonthlyTotalSalesByMerchantRow struct {
	Year       string `db:"year" json:"year"`
	Month      string `db:"month" json:"month"`
	TotalSales int32  `db:"total_sales" json:"total_sales"`
}

type GetMonthlyTotalSalesCashierParams struct {
	Extract     pgtype.Date      `db:"extract" json:"extract"`
	CreatedAt   pgtype.Timestamp `db:"created_at" json:"created_at"`
	CreatedAt_2 pgtype.Timestamp `db:"created_at_2" json:"created_at_2"`
	CreatedAt_3 pgtype.Timestamp `db:"created_at_3" json:"created_at_3"`
}

type GetMonthlyTotalSalesCashierRow struct {
	Year       string `db:"year" json:"year"`
	Month      string `db:"month" json:"month"`
	TotalSales int32  `db:"total_sales" json:"total_sales"`
}

type GetYearlyCashierRow struct {
	Year        string `db:"year" json:"year"`
	CashierID   int32  `db:"cashier_id" json:"cashier_id"`
	CashierName string `db:"cashier_name" json:"cashier_name"`
	OrderCount  int64  `db:"order_count" json:"order_count"`
	TotalSales  int64  `db:"total_sales" json:"total_sales"`
}

type GetYearlyCashierByCashierIdParams struct {
	Column1   time.Time `db:"column1" json:"column_1"`
	CashierID int32     `db:"cashier_id" json:"cashier_id"`
}

type GetYearlyCashierByCashierIdRow struct {
	Year        string `db:"year" json:"year"`
	CashierID   int32  `db:"cashier_id" json:"cashier_id"`
	CashierName string `db:"cashier_name" json:"cashier_name"`
	OrderCount  int64  `db:"order_count" json:"order_count"`
	TotalSales  int64  `db:"total_sales" json:"total_sales"`
}

type GetYearlyCashierByMerchantParams struct {
	Column1    time.Time `db:"column1" json:"column_1"`
	MerchantID int32     `db:"merchant_id" json:"merchant_id"`
}

type GetYearlyCashierByMerchantRow struct {
	Year        string `db:"year" json:"year"`
	CashierID   int32  `db:"cashier_id" json:"cashier_id"`
	CashierName string `db:"cashier_name" json:"cashier_name"`
	OrderCount  int64  `db:"order_count" json:"order_count"`
	TotalSales  int64  `db:"total_sales" json:"total_sales"`
}

type GetYearlyTotalSalesByIdParams struct {
	Column1   int32 `db:"column1" json:"column_1"`
	CashierID int32 `db:"cashier_id" json:"cashier_id"`
}

type GetYearlyTotalSalesByIdRow struct {
	Year       string `db:"year" json:"year"`
	TotalSales int32  `db:"total_sales" json:"total_sales"`
}

type GetYearlyTotalSalesByMerchantParams struct {
	Column1    int32 `db:"column1" json:"column_1"`
	MerchantID int32 `db:"merchant_id" json:"merchant_id"`
}

type GetYearlyTotalSalesByMerchantRow struct {
	Year       string `db:"year" json:"year"`
	TotalSales int32  `db:"total_sales" json:"total_sales"`
}

type GetYearlyTotalSalesCashierRow struct {
	Year       string `db:"year" json:"year"`
	TotalSales int32  `db:"total_sales" json:"total_sales"`
}

type UpdateCashierParams struct {
	CashierID int32  `db:"cashier_id" json:"cashier_id"`
	Name      string `db:"name" json:"name"`
}
