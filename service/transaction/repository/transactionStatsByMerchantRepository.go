package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	"github.com/MamangRust/monolith-pointofsale-shared/errors/transaction_errors"
	"github.com/jmoiron/sqlx"
)

const getMonthlyAmountTransactionSuccessByMerchant = `WITH
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
            AND t.merchant_id = $5
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

const getYearlyAmountTransactionSuccessByMerchant = `WITH
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
            AND t.merchant_id = $2
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

const getMonthlyAmountTransactionFailedByMerchant = `WITH
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
            AND t.merchant_id = $5
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

const getYearlyAmountTransactionFailedByMerchant = `WITH
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
            AND t.merchant_id = $2
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

const getMonthlyTransactionMethodsByMerchantSuccess = `WITH
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
            AND t.merchant_id = $5
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

const getMonthlyTransactionMethodsByMerchantFailed = `WITH
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
            AND t.merchant_id = $5
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

const getYearlyTransactionMethodsByMerchantSuccess = `WITH
    year_series AS (
        SELECT generate_series(
            EXTRACT(YEAR FROM $1::timestamp)::integer - 2,
            EXTRACT(YEAR FROM $1::timestamp)::integer,
            1
        ) AS year
    ),
    yearly_transactions AS (
        SELECT
            EXTRACT(YEAR FROM t.created_at)::integer AS year,
            t.payment_method,
            COUNT(t.transaction_id) AS total_transactions,
            SUM(t.amount)::NUMERIC AS total_amount
        FROM transactions t
        WHERE
            t.deleted_at IS NULL
            AND t.payment_status = 'success'
            AND t.merchant_id = $2
            AND EXTRACT(YEAR FROM t.created_at) BETWEEN (EXTRACT(YEAR FROM $1::timestamp) - 1) AND EXTRACT(YEAR FROM $1::timestamp)
        GROUP BY
            year,
            t.payment_method
    ),
    payment_methods AS (
        SELECT DISTINCT payment_method
        FROM transactions
        WHERE deleted_at IS NULL
    )
SELECT
    ys.year::text AS year,
    pm.payment_method,
    COALESCE(yt.total_transactions, 0) AS total_transactions,
    COALESCE(yt.total_amount, 0) AS total_amount
FROM year_series ys
CROSS JOIN payment_methods pm
LEFT JOIN yearly_transactions yt
    ON ys.year = yt.year
    AND pm.payment_method = yt.payment_method
ORDER BY ys.year, pm.payment_method`

const getYearlyTransactionMethodsByMerchantFailed = `WITH
    year_series AS (
        SELECT generate_series(
            EXTRACT(YEAR FROM $1::timestamp)::integer - 1,
            EXTRACT(YEAR FROM $1::timestamp)::integer,
            1
        ) AS year
    ),
    yearly_transactions AS (
        SELECT
            EXTRACT(YEAR FROM t.created_at)::integer AS year,
            t.payment_method,
            COUNT(t.transaction_id) AS total_transactions,
            SUM(t.amount)::NUMERIC AS total_amount
        FROM transactions t
        WHERE
            t.deleted_at IS NULL
            AND t.payment_status = 'failed'
            AND t.merchant_id = $2
            AND EXTRACT(YEAR FROM t.created_at) BETWEEN (EXTRACT(YEAR FROM $1::timestamp) - 1) AND EXTRACT(YEAR FROM $1::timestamp)
        GROUP BY
            year,
            t.payment_method
    ),
    payment_methods AS (
        SELECT DISTINCT payment_method
        FROM transactions
        WHERE deleted_at IS NULL
    )
SELECT
    ys.year::text AS year,
    pm.payment_method,
    COALESCE(yt.total_transactions, 0) AS total_transactions,
    COALESCE(yt.total_amount, 0) AS total_amount
FROM year_series ys
CROSS JOIN payment_methods pm
LEFT JOIN yearly_transactions yt
    ON ys.year = yt.year
    AND pm.payment_method = yt.payment_method
ORDER BY ys.year, pm.payment_method`

type transactionStatsByMerchantRepository struct {
	db *sqlx.DB
}

func NewTransactionStatsByMerchantRepository(db *sqlx.DB) *transactionStatsByMerchantRepository {
	return &transactionStatsByMerchantRepository{
		db: db,
	}
}

