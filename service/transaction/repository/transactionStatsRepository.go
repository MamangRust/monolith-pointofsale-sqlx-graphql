package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/transaction_errors"
	"github.com/jmoiron/sqlx"
)

const getMonthlyAmountTransactionSuccess = `WITH
    monthly_data AS (
        SELECT
            EXTRACT(
                YEAR
                FROM t.created_at
            )::integer AS year,
            EXTRACT(
                MONTH
                FROM t.created_at
            )::integer AS month,
            COUNT(*) AS total_success,
            COALESCE(SUM(t.amount), 0)::integer AS total_amount
        FROM transactions t
        WHERE
            t.deleted_at IS NULL
            AND t.payment_status = 'success'
            AND (
                (
                    t.created_at >= $1::timestamp
                    AND t.created_at <= $2::timestamp
                )
                OR (
                    t.created_at >= $3::timestamp
                    AND t.created_at <= $4::timestamp
                )
            )
        GROUP BY
            EXTRACT(
                YEAR
                FROM t.created_at
            ),
            EXTRACT(
                MONTH
                FROM t.created_at
            )
    ),
    formatted_data AS (
        SELECT
            year::text,
            TO_CHAR(
                TO_DATE(month::text, 'MM'),
                'Mon'
            ) AS month,
            total_success,
            total_amount
        FROM monthly_data
        UNION ALL
        SELECT
            EXTRACT(
                YEAR
                FROM $1::timestamp
            )::text AS year,
            TO_CHAR($1::timestamp, 'Mon') AS month,
            0 AS total_success,
            0 AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM monthly_data
                WHERE
                    year = EXTRACT(
                        YEAR
                        FROM $1::timestamp
                    )::integer
                    AND month = EXTRACT(
                        MONTH
                        FROM $1::timestamp
                    )::integer
            )
        UNION ALL
        SELECT
            EXTRACT(
                YEAR
                FROM $3::timestamp
            )::text AS year,
            TO_CHAR($3::timestamp, 'Mon') AS month,
            0 AS total_success,
            0 AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM monthly_data
                WHERE
                    year = EXTRACT(
                        YEAR
                        FROM $3::timestamp
                    )::integer
                    AND month = EXTRACT(
                        MONTH
                        FROM $3::timestamp
                    )::integer
            )
    )
SELECT *
FROM formatted_data
ORDER BY year DESC, TO_DATE(month, 'Mon') DESC`

const getYearlyAmountTransactionSuccess = `WITH
    yearly_data AS (
        SELECT
            EXTRACT(
                YEAR
                FROM t.created_at
            )::integer AS year,
            COUNT(*) AS total_success,
            COALESCE(SUM(t.amount), 0)::integer AS total_amount
        FROM transactions t
        WHERE
            t.deleted_at IS NULL
            AND t.payment_status = 'success'
            AND (
                EXTRACT(
                    YEAR
                    FROM t.created_at
                ) = $1::integer
                OR EXTRACT(
                    YEAR
                    FROM t.created_at
                ) = $1::integer - 1
            )
        GROUP BY
            EXTRACT(
                YEAR
                FROM t.created_at
            )
    ),
    formatted_data AS (
        SELECT
            year::text,
            total_success::integer,
            total_amount::integer
        FROM yearly_data
        UNION ALL
        SELECT
            $1::text AS year,
            0::integer AS total_success,
            0::integer AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM yearly_data
                WHERE
                    year = $1::integer
            )
        UNION ALL
        SELECT ($1::integer - 1)::text AS year,
            0::integer AS total_success,
            0::integer AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM yearly_data
                WHERE
                    year = $1::integer - 1
            )
    )
SELECT *
FROM formatted_data
ORDER BY year DESC`

