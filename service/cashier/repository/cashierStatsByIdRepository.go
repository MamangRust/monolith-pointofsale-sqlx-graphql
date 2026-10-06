package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/cashier_errors"
	"github.com/jmoiron/sqlx"
)

const getMonthlyTotalSalesById = `WITH monthly_totals AS (
    SELECT
        EXTRACT(YEAR FROM o.created_at)::TEXT AS year,
        EXTRACT(MONTH FROM o.created_at)::integer AS month,
        COALESCE(SUM(o.total_price), 0)::INTEGER AS total_sales
    FROM
        orders o
    JOIN
        cashiers c ON o.cashier_id = c.cashier_id
    WHERE
        o.deleted_at IS NULL
        AND c.deleted_at IS NULL
        AND (
            (o.created_at >= $1 AND o.created_at <= $2)
            OR (o.created_at >= $3 AND o.created_at <= $4)
        )
        AND c.cashier_id = $5
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
    COALESCE(mt.total_sales, 0) AS total_sales
FROM
    all_months am
LEFT JOIN
    monthly_totals mt ON am.year = mt.year AND am.month = mt.month
ORDER BY
    am.year::INT DESC,
    am.month DESC`

const getYearlyTotalSalesById = `WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM o.created_at)::integer AS year,
        COALESCE(SUM(o.total_price), 0)::INTEGER AS total_sales
    FROM
        orders o
    JOIN
        cashiers c ON o.cashier_id = c.cashier_id
    WHERE
        o.deleted_at IS NULL
        AND c.deleted_at IS NULL
        AND (
            EXTRACT(YEAR FROM o.created_at) = $1::integer
            OR EXTRACT(YEAR FROM o.created_at) = $1::integer - 1
        )
        AND c.cashier_id = $2
    GROUP BY
        EXTRACT(YEAR FROM o.created_at)
),
all_years AS (
    SELECT $1::integer AS year
    UNION
    SELECT $1::integer - 1 AS year
)
SELECT
    a.year::text AS year,
    COALESCE(yd.total_sales, 0) AS total_sales
FROM
    all_years a
LEFT JOIN
    yearly_data yd ON a.year = yd.year
ORDER BY
    a.year DESC`

const getMonthlyCashierByCashierId = `WITH date_range AS (
    SELECT
        date_trunc('month', $1::timestamp) AS start_date,
        date_trunc('month', $1::timestamp) + interval '1 year' - interval '1 day' AS end_date
),
cashier_activity AS (
    SELECT
        c.cashier_id,
        c.name AS cashier_name,
        date_trunc('month', o.created_at) AS activity_month,
        COUNT(o.order_id) AS order_count,
        SUM(o.total_price) AS total_sales
    FROM
        orders o
    JOIN
        cashiers c ON o.cashier_id = c.cashier_id
    WHERE
        o.deleted_at IS NULL
        AND c.deleted_at IS NULL
        AND c.cashier_id = $2
        AND o.created_at BETWEEN (SELECT start_date FROM date_range)
                             AND (SELECT end_date FROM date_range)
    GROUP BY
        c.cashier_id, c.name, activity_month
)
SELECT
    ca.cashier_id,
    ca.cashier_name,
    TO_CHAR(ca.activity_month, 'Mon') AS month,
    ca.order_count,
    ca.total_sales
FROM
    cashier_activity ca
ORDER BY
    ca.activity_month, ca.cashier_id`

const getYearlyCashierByCashierId = `WITH last_five_years AS (
    SELECT
        c.cashier_id,
        c.name AS cashier_name,
        EXTRACT(YEAR FROM o.created_at)::text AS year,
        COUNT(o.order_id) AS order_count,
        SUM(o.total_price) AS total_sales
    FROM
        orders o
    JOIN
        cashiers c ON o.cashier_id = c.cashier_id
    WHERE
        o.deleted_at IS NULL
        AND c.deleted_at IS NULL
        AND c.cashier_id = $2
        AND EXTRACT(YEAR FROM o.created_at) BETWEEN (EXTRACT(YEAR FROM $1::timestamp) - 4) AND EXTRACT(YEAR FROM $1::timestamp)
    GROUP BY
        c.cashier_id, c.name, EXTRACT(YEAR FROM o.created_at)
)
SELECT
    year,
    cashier_id,
    cashier_name,
    order_count,
    total_sales
FROM
    last_five_years
ORDER BY
    year, cashier_id`

type cashierStatsByIdRepository struct {
	db *sqlx.DB
}

func NewCashierStatsByIdRepository(db *sqlx.DB) CashierStatByIdRepository {
	return &cashierStatsByIdRepository{
		db: db,
	}
}

func (r *cashierStatsByIdRepository) GetMonthlyTotalSalesById(ctx context.Context, req *requests.MonthTotalSalesCashier) ([]*db.GetMonthlyTotalSalesByIdRow, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)

	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	var rows []*db.GetMonthlyTotalSalesByIdRow
	err := r.db.SelectContext(ctx, &rows, getMonthlyTotalSalesById,
		currentMonthStart,
		currentMonthEnd,
		prevMonthStart,
		prevMonthEnd,
		int32(req.CashierID),
	)
	if err != nil {
		return nil, cashier_errors.ErrGetMonthlyTotalSalesById
	}

	return rows, nil
}

func (r *cashierStatsByIdRepository) GetYearlyTotalSalesById(ctx context.Context, req *requests.YearTotalSalesCashier) ([]*db.GetYearlyTotalSalesByIdRow, error) {
	var rows []*db.GetYearlyTotalSalesByIdRow
	err := r.db.SelectContext(ctx, &rows, getYearlyTotalSalesById,
		int32(req.Year),
		int32(req.CashierID),
	)
	if err != nil {
		return nil, cashier_errors.ErrGetYearlyTotalSalesById
	}

	return rows, nil
}

func (r *cashierStatsByIdRepository) GetMonthlyCashierById(ctx context.Context, req *requests.MonthCashierId) ([]*db.GetMonthlyCashierByCashierIdRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	var rows []*db.GetMonthlyCashierByCashierIdRow
	err := r.db.SelectContext(ctx, &rows, getMonthlyCashierByCashierId,
		yearStart,
		int32(req.CashierID),
	)
	if err != nil {
		return nil, cashier_errors.ErrGetMonthlyCashierById
	}

	return rows, nil
}

func (r *cashierStatsByIdRepository) GetYearlyCashierById(ctx context.Context, req *requests.YearCashierId) ([]*db.GetYearlyCashierByCashierIdRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	var rows []*db.GetYearlyCashierByCashierIdRow
	err := r.db.SelectContext(ctx, &rows, getYearlyCashierByCashierId,
		yearStart,
		int32(req.CashierID),
	)
	if err != nil {
		return nil, cashier_errors.ErrGetYearlyCashierById
	}

	return rows, nil
}
