package merchant_errors

import (
	"net/http"

	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiMerchantNotFound = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "merchant not found", http.StatusNotFound)
	}
	ErrApiMerchantInvalidId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid merchant id", http.StatusBadRequest)
	}

	ErrApiMerchantFailedFindAll = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find all merchants", http.StatusInternalServerError)
	}
	ErrApiMerchantFailedFindById = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find merchant by id", http.StatusInternalServerError)
	}
	ErrApiMerchantFailedFindByActive = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find active merchants", http.StatusInternalServerError)
	}
	ErrApiMerchantFailedFindByTrashed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find trashed merchants", http.StatusInternalServerError)
	}

	ErrApiMerchantFailedCreate = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to create merchant", http.StatusInternalServerError)
	}
	ErrApiMerchantFailedUpdate = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to update merchant", http.StatusInternalServerError)
	}

	ErrApiMerchantFailedUpdateStatus = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to update merchant status", http.StatusInternalServerError)
	}

	ErrApiValidateCreateMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid create bank request", http.StatusBadRequest)
	}

	ErrApiValidateUpdateMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid update bank request", http.StatusBadRequest)
	}

	ErrApiValidateUpdateMerchantStatus = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid update merchant status request", http.StatusBadRequest)
	}

	ErrApiBindUpdateMerchantStatus = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid update merchant status request", http.StatusBadRequest)
	}

	ErrApiBindCreateMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid create bank request", http.StatusBadRequest)
	}

	ErrApiBindUpdateMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid update bank request", http.StatusBadRequest)
	}

	ErrApiMerchantFailedTrashed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to trash merchant", http.StatusInternalServerError)
	}
	ErrApiMerchantFailedRestore = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore merchant", http.StatusInternalServerError)
	}
	ErrApiMerchantFailedDeletePermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete merchant", http.StatusInternalServerError)
	}
	ErrApiMerchantFailedRestoreAll = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore all merchants", http.StatusInternalServerError)
	}
	ErrApiMerchantFailedDeleteAll = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete all merchants", http.StatusInternalServerError)
	}
)
