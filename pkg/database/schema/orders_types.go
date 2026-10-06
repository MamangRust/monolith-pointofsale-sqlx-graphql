// Code migrated from sqlc-generated sources; domain types kept for repository layer.
// source: orders.sql

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type CreateOrderParams struct {
	MerchantID int32 `db:"merchant_id" json:"merchant_id"`
	CashierID  int32 `db:"cashier_id" json:"cashier_id"`
	TotalPrice int64 `db:"total_price" json:"total_price"`
}

type DeleteOrderPermanentlyRow struct {
	OrderID    int32            `db:"order_id" json:"order_id"`
	MerchantID int32            `db:"merchant_id" json:"merchant_id"`
	CashierID  int32            `db:"cashier_id" json:"cashier_id"`
	TotalPrice int64            `db:"total_price" json:"total_price"`
	CreatedAt  pgtype.Timestamp `db:"created_at" json:"created_at"`
	UpdatedAt  pgtype.Timestamp `db:"updated_at" json:"updated_at"`
	DeletedAt  pgtype.Timestamp `db:"deleted_at" json:"deleted_at"`
}

type GetMonthlyOrderRow struct {
	Month          string  `db:"month" json:"month"`
	OrderCount     int64   `db:"order_count" json:"order_count"`
	TotalRevenue   float64 `db:"total_revenue" json:"total_revenue"`
	TotalItemsSold int64   `db:"total_items_sold" json:"total_items_sold"`
}

type GetMonthlyOrderByMerchantParams struct {
	Column1    time.Time `db:"column1" json:"column_1"`
	MerchantID int32     `db:"merchant_id" json:"merchant_id"`
}

type GetMonthlyOrderByMerchantRow struct {
	Month          string  `db:"month" json:"month"`
	OrderCount     int64   `db:"order_count" json:"order_count"`
	TotalRevenue   float64 `db:"total_revenue" json:"total_revenue"`
	TotalItemsSold int64   `db:"total_items_sold" json:"total_items_sold"`
}

type GetMonthlyTotalRevenueParams struct {
	Extract     pgtype.Date      `db:"extract" json:"extract"`
	CreatedAt   pgtype.Timestamp `db:"created_at" json:"created_at"`
	CreatedAt_2 pgtype.Timestamp `db:"created_at_2" json:"created_at_2"`
	CreatedAt_3 pgtype.Timestamp `db:"created_at_3" json:"created_at_3"`
}

type GetMonthlyTotalRevenueRow struct {
	Year         string `db:"year" json:"year"`
	Month        string `db:"month" json:"month"`
	TotalRevenue int32  `db:"total_revenue" json:"total_revenue"`
}

type GetMonthlyTotalRevenueByIdParams struct {
	Extract     pgtype.Date      `db:"extract" json:"extract"`
	CreatedAt   pgtype.Timestamp `db:"created_at" json:"created_at"`
	CreatedAt_2 pgtype.Timestamp `db:"created_at_2" json:"created_at_2"`
	CreatedAt_3 pgtype.Timestamp `db:"created_at_3" json:"created_at_3"`
	OrderID     int32            `db:"order_id" json:"order_id"`
}

type GetMonthlyTotalRevenueByIdRow struct {
	Year         string `db:"year" json:"year"`
	Month        string `db:"month" json:"month"`
	TotalRevenue int32  `db:"total_revenue" json:"total_revenue"`
}

type GetMonthlyTotalRevenueByMerchantParams struct {
	Extract     pgtype.Date      `db:"extract" json:"extract"`
	CreatedAt   pgtype.Timestamp `db:"created_at" json:"created_at"`
	CreatedAt_2 pgtype.Timestamp `db:"created_at_2" json:"created_at_2"`
	CreatedAt_3 pgtype.Timestamp `db:"created_at_3" json:"created_at_3"`
	MerchantID  int32            `db:"merchant_id" json:"merchant_id"`
}

type GetMonthlyTotalRevenueByMerchantRow struct {
	Year         string `db:"year" json:"year"`
	Month        string `db:"month" json:"month"`
	TotalRevenue int32  `db:"total_revenue" json:"total_revenue"`
}

type GetOrdersParams struct {
	Column1 string `db:"column1" json:"column_1"`
	Limit   int32  `db:"limit" json:"limit"`
	Offset  int32  `db:"offset" json:"offset"`
}

type GetOrdersRow struct {
	OrderID    int32            `db:"order_id" json:"order_id"`
	MerchantID int32            `db:"merchant_id" json:"merchant_id"`
	CashierID  int32            `db:"cashier_id" json:"cashier_id"`
	TotalPrice int64            `db:"total_price" json:"total_price"`
	CreatedAt  pgtype.Timestamp `db:"created_at" json:"created_at"`
	UpdatedAt  pgtype.Timestamp `db:"updated_at" json:"updated_at"`
	DeletedAt  pgtype.Timestamp `db:"deleted_at" json:"deleted_at"`
	TotalCount int64            `db:"total_count" json:"total_count"`
}

