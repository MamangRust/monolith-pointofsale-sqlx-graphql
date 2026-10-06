// Code migrated from sqlc-generated sources; domain types kept for repository layer.
// source: categories.sql

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type CreateCategoryParams struct {
	Name         string  `db:"name" json:"name"`
	Description  *string `db:"description" json:"description"`
	SlugCategory *string `db:"slug_category" json:"slug_category"`
}

type GetCategoriesParams struct {
	Column1 string `db:"column1" json:"column_1"`
	Limit   int32  `db:"limit" json:"limit"`
	Offset  int32  `db:"offset" json:"offset"`
}

type GetCategoriesRow struct {
	CategoryID   int32            `db:"category_id" json:"category_id"`
	Name         string           `db:"name" json:"name"`
	Description  *string          `db:"description" json:"description"`
	SlugCategory *string          `db:"slug_category" json:"slug_category"`
	CreatedAt    pgtype.Timestamp `db:"created_at" json:"created_at"`
	UpdatedAt    pgtype.Timestamp `db:"updated_at" json:"updated_at"`
	DeletedAt    pgtype.Timestamp `db:"deleted_at" json:"deleted_at"`
	TotalCount   int64            `db:"total_count" json:"total_count"`
}

type GetCategoriesActiveParams struct {
	Column1 string `db:"column1" json:"column_1"`
	Limit   int32  `db:"limit" json:"limit"`
	Offset  int32  `db:"offset" json:"offset"`
}

type GetCategoriesActiveRow struct {
	CategoryID   int32            `db:"category_id" json:"category_id"`
	Name         string           `db:"name" json:"name"`
	Description  *string          `db:"description" json:"description"`
	SlugCategory *string          `db:"slug_category" json:"slug_category"`
	CreatedAt    pgtype.Timestamp `db:"created_at" json:"created_at"`
	UpdatedAt    pgtype.Timestamp `db:"updated_at" json:"updated_at"`
	DeletedAt    pgtype.Timestamp `db:"deleted_at" json:"deleted_at"`
	TotalCount   int64            `db:"total_count" json:"total_count"`
}

type GetCategoriesTrashedParams struct {
	Column1 string `db:"column1" json:"column_1"`
	Limit   int32  `db:"limit" json:"limit"`
	Offset  int32  `db:"offset" json:"offset"`
}

type GetCategoriesTrashedRow struct {
	CategoryID   int32            `db:"category_id" json:"category_id"`
	Name         string           `db:"name" json:"name"`
	Description  *string          `db:"description" json:"description"`
	SlugCategory *string          `db:"slug_category" json:"slug_category"`
	CreatedAt    pgtype.Timestamp `db:"created_at" json:"created_at"`
	UpdatedAt    pgtype.Timestamp `db:"updated_at" json:"updated_at"`
	DeletedAt    pgtype.Timestamp `db:"deleted_at" json:"deleted_at"`
	TotalCount   int64            `db:"total_count" json:"total_count"`
}

type GetCategoryByNameAndIdParams struct {
	Name       string `db:"name" json:"name"`
	CategoryID int32  `db:"category_id" json:"category_id"`
}

type GetMonthlyCategoryRow struct {
	Month        string `db:"month" json:"month"`
	CategoryID   int32  `db:"category_id" json:"category_id"`
	CategoryName string `db:"category_name" json:"category_name"`
	OrderCount   int64  `db:"order_count" json:"order_count"`
	ItemsSold    int64  `db:"items_sold" json:"items_sold"`
	TotalRevenue int32  `db:"total_revenue" json:"total_revenue"`
}

type GetMonthlyCategoryByIdParams struct {
	Column1    time.Time `db:"column1" json:"column_1"`
	CategoryID int32     `db:"category_id" json:"category_id"`
}

type GetMonthlyCategoryByIdRow struct {
	Month        string `db:"month" json:"month"`
	CategoryID   int32  `db:"category_id" json:"category_id"`
	CategoryName string `db:"category_name" json:"category_name"`
	OrderCount   int64  `db:"order_count" json:"order_count"`
	ItemsSold    int64  `db:"items_sold" json:"items_sold"`
	TotalRevenue int32  `db:"total_revenue" json:"total_revenue"`
}

type GetMonthlyCategoryByMerchantParams struct {
	Column1    time.Time `db:"column1" json:"column_1"`
	MerchantID int32     `db:"merchant_id" json:"merchant_id"`
}

type GetMonthlyCategoryByMerchantRow struct {
	Month        string `db:"month" json:"month"`
	CategoryID   int32  `db:"category_id" json:"category_id"`
	CategoryName string `db:"category_name" json:"category_name"`
	OrderCount   int64  `db:"order_count" json:"order_count"`
	ItemsSold    int64  `db:"items_sold" json:"items_sold"`
	TotalRevenue int32  `db:"total_revenue" json:"total_revenue"`
}

type GetMonthlyTotalPriceParams struct {
	Extract     pgtype.Date      `db:"extract" json:"extract"`
	CreatedAt   pgtype.Timestamp `db:"created_at" json:"created_at"`
	CreatedAt_2 pgtype.Timestamp `db:"created_at_2" json:"created_at_2"`
	CreatedAt_3 pgtype.Timestamp `db:"created_at_3" json:"created_at_3"`
}

