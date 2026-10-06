package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/order_errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
)

const getMonthlyTotalRevenueByMerchant = `WITH monthly_revenue AS (
    SELECT
        EXTRACT(YEAR FROM o.created_at)::TEXT AS year,
        EXTRACT(MONTH FROM o.created_at)::integer AS month,
        COALESCE(SUM(o.total_price), 0)::INTEGER AS total_revenue
    FROM
        orders o
    JOIN
        order_items oi ON o.order_id = oi.order_id
    WHERE
        o.deleted_at IS NULL
        AND oi.deleted_at IS NULL
        AND (
            (o.created_at >= $1 AND o.created_at <= $2)
            OR (o.created_at >= $3 AND o.created_at <= $4)
        )
        AND o.merchant_id = $5
    GROUP BY
        EXTRACT(YEAR FROM o.created_at),
        EXTRACT(MONTH FROM o.created_at)
),
all_months AS (
    SELECT
        EXTRACT(YEAR FROM $1)::TEXT AS year,
        EXTRACT(MONTH FROM $1)::integer AS month,
        TO_CHAR($1, 'FMMonth') AS month_name

    UNION

    SELECT
        EXTRACT(YEAR FROM $3)::TEXT AS year,
        EXTRACT(MONTH FROM $3)::integer AS month,
        TO_CHAR($3, 'FMMonth') AS month_name
)
SELECT
    COALESCE(am.year, EXTRACT(YEAR FROM $1)::TEXT) AS year,
    COALESCE(am.month_name, TO_CHAR($1, 'FMMonth')) AS month,
    COALESCE(mr.total_revenue, 0) AS total_revenue
FROM
    all_months am
LEFT JOIN
    monthly_revenue mr ON am.year = mr.year
                      AND am.month = mr.month
ORDER BY
    am.year DESC,
    am.month DESC`

const getYearlyTotalRevenueByMerchant = `WITH yearly_revenue AS (
    SELECT
        EXTRACT(YEAR FROM o.created_at)::integer AS year,
        COALESCE(SUM(o.total_price), 0)::INTEGER AS total_revenue
    FROM
        orders o
    JOIN
        order_items oi ON o.order_id = oi.order_id
    WHERE
        o.deleted_at IS NULL
        AND oi.deleted_at IS NULL
        AND (
            EXTRACT(YEAR FROM o.created_at) = $1::integer
            OR EXTRACT(YEAR FROM o.created_at) = $1::integer - 1
        )
        AND o.merchant_id = $2
    GROUP BY
        EXTRACT(YEAR FROM o.created_at)
),
all_years AS (
    SELECT $1 AS year
    UNION
    SELECT $1 - 1 AS year
)
SELECT
    ay.year::text AS year,
    COALESCE(yr.total_revenue, 0) AS total_revenue
FROM
    all_years ay
LEFT JOIN
    yearly_revenue yr ON ay.year = yr.year
ORDER BY
    ay.year DESC`

const getMonthlyOrderByMerchant = `WITH date_range AS (
    SELECT
        date_trunc('month', $1::timestamp) AS start_date,
        date_trunc('month', $1::timestamp) + interval '1 year' - interval '1 day' AS end_date
),
monthly_orders AS (
    SELECT
        date_trunc('month', o.created_at) AS activity_month,
        COUNT(o.order_id) AS order_count,
        SUM(o.total_price)::NUMERIC AS total_revenue,
        SUM(oi.quantity) AS total_items_sold
    FROM
        orders o
    JOIN
        order_items oi ON o.order_id = oi.order_id
    WHERE
        o.deleted_at IS NULL
        AND oi.deleted_at IS NULL
        AND o.created_at BETWEEN (SELECT start_date FROM date_range)
                             AND (SELECT end_date FROM date_range)
        AND o.merchant_id = $2
    GROUP BY
        activity_month
)
SELECT
    TO_CHAR(mo.activity_month, 'Mon') AS month,
    mo.order_count,
    mo.total_revenue,
    mo.total_items_sold
FROM
    monthly_orders mo
ORDER BY
    mo.activity_month`

