package transaction_errors

import (
	"net/http"

	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiTransactionInvalidYear = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid year", http.StatusBadRequest)
	}
	ErrApiTransactionInvalidMonth = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid month", http.StatusBadRequest)
	}

	ErrApiTransactionInvalidMerchantId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid merchant ID", http.StatusBadRequest)
	}

	ErrApiTransactionFailedFindAll = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find all transactions", http.StatusInternalServerError)
	}
	ErrApiTransactionFailedFindById = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find transaction by ID", http.StatusInternalServerError)
	}
	ErrApiTransactionFailedFindByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find transaction by merchant", http.StatusInternalServerError)
	}
	ErrApiTransactionFailedFindByActive = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find active transactions", http.StatusInternalServerError)
	}
	ErrApiTransactionFailedFindByTrashed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find trashed transactions", http.StatusInternalServerError)
	}

	ErrApiTransactionFailedFindMonthSuccess = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find monthly successful transactions", http.StatusInternalServerError)
	}
	ErrApiTransactionFailedFindYearSuccess = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find yearly successful transactions", http.StatusInternalServerError)
	}
	ErrApiTransactionFailedFindMonthFailed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find monthly failed transactions", http.StatusInternalServerError)
	}
	ErrApiTransactionFailedFindYearFailed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find yearly failed transactions", http.StatusInternalServerError)
	}

	ErrApiTransactionFailedFindMonthSuccessByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find monthly successful transactions by merchant", http.StatusInternalServerError)
	}
	ErrApiTransactionFailedFindYearSuccessByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find yearly successful transactions by merchant", http.StatusInternalServerError)
	}
	ErrApiTransactionFailedFindMonthFailedByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find monthly failed transactions by merchant", http.StatusInternalServerError)
	}
	ErrApiTransactionFailedFindYearFailedByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find yearly failed transactions by merchant", http.StatusInternalServerError)
	}

	ErrApiTransactionSuccessFindMonthMethod = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find monthly successful transaction methods", http.StatusInternalServerError)
	}

	ErrApiTransactionSuccessFindYearMethod = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find yearly successful transaction methods", http.StatusInternalServerError)
	}

	ErrApiTransactionSuccessFindMonthMethodByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find monthly successful transaction methods by merchant", http.StatusInternalServerError)
	}

	ErrApiTransactionSuccessFindYearMethodByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find yearly successful transaction methods by merchant", http.StatusInternalServerError)
	}

	ErrApiTransactionFailedFindMonthMethod = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find monthly failed transaction methods", http.StatusInternalServerError)
	}

	ErrApiTransactionFailedFindYearMethod = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find yearly failed transaction methods", http.StatusInternalServerError)
	}

	ErrApiTransactionFailedFindMonthMethodByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find monthly failed transaction methods by merchant", http.StatusInternalServerError)
	}

	ErrApiTransactionFailedFindYearMethodByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find yearly failed transaction methods by merchant", http.StatusInternalServerError)
	}

	ErrApiTransactionFailedCreate = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to create transaction", http.StatusInternalServerError)
	}
	ErrApiTransactionFailedUpdate = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to update transaction", http.StatusInternalServerError)
	}

	ErrApiValidateCreateTransaction = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid create bank request", http.StatusBadRequest)
	}

	ErrApiValidateUpdateTransaction = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid update bank request", http.StatusBadRequest)
	}

	ErrApiBindCreateTransaction = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid create bank request", http.StatusBadRequest)
	}

	ErrApiBindUpdateTransaction = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid update bank request", http.StatusBadRequest)
	}

	ErrApiTransactionFailedTrashed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to trashed transaction", http.StatusInternalServerError)
	}
	ErrApiTransactionFailedRestore = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore transaction", http.StatusInternalServerError)
	}
	ErrApiTransactionFailedDeletePermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete transaction", http.StatusInternalServerError)
	}
	ErrApiTransactionFailedRestoreAll = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore all transactions", http.StatusInternalServerError)
	}
	ErrApiTransactionFailedDeleteAllPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete all transactions", http.StatusInternalServerError)
	}

	ErrApiTransactionNotFound = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "transaction not found", http.StatusNotFound)
	}
	ErrApiTransactionInvalidId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid transaction ID", http.StatusBadRequest)
	}
)