type GetMonthlyTotalPriceRow struct {
	Year         string `db:"year" json:"year"`
	Month        string `db:"month" json:"month"`
	TotalRevenue int32  `db:"total_revenue" json:"total_revenue"`
}

type GetMonthlyTotalPriceByIdParams struct {
	Extract     pgtype.Date      `db:"extract" json:"extract"`
	CreatedAt   pgtype.Timestamp `db:"created_at" json:"created_at"`
	CreatedAt_2 pgtype.Timestamp `db:"created_at_2" json:"created_at_2"`
	CreatedAt_3 pgtype.Timestamp `db:"created_at_3" json:"created_at_3"`
	CategoryID  int32            `db:"category_id" json:"category_id"`
}

type GetMonthlyTotalPriceByIdRow struct {
	Year         string `db:"year" json:"year"`
	Month        string `db:"month" json:"month"`
	TotalRevenue int32  `db:"total_revenue" json:"total_revenue"`
}

type GetMonthlyTotalPriceByMerchantParams struct {
	Extract     pgtype.Date      `db:"extract" json:"extract"`
	CreatedAt   pgtype.Timestamp `db:"created_at" json:"created_at"`
	CreatedAt_2 pgtype.Timestamp `db:"created_at_2" json:"created_at_2"`
	CreatedAt_3 pgtype.Timestamp `db:"created_at_3" json:"created_at_3"`
	MerchantID  int32            `db:"merchant_id" json:"merchant_id"`
}

type GetMonthlyTotalPriceByMerchantRow struct {
	Year         string `db:"year" json:"year"`
	Month        string `db:"month" json:"month"`
	TotalRevenue int32  `db:"total_revenue" json:"total_revenue"`
}

type GetYearlyCategoryRow struct {
	Year               string `db:"year" json:"year"`
	CategoryID         int32  `db:"category_id" json:"category_id"`
	CategoryName       string `db:"category_name" json:"category_name"`
	OrderCount         int64  `db:"order_count" json:"order_count"`
	ItemsSold          int64  `db:"items_sold" json:"items_sold"`
	TotalRevenue       int32  `db:"total_revenue" json:"total_revenue"`
	UniqueProductsSold int64  `db:"unique_products_sold" json:"unique_products_sold"`
}

type GetYearlyCategoryByIdParams struct {
	Column1    time.Time `db:"column1" json:"column_1"`
	CategoryID int32     `db:"category_id" json:"category_id"`
}

type GetYearlyCategoryByIdRow struct {
	Year               string `db:"year" json:"year"`
	CategoryID         int32  `db:"category_id" json:"category_id"`
	CategoryName       string `db:"category_name" json:"category_name"`
	OrderCount         int64  `db:"order_count" json:"order_count"`
	ItemsSold          int64  `db:"items_sold" json:"items_sold"`
	TotalRevenue       int32  `db:"total_revenue" json:"total_revenue"`
	UniqueProductsSold int64  `db:"unique_products_sold" json:"unique_products_sold"`
}

type GetYearlyCategoryByMerchantParams struct {
	Column1    time.Time `db:"column1" json:"column_1"`
	MerchantID int32     `db:"merchant_id" json:"merchant_id"`
}

type GetYearlyCategoryByMerchantRow struct {
	Year               string `db:"year" json:"year"`
	CategoryID         int32  `db:"category_id" json:"category_id"`
	CategoryName       string `db:"category_name" json:"category_name"`
	OrderCount         int64  `db:"order_count" json:"order_count"`
	ItemsSold          int64  `db:"items_sold" json:"items_sold"`
	TotalRevenue       int32  `db:"total_revenue" json:"total_revenue"`
	UniqueProductsSold int64  `db:"unique_products_sold" json:"unique_products_sold"`
}

type GetYearlyTotalPriceRow struct {
	Year         string `db:"year" json:"year"`
	TotalRevenue int32  `db:"total_revenue" json:"total_revenue"`
}

type GetYearlyTotalPriceByIdParams struct {
	Column1    int32 `db:"column1" json:"column_1"`
	CategoryID int32 `db:"category_id" json:"category_id"`
}

type GetYearlyTotalPriceByIdRow struct {
	Year         string `db:"year" json:"year"`
	TotalRevenue int32  `db:"total_revenue" json:"total_revenue"`
}

type GetYearlyTotalPriceByMerchantParams struct {
	Column1    int32 `db:"column1" json:"column_1"`
	MerchantID int32 `db:"merchant_id" json:"merchant_id"`
}

type GetYearlyTotalPriceByMerchantRow struct {
	Year         string `db:"year" json:"year"`
	TotalRevenue int32  `db:"total_revenue" json:"total_revenue"`
}

type UpdateCategoryParams struct {
	CategoryID   int32   `db:"category_id" json:"category_id"`
	Name         string  `db:"name" json:"name"`
	Description  *string `db:"description" json:"description"`
	SlugCategory *string `db:"slug_category" json:"slug_category"`
}