const getMonthlyAmountTransactionFailed = `WITH
    monthly_data AS (
        SELECT
            EXTRACT(
                YEAR
                FROM t.created_at
            )::integer AS year,
            EXTRACT(
                MONTH
                FROM t.created_at
            )::integer AS month,
            COUNT(*) AS total_failed,
            COALESCE(SUM(t.amount), 0)::integer AS total_amount
        FROM transactions t
        WHERE
            t.deleted_at IS NULL
            AND t.payment_status = 'failed'
            AND (
                (
                    t.created_at >= $1::timestamp
                    AND t.created_at <= $2::timestamp
                )
                OR (
                    t.created_at >= $3::timestamp
                    AND t.created_at <= $4::timestamp
                )
            )
        GROUP BY
            EXTRACT(
                YEAR
                FROM t.created_at
            ),
            EXTRACT(
                MONTH
                FROM t.created_at
            )
    ),
    formatted_data AS (
        SELECT
            year::text,
            TO_CHAR(
                TO_DATE(month::text, 'MM'),
                'Mon'
            ) AS month,
            total_failed,
            total_amount
        FROM monthly_data
        UNION ALL
        SELECT
            EXTRACT(
                YEAR
                FROM $1::timestamp
            )::text AS year,
            TO_CHAR($1::timestamp, 'Mon') AS month,
            0 AS total_failed,
            0 AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM monthly_data
                WHERE
                    year = EXTRACT(
                        YEAR
                        FROM $1::timestamp
                    )::integer
                    AND month = EXTRACT(
                        MONTH
                        FROM $1::timestamp
                    )::integer
            )
        UNION ALL
        SELECT
            EXTRACT(
                YEAR
                FROM $3::timestamp
            )::text AS year,
            TO_CHAR($3::timestamp, 'Mon') AS month,
            0 AS total_failed,
            0 AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM monthly_data
                WHERE
                    year = EXTRACT(
                        YEAR
                        FROM $3::timestamp
                    )::integer
                    AND month = EXTRACT(
                        MONTH
                        FROM $3::timestamp
                    )::integer
            )
    )
SELECT *
FROM formatted_data
ORDER BY year DESC, TO_DATE(month, 'Mon') DESC`

const getYearlyAmountTransactionFailed = `WITH
    yearly_data AS (
        SELECT
            EXTRACT(
                YEAR
                FROM t.created_at
            )::integer AS year,
            COUNT(*) AS total_failed,
            COALESCE(SUM(t.amount), 0)::integer AS total_amount
        FROM transactions t
        WHERE
            t.deleted_at IS NULL
            AND t.payment_status = 'failed'
            AND (
                EXTRACT(
                    YEAR
                    FROM t.created_at
                ) = $1::integer
                OR EXTRACT(
                    YEAR
                    FROM t.created_at
                ) = $1::integer - 1
            )
        GROUP BY
            EXTRACT(
                YEAR
                FROM t.created_at
            )
    ),
    formatted_data AS (
        SELECT
            year::text,
            total_failed::integer,
            total_amount::integer
        FROM yearly_data
        UNION ALL
        SELECT
            $1::text AS year,
            0::integer AS total_failed,
            0::integer AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM yearly_data
                WHERE
                    year = $1::integer
            )
        UNION ALL
        SELECT ($1::integer - 1)::text AS year,
            0::integer AS total_failed,
            0::integer AS total_amount
        WHERE
            NOT EXISTS (
                SELECT 1
                FROM yearly_data
                WHERE
                    year = $1::integer - 1
            )
    )
SELECT *
FROM formatted_data
ORDER BY year DESC`

const getMonthlyTransactionMethodsSuccess = `WITH
    date_ranges AS (
        SELECT
            $1::timestamp AS range1_start,
            $2::timestamp AS range1_end,
            $3::timestamp AS range2_start,
            $4::timestamp AS range2_end
    ),
    payment_methods AS (
        SELECT DISTINCT
            payment_method
        FROM transactions
        WHERE deleted_at IS NULL
    ),
    all_months AS (
        SELECT generate_series(
            date_trunc('month', LEAST(
                (SELECT range1_start FROM date_ranges),
                (SELECT range2_start FROM date_ranges)
            )),
            date_trunc('month', GREATEST(
                (SELECT range1_end FROM date_ranges),
                (SELECT range2_end FROM date_ranges)
            )),
            interval '1 month'
        )::date AS activity_month
    ),
    all_combinations AS (
        SELECT
            am.activity_month,
            pm.payment_method
        FROM all_months am
        CROSS JOIN payment_methods pm
    ),
    monthly_transactions AS (
        SELECT
            date_trunc('month', t.created_at)::date AS activity_month,
            t.payment_method,
            COUNT(t.transaction_id) AS total_transactions,
            COALESCE(SUM(t.amount), 0)::NUMERIC AS total_amount
        FROM transactions t
        JOIN date_ranges dr ON (
            t.created_at BETWEEN dr.range1_start AND dr.range1_end
            OR t.created_at BETWEEN dr.range2_start AND dr.range2_end
        )
        WHERE
            t.deleted_at IS NULL
            AND t.payment_status = 'success'
        GROUP BY
            date_trunc('month', t.created_at),
            t.payment_method
    )
SELECT
    TO_CHAR(ac.activity_month, 'Mon') AS month,
    ac.payment_method,
    COALESCE(mt.total_transactions, 0) AS total_transactions,
    COALESCE(mt.total_amount, 0) AS total_amount
FROM all_combinations ac
LEFT JOIN monthly_transactions mt ON
    ac.activity_month = mt.activity_month
    AND ac.payment_method = mt.payment_method
ORDER BY
    ac.activity_month,
    ac.payment_method`