const getYearlyOrderByMerchant = `WITH last_five_years AS (
    SELECT
        EXTRACT(YEAR FROM o.created_at)::text AS year,
        COUNT(o.order_id) AS order_count,
        SUM(o.total_price)::NUMERIC AS total_revenue,
        SUM(oi.quantity) AS total_items_sold,
        COUNT(DISTINCT o.cashier_id) AS active_cashiers,
        COUNT(DISTINCT oi.product_id) AS unique_products_sold
    FROM
        orders o
    JOIN
        order_items oi ON o.order_id = oi.order_id
    WHERE
        o.deleted_at IS NULL
        AND oi.deleted_at IS NULL
        AND EXTRACT(YEAR FROM o.created_at) BETWEEN (EXTRACT(YEAR FROM $1::timestamp) - 4) AND EXTRACT(YEAR FROM $1::timestamp)
        AND o.merchant_id = $2
    GROUP BY
        EXTRACT(YEAR FROM o.created_at)
)
SELECT
    year,
    order_count,
    total_revenue,
    total_items_sold,
    active_cashiers,
    unique_products_sold
FROM
    last_five_years
ORDER BY
    year`

type orderStatsByMerchantRepository struct {
	db *sqlx.DB
}

func NewOrderStatsByMerchantRepository(db *sqlx.DB) OrderStatByMerchantRepository {
	return &orderStatsByMerchantRepository{
		db: db,
	}
}

func (r *orderStatsByMerchantRepository) GetMonthlyTotalRevenueByMerchant(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*db.GetMonthlyTotalRevenueByMerchantRow, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	var rows []*db.GetMonthlyTotalRevenueByMerchantRow
	err := r.db.SelectContext(ctx, &rows, getMonthlyTotalRevenueByMerchant,
		pgtype.Date{Time: currentMonthStart, Valid: true},
		pgtype.Timestamp{Time: currentMonthEnd, Valid: true},
		pgtype.Timestamp{Time: prevMonthStart, Valid: true},
		pgtype.Timestamp{Time: prevMonthEnd, Valid: true},
		int32(req.MerchantID),
	)
	if err != nil {
		return nil, order_errors.ErrGetMonthlyTotalRevenueByMerchant
	}

	return rows, nil
}

func (r *orderStatsByMerchantRepository) GetYearlyTotalRevenueByMerchant(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*db.GetYearlyTotalRevenueByMerchantRow, error) {
	var rows []*db.GetYearlyTotalRevenueByMerchantRow
	err := r.db.SelectContext(ctx, &rows, getYearlyTotalRevenueByMerchant,
		int32(req.Year),
		int32(req.MerchantID),
	)
	if err != nil {
		return nil, order_errors.ErrGetYearlyTotalRevenueByMerchant
	}

	return rows, nil
}

func (r *orderStatsByMerchantRepository) GetMonthlyOrderByMerchant(ctx context.Context, req *requests.MonthOrderMerchant) ([]*db.GetMonthlyOrderByMerchantRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	var rows []*db.GetMonthlyOrderByMerchantRow
	err := r.db.SelectContext(ctx, &rows, getMonthlyOrderByMerchant,
		yearStart,
		int32(req.MerchantID),
	)
	if err != nil {
		return nil, order_errors.ErrGetMonthlyOrderByMerchant
	}

	return rows, nil
}

func (r *orderStatsByMerchantRepository) GetYearlyOrderByMerchant(ctx context.Context, req *requests.YearOrderMerchant) ([]*db.GetYearlyOrderByMerchantRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	var rows []*db.GetYearlyOrderByMerchantRow
	err := r.db.SelectContext(ctx, &rows, getYearlyOrderByMerchant,
		yearStart,
		int32(req.MerchantID),
	)
	if err != nil {
		return nil, order_errors.ErrGetYearlyOrderByMerchant
	}

	return rows, nil
}