type GetOrdersActiveParams struct {
	Column1 string `db:"column1" json:"column_1"`
	Limit   int32  `db:"limit" json:"limit"`
	Offset  int32  `db:"offset" json:"offset"`
}

type GetOrdersActiveRow struct {
	OrderID    int32            `db:"order_id" json:"order_id"`
	MerchantID int32            `db:"merchant_id" json:"merchant_id"`
	CashierID  int32            `db:"cashier_id" json:"cashier_id"`
	TotalPrice int64            `db:"total_price" json:"total_price"`
	CreatedAt  pgtype.Timestamp `db:"created_at" json:"created_at"`
	UpdatedAt  pgtype.Timestamp `db:"updated_at" json:"updated_at"`
	DeletedAt  pgtype.Timestamp `db:"deleted_at" json:"deleted_at"`
	TotalCount int64            `db:"total_count" json:"total_count"`
}

type GetOrdersByMerchantParams struct {
	Column1 string `db:"column1" json:"column_1"`
	Limit   int32  `db:"limit" json:"limit"`
	Offset  int32  `db:"offset" json:"offset"`
	Column4 int32  `db:"column4" json:"column_4"`
}

type GetOrdersByMerchantRow struct {
	OrderID    int32            `db:"order_id" json:"order_id"`
	MerchantID int32            `db:"merchant_id" json:"merchant_id"`
	CashierID  int32            `db:"cashier_id" json:"cashier_id"`
	TotalPrice int64            `db:"total_price" json:"total_price"`
	CreatedAt  pgtype.Timestamp `db:"created_at" json:"created_at"`
	UpdatedAt  pgtype.Timestamp `db:"updated_at" json:"updated_at"`
	DeletedAt  pgtype.Timestamp `db:"deleted_at" json:"deleted_at"`
	TotalCount int64            `db:"total_count" json:"total_count"`
}

type GetOrdersTrashedParams struct {
	Column1 string `db:"column1" json:"column_1"`
	Limit   int32  `db:"limit" json:"limit"`
	Offset  int32  `db:"offset" json:"offset"`
}

type GetOrdersTrashedRow struct {
	OrderID    int32            `db:"order_id" json:"order_id"`
	MerchantID int32            `db:"merchant_id" json:"merchant_id"`
	CashierID  int32            `db:"cashier_id" json:"cashier_id"`
	TotalPrice int64            `db:"total_price" json:"total_price"`
	CreatedAt  pgtype.Timestamp `db:"created_at" json:"created_at"`
	UpdatedAt  pgtype.Timestamp `db:"updated_at" json:"updated_at"`
	DeletedAt  pgtype.Timestamp `db:"deleted_at" json:"deleted_at"`
	TotalCount int64            `db:"total_count" json:"total_count"`
}

type GetYearlyOrderRow struct {
	Year               string  `db:"year" json:"year"`
	OrderCount         int64   `db:"order_count" json:"order_count"`
	TotalRevenue       float64 `db:"total_revenue" json:"total_revenue"`
	TotalItemsSold     int64   `db:"total_items_sold" json:"total_items_sold"`
	ActiveCashiers     int64   `db:"active_cashiers" json:"active_cashiers"`
	UniqueProductsSold int64   `db:"unique_products_sold" json:"unique_products_sold"`
}

type GetYearlyOrderByMerchantParams struct {
	Column1    time.Time `db:"column1" json:"column_1"`
	MerchantID int32     `db:"merchant_id" json:"merchant_id"`
}

type GetYearlyOrderByMerchantRow struct {
	Year               string  `db:"year" json:"year"`
	OrderCount         int64   `db:"order_count" json:"order_count"`
	TotalRevenue       float64 `db:"total_revenue" json:"total_revenue"`
	TotalItemsSold     int64   `db:"total_items_sold" json:"total_items_sold"`
	ActiveCashiers     int64   `db:"active_cashiers" json:"active_cashiers"`
	UniqueProductsSold int64   `db:"unique_products_sold" json:"unique_products_sold"`
}

type GetYearlyTotalRevenueRow struct {
	Year         string `db:"year" json:"year"`
	TotalRevenue int32  `db:"total_revenue" json:"total_revenue"`
}

type GetYearlyTotalRevenueByIdParams struct {
	Column1 int32 `db:"column1" json:"column_1"`
	OrderID int32 `db:"order_id" json:"order_id"`
}

type GetYearlyTotalRevenueByIdRow struct {
	Year         string `db:"year" json:"year"`
	TotalRevenue int32  `db:"total_revenue" json:"total_revenue"`
}

type GetYearlyTotalRevenueByMerchantParams struct {
	Column1    int32 `db:"column1" json:"column_1"`
	MerchantID int32 `db:"merchant_id" json:"merchant_id"`
}

type GetYearlyTotalRevenueByMerchantRow struct {
	Year         string `db:"year" json:"year"`
	TotalRevenue int32  `db:"total_revenue" json:"total_revenue"`
}

type UpdateOrderParams struct {
	OrderID    int32 `db:"order_id" json:"order_id"`
	TotalPrice int64 `db:"total_price" json:"total_price"`
}
