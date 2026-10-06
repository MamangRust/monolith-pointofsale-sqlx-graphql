package repository

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-pointofsale-shared/errors"
	"github.com/jmoiron/sqlx"
)

type merchantDocumentQueryRepository struct {
	db *sqlx.DB
}

func NewMerchantDocumentQueryRepository(db *sqlx.DB) MerchantDocumentQueryRepository {
	return &merchantDocumentQueryRepository{
		db: db,
	}
}

const getMerchantDocuments = `SELECT *, COUNT(*) OVER() AS total_count
FROM merchant_documents
WHERE deleted_at IS NULL
  AND (
    $1::TEXT IS NULL OR
    document_type ILIKE '%' || $1 || '%' OR
    status ILIKE '%' || $1 || '%' OR
    note ILIKE '%' || $1 || '%'
  )
ORDER BY document_id
LIMIT $2 OFFSET $3`

const getActiveMerchantDocuments = `SELECT *, COUNT(*) OVER() AS total_count
FROM merchant_documents
WHERE deleted_at IS NULL AND status != 'deleted'
  AND (
    $1::TEXT IS NULL OR
    document_type ILIKE '%' || $1 || '%' OR
    status ILIKE '%' || $1 || '%' OR
    note ILIKE '%' || $1 || '%'
  )
ORDER BY document_id
LIMIT $2 OFFSET $3`

const getTrashedMerchantDocuments = `SELECT *, COUNT(*) OVER() AS total_count
FROM merchant_documents
WHERE deleted_at IS NOT NULL
  AND (
    $1::TEXT IS NULL OR
    document_type ILIKE '%' || $1 || '%' OR
    status ILIKE '%' || $1 || '%' OR
    note ILIKE '%' || $1 || '%'
  )
ORDER BY document_id
LIMIT $2 OFFSET $3`

const getMerchantDocument = `SELECT *
FROM merchant_documents
WHERE document_id = $1 AND deleted_at IS NULL`

func (r *merchantDocumentQueryRepository) FindAllDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetMerchantDocumentsRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var docs []*db.GetMerchantDocumentsRow
	err := r.db.SelectContext(ctx, &docs, getMerchantDocuments,
		req.Search,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	totalCount := 0
	if len(docs) > 0 {
		totalCount = int(docs[0].TotalCount)
	}

	return docs, &totalCount, nil
}

func (r *merchantDocumentQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetActiveMerchantDocumentsRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var docs []*db.GetActiveMerchantDocumentsRow
	err := r.db.SelectContext(ctx, &docs, getActiveMerchantDocuments,
		req.Search,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	totalCount := 0
	if len(docs) > 0 {
		totalCount = int(docs[0].TotalCount)
	}

	return docs, &totalCount, nil
}

func (r *merchantDocumentQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetTrashedMerchantDocumentsRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var docs []*db.GetTrashedMerchantDocumentsRow
	err := r.db.SelectContext(ctx, &docs, getTrashedMerchantDocuments,
		req.Search,
		int32(req.PageSize),
		int32(offset),
	)
	if err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	totalCount := 0
	if len(docs) > 0 {
		totalCount = int(docs[0].TotalCount)
	}

	return docs, &totalCount, nil
}

func (r *merchantDocumentQueryRepository) FindById(ctx context.Context, id int) (*db.MerchantDocument, error) {
	var doc db.MerchantDocument
	if err := r.db.GetContext(ctx, &doc, getMerchantDocument, int32(id)); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &doc, nil
}