const getMonthlyTransactionMethodsFailed = `WITH
    date_ranges AS (
        SELECT
            $1::timestamp AS range1_start,
            $2::timestamp AS range1_end,
            $3::timestamp AS range2_start,
            $4::timestamp AS range2_end
    ),
    payment_methods AS (
        SELECT DISTINCT
            payment_method
        FROM transactions
        WHERE deleted_at IS NULL
    ),
    all_months AS (
        SELECT generate_series(
            date_trunc('month', LEAST(
                (SELECT range1_start FROM date_ranges),
                (SELECT range2_start FROM date_ranges)
            )),
            date_trunc('month', GREATEST(
                (SELECT range1_end FROM date_ranges),
                (SELECT range2_end FROM date_ranges)
            )),
            interval '1 month'
        )::date AS activity_month
    ),
    all_combinations AS (
        SELECT
            am.activity_month,
            pm.payment_method
        FROM all_months am
        CROSS JOIN payment_methods pm
    ),
    monthly_transactions AS (
        SELECT
            date_trunc('month', t.created_at)::date AS activity_month,
            t.payment_method,
            COUNT(t.transaction_id) AS total_transactions,
            COALESCE(SUM(t.amount), 0)::NUMERIC AS total_amount
        FROM transactions t
        JOIN date_ranges dr ON (
            t.created_at BETWEEN dr.range1_start AND dr.range1_end
            OR t.created_at BETWEEN dr.range2_start AND dr.range2_end
        )
        WHERE
            t.deleted_at IS NULL
            AND t.payment_status = 'failed'
        GROUP BY
            date_trunc('month', t.created_at),
            t.payment_method
    )
SELECT
    TO_CHAR(ac.activity_month, 'Mon') AS month,
    ac.payment_method,
    COALESCE(mt.total_transactions, 0) AS total_transactions,
    COALESCE(mt.total_amount, 0) AS total_amount
FROM all_combinations ac
LEFT JOIN monthly_transactions mt ON
    ac.activity_month = mt.activity_month
    AND ac.payment_method = mt.payment_method
ORDER BY
    ac.activity_month,
    ac.payment_method`

const getYearlyTransactionMethodsSuccess = `WITH
    year_range AS (
        SELECT
            EXTRACT(YEAR FROM $1::timestamp)::int - 1 AS start_year,
            EXTRACT(YEAR FROM $1::timestamp)::int AS end_year
    ),
    payment_methods AS (
        SELECT DISTINCT
            payment_method
        FROM transactions
        WHERE deleted_at IS NULL
    ),
    all_years AS (
        SELECT generate_series(
            (SELECT start_year FROM year_range),
            (SELECT end_year FROM year_range)
        )::int AS year
    ),
    all_combinations AS (
        SELECT
            ay.year::text AS year,
            pm.payment_method
        FROM all_years ay
        CROSS JOIN payment_methods pm
    ),
    yearly_transactions AS (
        SELECT
            EXTRACT(YEAR FROM t.created_at)::text AS year,
            t.payment_method,
            COUNT(t.transaction_id) AS total_transactions,
            COALESCE(SUM(t.amount), 0)::NUMERIC AS total_amount
        FROM transactions t
        WHERE
            t.deleted_at IS NULL
            AND t.payment_status = 'success'
            AND EXTRACT(YEAR FROM t.created_at) BETWEEN (SELECT start_year FROM year_range) AND (SELECT end_year FROM year_range)
        GROUP BY
            EXTRACT(YEAR FROM t.created_at),
            t.payment_method
    )
SELECT
    ac.year,
    ac.payment_method,
    COALESCE(yt.total_transactions, 0) AS total_transactions,
    COALESCE(yt.total_amount, 0) AS total_amount
FROM all_combinations ac
LEFT JOIN yearly_transactions yt ON
    ac.year = yt.year
    AND ac.payment_method = yt.payment_method
ORDER BY
    ac.year,
    ac.payment_method`

