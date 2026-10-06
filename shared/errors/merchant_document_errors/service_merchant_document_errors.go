package merchantdocument_errors

import (
	"github.com/MamangRust/monolith-pointofsale-shared/errors"
)

var (
	ErrMerchantDocumentNotFoundRes        = errors.ErrNotFound.WithMessage("Merchant Document not found")
	ErrFailedFindAllMerchantDocuments     = errors.ErrInternal.WithMessage("Failed to fetch Merchant Documents")
	ErrFailedFindActiveMerchantDocuments  = errors.ErrInternal.WithMessage("Failed to fetch active Merchant Documents")
	ErrFailedFindTrashedMerchantDocuments = errors.ErrInternal.WithMessage("Failed to fetch trashed Merchant Documents")
	ErrFailedFindMerchantDocumentById     = errors.ErrInternal.WithMessage("Failed to find Merchant Document by ID")

	ErrFailedCreateMerchantDocument = errors.ErrInternal.WithMessage("Failed to create Merchant Document")
	ErrFailedUpdateMerchantDocument = errors.ErrInternal.WithMessage("Failed to update Merchant Document")

	ErrFailedTrashMerchantDocument   = errors.ErrInternal.WithMessage("Failed to trash Merchant Document")
	ErrFailedRestoreMerchantDocument = errors.ErrInternal.WithMessage("Failed to restore Merchant Document")
	ErrFailedDeleteMerchantDocument  = errors.ErrInternal.WithMessage("Failed to delete Merchant Document permanently")

	ErrFailedRestoreAllMerchantDocuments = errors.ErrInternal.WithMessage("Failed to restore all Merchant Documents")
	ErrFailedDeleteAllMerchantDocuments  = errors.ErrInternal.WithMessage("Failed to delete all Merchant Documents permanently")
)
