// Code migrated from sqlc-generated sources; domain types kept for repository layer.
// source: merchant.sql

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateMerchantParams struct {
	UserID       int32   `db:"user_id" json:"user_id"`
	Name         string  `db:"name" json:"name"`
	Description  *string `db:"description" json:"description"`
	Address      *string `db:"address" json:"address"`
	ContactEmail *string `db:"contact_email" json:"contact_email"`
	ContactPhone *string `db:"contact_phone" json:"contact_phone"`
	Status       string  `db:"status" json:"status"`
}

type GetMerchantsParams struct {
	Column1 string `db:"column1" json:"column_1"`
	Limit   int32  `db:"limit" json:"limit"`
	Offset  int32  `db:"offset" json:"offset"`
}

type GetMerchantsRow struct {
	MerchantID   int32            `db:"merchant_id" json:"merchant_id"`
	UserID       int32            `db:"user_id" json:"user_id"`
	Name         string           `db:"name" json:"name"`
	Description  *string          `db:"description" json:"description"`
	Address      *string          `db:"address" json:"address"`
	ContactEmail *string          `db:"contact_email" json:"contact_email"`
	ContactPhone *string          `db:"contact_phone" json:"contact_phone"`
	Status       string           `db:"status" json:"status"`
	CreatedAt    pgtype.Timestamp `db:"created_at" json:"created_at"`
	UpdatedAt    pgtype.Timestamp `db:"updated_at" json:"updated_at"`
	DeletedAt    pgtype.Timestamp `db:"deleted_at" json:"deleted_at"`
	TotalCount   int64            `db:"total_count" json:"total_count"`
}

type GetMerchantsActiveParams struct {
	Column1 string `db:"column1" json:"column_1"`
	Limit   int32  `db:"limit" json:"limit"`
	Offset  int32  `db:"offset" json:"offset"`
}

type GetMerchantsActiveRow struct {
	MerchantID   int32            `db:"merchant_id" json:"merchant_id"`
	UserID       int32            `db:"user_id" json:"user_id"`
	Name         string           `db:"name" json:"name"`
	Description  *string          `db:"description" json:"description"`
	Address      *string          `db:"address" json:"address"`
	ContactEmail *string          `db:"contact_email" json:"contact_email"`
	ContactPhone *string          `db:"contact_phone" json:"contact_phone"`
	Status       string           `db:"status" json:"status"`
	CreatedAt    pgtype.Timestamp `db:"created_at" json:"created_at"`
	UpdatedAt    pgtype.Timestamp `db:"updated_at" json:"updated_at"`
	DeletedAt    pgtype.Timestamp `db:"deleted_at" json:"deleted_at"`
	TotalCount   int64            `db:"total_count" json:"total_count"`
}

type GetMerchantsTrashedParams struct {
	Column1 string `db:"column1" json:"column_1"`
	Limit   int32  `db:"limit" json:"limit"`
	Offset  int32  `db:"offset" json:"offset"`
}

type GetMerchantsTrashedRow struct {
	MerchantID   int32            `db:"merchant_id" json:"merchant_id"`
	UserID       int32            `db:"user_id" json:"user_id"`
	Name         string           `db:"name" json:"name"`
	Description  *string          `db:"description" json:"description"`
	Address      *string          `db:"address" json:"address"`
	ContactEmail *string          `db:"contact_email" json:"contact_email"`
	ContactPhone *string          `db:"contact_phone" json:"contact_phone"`
	Status       string           `db:"status" json:"status"`
	CreatedAt    pgtype.Timestamp `db:"created_at" json:"created_at"`
	UpdatedAt    pgtype.Timestamp `db:"updated_at" json:"updated_at"`
	DeletedAt    pgtype.Timestamp `db:"deleted_at" json:"deleted_at"`
	TotalCount   int64            `db:"total_count" json:"total_count"`
}

type UpdateMerchantParams struct {
	MerchantID   int32   `db:"merchant_id" json:"merchant_id"`
	Name         string  `db:"name" json:"name"`
	Description  *string `db:"description" json:"description"`
	Address      *string `db:"address" json:"address"`
	ContactEmail *string `db:"contact_email" json:"contact_email"`
	ContactPhone *string `db:"contact_phone" json:"contact_phone"`
	Status       string  `db:"status" json:"status"`
}

type UpdateMerchantStatusParams struct {
	MerchantID int32  `db:"merchant_id" json:"merchant_id"`
	Status     string `db:"status" json:"status"`
}
