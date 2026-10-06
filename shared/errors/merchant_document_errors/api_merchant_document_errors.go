package merchantdocument_errors

import (
	"net/http"

	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"
	"github.com/labstack/echo/v4"
)

var (
	ErrApiInvalidMerchantDocumentID = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Invalid merchant document ID", http.StatusBadRequest)
	}

	ErrApiFailedFindAllMerchantDocuments = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch all merchant documents", http.StatusInternalServerError)
	}

	ErrApiFailedFindByIdMerchantDocument = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch merchant document by ID", http.StatusInternalServerError)
	}

	ErrApiFailedFindAllActiveMerchantDocuments = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch all active merchant documents", http.StatusInternalServerError)
	}

	ErrApiFailedFindAllTrashedMerchantDocuments = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to fetch all trashed merchant documents", http.StatusInternalServerError)
	}

	ErrApiFailedCreateMerchantDocument = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to create merchant document", http.StatusInternalServerError)
	}

	ErrApiFailedUpdateMerchantDocument = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to update merchant document", http.StatusInternalServerError)
	}

	ErrApiFailedUpdateMerchantDocumentStatus = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to update merchant document status", http.StatusInternalServerError)
	}

	ErrApiFailedTrashMerchantDocument = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to trash merchant document", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreMerchantDocument = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to restore merchant document", http.StatusInternalServerError)
	}

	ErrApiFailedDeleteMerchantDocumentPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to permanently delete merchant document", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreAllMerchantDocuments = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to restore all merchant documents", http.StatusInternalServerError)
	}

	ErrApiFailedDeleteAllMerchantDocumentsPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Failed to permanently delete all merchant documents", http.StatusInternalServerError)
	}

	ErrApiValidateCreateMerchantDocument = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid create merchant document request", http.StatusBadRequest)
	}

	ErrApiValidateUpdateMerchantDocument = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid update merchant document request", http.StatusBadRequest)
	}

	ErrApiValidateUpdateMerchantDocumentStatus = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid update merchant document status request", http.StatusBadRequest)
	}

	ErrApiBindCreateMerchantDocument = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid create merchant document request", http.StatusBadRequest)
	}

	ErrApiBindUpdateMerchantDocument = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid update merchant document request", http.StatusBadRequest)
	}

	ErrApiBindUpdateMerchantDocumentStatus = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid update merchant document status request", http.StatusBadRequest)
	}
)
