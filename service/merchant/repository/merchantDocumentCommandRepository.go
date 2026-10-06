package repository

import (
	"context"

	db "github.com/MamangRust/monolith-pointofsale-pkg/database/schema"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-pointofsale-shared/errors"
	"github.com/jmoiron/sqlx"
)

type merchantDocumentCommandRepository struct {
	db *sqlx.DB
}

func NewMerchantDocumentCommandRepository(db *sqlx.DB) MerchantDocumentCommandRepository {
	return &merchantDocumentCommandRepository{
		db: db,
	}
}

const createMerchantDocument = `INSERT INTO merchant_documents (
    merchant_id, document_type, document_url, status, note, uploaded_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, current_timestamp, current_timestamp
)
RETURNING *`

const updateMerchantDocument = `UPDATE merchant_documents
SET
    document_type = $2,
    document_url = $3,
    status = $4,
    note = $5,
    updated_at = current_timestamp
WHERE
    document_id = $1 AND deleted_at IS NULL
RETURNING *`

const updateMerchantDocumentStatus = `UPDATE merchant_documents
SET
    status = $2,
    note = $3,
    updated_at = current_timestamp
WHERE
    document_id = $1 AND deleted_at IS NULL
RETURNING *`

const trashMerchantDocument = `UPDATE merchant_documents
SET
    deleted_at = current_timestamp
WHERE
    document_id = $1 AND deleted_at IS NULL
RETURNING *`

const restoreMerchantDocument = `UPDATE merchant_documents
SET
    deleted_at = NULL
WHERE
    document_id = $1 AND deleted_at IS NOT NULL
RETURNING *`

const deleteMerchantDocumentPermanently = `DELETE FROM merchant_documents
WHERE
    document_id = $1 AND deleted_at IS NOT NULL`

const restoreAllMerchantDocuments = `UPDATE merchant_documents
SET deleted_at = NULL
WHERE deleted_at IS NOT NULL`

const deleteAllPermanentMerchantDocuments = `DELETE FROM merchant_documents
WHERE deleted_at IS NOT NULL`

func (r *merchantDocumentCommandRepository) CreateMerchantDocument(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*db.MerchantDocument, error) {
	return r.createMerchantDocument(ctx, r.db, request)
}

// CreateMerchantDocumentInTx persists the document inside the given database
// transaction so the caller can commit the business write and its outbox event
// atomically (Phase 6 — transactional outbox).
func (r *merchantDocumentCommandRepository) CreateMerchantDocumentInTx(ctx context.Context, tx *sqlx.Tx, request *requests.CreateMerchantDocumentRequest) (*db.MerchantDocument, error) {
	return r.createMerchantDocument(ctx, tx, request)
}

func (r *merchantDocumentCommandRepository) createMerchantDocument(ctx context.Context, q sqlx.QueryerContext, request *requests.CreateMerchantDocumentRequest) (*db.MerchantDocument, error) {
	emptyNote := ""

	var doc db.MerchantDocument
	if err := sqlx.GetContext(ctx, q, &doc, createMerchantDocument,
		int32(request.MerchantID),
		request.DocumentType,
		request.DocumentUrl,
		"pending",
		&emptyNote,
	); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return &doc, nil
}

// UpdateMerchantDocumentStatusInTx updates the document status inside the given
// database transaction so the caller can commit the business write and its
// outbox event atomically (Phase 6 — transactional outbox).
func (r *merchantDocumentCommandRepository) UpdateMerchantDocumentStatusInTx(ctx context.Context, tx *sqlx.Tx, request *requests.UpdateMerchantDocumentStatusRequest) (*db.MerchantDocument, error) {
	return r.updateMerchantDocumentStatus(ctx, tx, request)
}

func (r *merchantDocumentCommandRepository) UpdateMerchantDocumentStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*db.MerchantDocument, error) {
	return r.updateMerchantDocumentStatus(ctx, r.db, request)
}

func (r *merchantDocumentCommandRepository) updateMerchantDocumentStatus(ctx context.Context, q sqlx.QueryerContext, request *requests.UpdateMerchantDocumentStatusRequest) (*db.MerchantDocument, error) {
	if request.DocumentID == nil {
		return nil, sharedErrors.ErrBadRequest.WithMessage("document id is required")
	}

	var doc db.MerchantDocument
	if err := sqlx.GetContext(ctx, q, &doc, updateMerchantDocumentStatus,
		int32(*request.DocumentID),
		request.Status,
		&request.Note,
	); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return &doc, nil
}

func (r *merchantDocumentCommandRepository) UpdateMerchantDocument(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*db.MerchantDocument, error) {
	if request.DocumentID == nil {
		return nil, sharedErrors.ErrBadRequest.WithMessage("document id is required")
	}

	var doc db.MerchantDocument
	if err := r.db.GetContext(ctx, &doc, updateMerchantDocument,
		int32(*request.DocumentID),
		request.DocumentType,
		request.DocumentUrl,
		request.Status,
		&request.Note,
	); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return &doc, nil
}

func (r *merchantDocumentCommandRepository) TrashedMerchantDocument(ctx context.Context, documentID int) (*db.MerchantDocument, error) {
	var doc db.MerchantDocument
	if err := r.db.GetContext(ctx, &doc, trashMerchantDocument, int32(documentID)); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return &doc, nil
}

func (r *merchantDocumentCommandRepository) RestoreMerchantDocument(ctx context.Context, documentID int) (*db.MerchantDocument, error) {
	var doc db.MerchantDocument
	if err := r.db.GetContext(ctx, &doc, restoreMerchantDocument, int32(documentID)); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return &doc, nil
}

func (r *merchantDocumentCommandRepository) DeleteMerchantDocumentPermanent(ctx context.Context, documentID int) (bool, error) {
	if _, err := r.db.ExecContext(ctx, deleteMerchantDocumentPermanently, int32(documentID)); err != nil {
		return false, sharedErrors.ErrInternal.WithInternal(err)
	}

	return true, nil
}

func (r *merchantDocumentCommandRepository) RestoreAllMerchantDocument(ctx context.Context) (bool, error) {
	if _, err := r.db.ExecContext(ctx, restoreAllMerchantDocuments); err != nil {
		return false, sharedErrors.ErrInternal.WithInternal(err)
	}

	return true, nil
}

func (r *merchantDocumentCommandRepository) DeleteAllMerchantDocumentPermanent(ctx context.Context) (bool, error) {
	if _, err := r.db.ExecContext(ctx, deleteAllPermanentMerchantDocuments); err != nil {
		return false, sharedErrors.ErrInternal.WithInternal(err)
	}

	return true, nil
}