const getYearlyTransactionMethodsFailed = `WITH
    year_range AS (
        SELECT
            EXTRACT(YEAR FROM $1::timestamp)::int - 1 AS start_year,
            EXTRACT(YEAR FROM $1::timestamp)::int AS end_year
    ),
    payment_methods AS (
        SELECT DISTINCT
            payment_method
        FROM transactions
        WHERE deleted_at IS NULL
    ),
    all_years AS (
        SELECT generate_series(
            (SELECT start_year FROM year_range),
            (SELECT end_year FROM year_range)
        )::int AS year
    ),
    all_combinations AS (
        SELECT
            ay.year::text AS year,
            pm.payment_method
        FROM all_years ay
        CROSS JOIN payment_methods pm
    ),
    yearly_transactions AS (
        SELECT
            EXTRACT(YEAR FROM t.created_at)::text AS year,
            t.payment_method,
            COUNT(t.transaction_id) AS total_transactions,
            COALESCE(SUM(t.amount), 0)::NUMERIC AS total_amount
        FROM transactions t
        WHERE
            t.deleted_at IS NULL
            AND t.payment_status = 'failed'
            AND EXTRACT(YEAR FROM t.created_at) BETWEEN (SELECT start_year FROM year_range) AND (SELECT end_year FROM year_range)
        GROUP BY
            EXTRACT(YEAR FROM t.created_at),
            t.payment_method
    )
SELECT
    ac.year,
    ac.payment_method,
    COALESCE(yt.total_transactions, 0) AS total_transactions,
    COALESCE(yt.total_amount, 0) AS total_amount
FROM all_combinations ac
LEFT JOIN yearly_transactions yt ON
    ac.year = yt.year
    AND ac.payment_method = yt.payment_method
ORDER BY
    ac.year,
    ac.payment_method`

type transactonStatsRepository struct {
	db *sqlx.DB
}

func NewTransactionStatsRepository(db *sqlx.DB) *transactonStatsRepository {
	return &transactonStatsRepository{
		db: db,
	}
}

func (r *transactonStatsRepository) GetMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction) ([]*db.GetMonthlyAmountTransactionSuccessRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	var rows []*db.GetMonthlyAmountTransactionSuccessRow
	err := r.db.SelectContext(ctx, &rows, getMonthlyAmountTransactionSuccess,
		currentDate,
		lastDayCurrentMonth,
		prevDate,
		lastDayPrevMonth,
	)
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountSuccess
	}

	return rows, nil
}

func (r *transactonStatsRepository) GetYearlyAmountSuccess(ctx context.Context, year int) ([]*db.GetYearlyAmountTransactionSuccessRow, error) {
	var rows []*db.GetYearlyAmountTransactionSuccessRow
	err := r.db.SelectContext(ctx, &rows, getYearlyAmountTransactionSuccess, int32(year))
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountSuccess
	}

	return rows, nil
}

func (r *transactonStatsRepository) GetMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction) ([]*db.GetMonthlyAmountTransactionFailedRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	var rows []*db.GetMonthlyAmountTransactionFailedRow
	err := r.db.SelectContext(ctx, &rows, getMonthlyAmountTransactionFailed,
		currentDate,
		lastDayCurrentMonth,
		prevDate,
		lastDayPrevMonth,
	)
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountFailed
	}

	return rows, nil
}

func (r *transactonStatsRepository) GetYearlyAmountFailed(ctx context.Context, year int) ([]*db.GetYearlyAmountTransactionFailedRow, error) {
	var rows []*db.GetYearlyAmountTransactionFailedRow
	err := r.db.SelectContext(ctx, &rows, getYearlyAmountTransactionFailed, int32(year))
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountFailed
	}

	return rows, nil
}

func (r *transactonStatsRepository) GetMonthlyTransactionMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction) ([]*db.GetMonthlyTransactionMethodsSuccessRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	var rows []*db.GetMonthlyTransactionMethodsSuccessRow
	err := r.db.SelectContext(ctx, &rows, getMonthlyTransactionMethodsSuccess,
		currentDate,
		lastDayCurrentMonth,
		prevDate,
		lastDayPrevMonth,
	)
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethod
	}

	return rows, nil
}

func (r *transactonStatsRepository) GetYearlyTransactionMethodSuccess(ctx context.Context, year int) ([]*db.GetYearlyTransactionMethodsSuccessRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	var rows []*db.GetYearlyTransactionMethodsSuccessRow
	err := r.db.SelectContext(ctx, &rows, getYearlyTransactionMethodsSuccess, yearStart)
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethod
	}

	return rows, nil
}

func (r *transactonStatsRepository) GetMonthlyTransactionMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction) ([]*db.GetMonthlyTransactionMethodsFailedRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	var rows []*db.GetMonthlyTransactionMethodsFailedRow
	err := r.db.SelectContext(ctx, &rows, getMonthlyTransactionMethodsFailed,
		currentDate,
		lastDayCurrentMonth,
		prevDate,
		lastDayPrevMonth,
	)
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethod
	}

	return rows, nil
}

func (r *transactonStatsRepository) GetYearlyTransactionMethodFailed(ctx context.Context, year int) ([]*db.GetYearlyTransactionMethodsFailedRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	var rows []*db.GetYearlyTransactionMethodsFailedRow
	err := r.db.SelectContext(ctx, &rows, getYearlyTransactionMethodsFailed, yearStart)
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethod
	}

	return rows, nil
}