func (r *transactionStatsByMerchantRepository) GetMonthlyAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*db.GetMonthlyAmountTransactionSuccessByMerchantRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	var rows []*db.GetMonthlyAmountTransactionSuccessByMerchantRow
	err := r.db.SelectContext(ctx, &rows, getMonthlyAmountTransactionSuccessByMerchant,
		currentDate,
		lastDayCurrentMonth,
		prevDate,
		lastDayPrevMonth,
		int32(req.MerchantID),
	)
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountSuccessByMerchant
	}

	return rows, nil
}

func (r *transactionStatsByMerchantRepository) GetYearlyAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*db.GetYearlyAmountTransactionSuccessByMerchantRow, error) {
	var rows []*db.GetYearlyAmountTransactionSuccessByMerchantRow
	err := r.db.SelectContext(ctx, &rows, getYearlyAmountTransactionSuccessByMerchant,
		int32(req.Year),
		int32(req.MerchantID),
	)
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountSuccessByMerchant
	}

	return rows, nil
}

func (r *transactionStatsByMerchantRepository) GetMonthlyAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*db.GetMonthlyAmountTransactionFailedByMerchantRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	var rows []*db.GetMonthlyAmountTransactionFailedByMerchantRow
	err := r.db.SelectContext(ctx, &rows, getMonthlyAmountTransactionFailedByMerchant,
		currentDate,
		lastDayCurrentMonth,
		prevDate,
		lastDayPrevMonth,
		int32(req.MerchantID),
	)
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountFailedByMerchant
	}

	return rows, nil
}

func (r *transactionStatsByMerchantRepository) GetYearlyAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*db.GetYearlyAmountTransactionFailedByMerchantRow, error) {
	var rows []*db.GetYearlyAmountTransactionFailedByMerchantRow
	err := r.db.SelectContext(ctx, &rows, getYearlyAmountTransactionFailedByMerchant,
		int32(req.Year),
		int32(req.MerchantID),
	)
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountFailedByMerchant
	}

	return rows, nil
}

func (r *transactionStatsByMerchantRepository) GetMonthlyTransactionMethodByMerchantSuccess(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*db.GetMonthlyTransactionMethodsByMerchantSuccessRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	var rows []*db.GetMonthlyTransactionMethodsByMerchantSuccessRow
	err := r.db.SelectContext(ctx, &rows, getMonthlyTransactionMethodsByMerchantSuccess,
		currentDate,
		lastDayCurrentMonth,
		prevDate,
		lastDayPrevMonth,
		int32(req.MerchantID),
	)
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethodByMerchant
	}

	return rows, nil
}

func (r *transactionStatsByMerchantRepository) GetYearlyTransactionMethodByMerchantSuccess(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*db.GetYearlyTransactionMethodsByMerchantSuccessRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	var rows []*db.GetYearlyTransactionMethodsByMerchantSuccessRow
	err := r.db.SelectContext(ctx, &rows, getYearlyTransactionMethodsByMerchantSuccess,
		yearStart,
		int32(req.MerchantID),
	)
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethodByMerchant
	}

	return rows, nil
}

func (r *transactionStatsByMerchantRepository) GetMonthlyTransactionMethodByMerchantFailed(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*db.GetMonthlyTransactionMethodsByMerchantFailedRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	var rows []*db.GetMonthlyTransactionMethodsByMerchantFailedRow
	err := r.db.SelectContext(ctx, &rows, getMonthlyTransactionMethodsByMerchantFailed,
		currentDate,
		lastDayCurrentMonth,
		prevDate,
		lastDayPrevMonth,
		int32(req.MerchantID),
	)
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethodByMerchant
	}

	return rows, nil
}

func (r *transactionStatsByMerchantRepository) GetYearlyTransactionMethodByMerchantFailed(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*db.GetYearlyTransactionMethodsByMerchantFailedRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	var rows []*db.GetYearlyTransactionMethodsByMerchantFailedRow
	err := r.db.SelectContext(ctx, &rows, getYearlyTransactionMethodsByMerchantFailed,
		yearStart,
		int32(req.MerchantID),
	)
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethodByMerchant
	}

	return rows, nil
}
