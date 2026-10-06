package cashier_errors

import (
	"net/http"

	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiCashierInvalidYear = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid year", http.StatusBadRequest)
	}

	ErrApiCashierInvalidMonth = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid month", http.StatusBadRequest)
	}

	ErrApiCashierNotFound = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "cashier not found", http.StatusNotFound)
	}
	ErrApiCashierInvalidId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid cashier id", http.StatusBadRequest)
	}

	ErrApiCashierInvalidMerchantId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid merchant id", http.StatusBadRequest)
	}

	ErrApiCashierFailedFindAll = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find all cashiers", http.StatusInternalServerError)
	}
	ErrApiCashierFailedFindById = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find cashier by id", http.StatusInternalServerError)
	}
	ErrApiCashierFailedFindByActive = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find active cashiers", http.StatusInternalServerError)
	}
	ErrApiCashierFailedFindByTrashed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find trashed cashiers", http.StatusInternalServerError)
	}

	ErrApiCashierFailedMonthlyTotalSales = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find monthly total sales", http.StatusInternalServerError)
	}
	ErrApiCashierFailedYearlyTotalSales = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find yearly total sales", http.StatusInternalServerError)
	}
	ErrApiCashierFailedMonthlyTotalSalesByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find monthly total sales by merchant", http.StatusInternalServerError)
	}
	ErrApiCashierFailedYearlyTotalSalesByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find yearly total sales by merchant", http.StatusInternalServerError)
	}
	ErrApiCashierFailedMonthlyTotalSalesById = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find monthly total sales by cashier id", http.StatusInternalServerError)
	}
	ErrApiCashierFailedYearlyTotalSalesById = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find yearly total sales by cashier id", http.StatusInternalServerError)
	}

	ErrApiCashierFailedMonthSales = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find monthly sales", http.StatusInternalServerError)
	}
	ErrApiCashierFailedYearSales = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find yearly sales", http.StatusInternalServerError)
	}
	ErrApiCashierFailedMonthSalesByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find monthly sales by merchant", http.StatusInternalServerError)
	}
	ErrApiCashierFailedYearSalesByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find yearly sales by merchant", http.StatusInternalServerError)
	}
	ErrApiCashierFailedMonthSalesById = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find monthly sales by cashier id", http.StatusInternalServerError)
	}
	ErrApiCashierFailedYearSalesById = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find yearly sales by cashier id", http.StatusInternalServerError)
	}

	ErrApiCashierFailedCreate = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to create cashier", http.StatusInternalServerError)
	}
	ErrApiCashierFailedUpdate = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to update cashier", http.StatusInternalServerError)
	}

	ErrApiValidateCreateCashier = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid create cashier request", http.StatusBadRequest)
	}

	ErrApiValidateUpdateCashier = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid update cashier request", http.StatusBadRequest)
	}

	ErrApiBindCreateCashier = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid create cashier request", http.StatusBadRequest)
	}

	ErrApiBindUpdateCashier = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid update cashier request", http.StatusBadRequest)
	}

	ErrApiCashierFailedTrashed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to trash cashier", http.StatusInternalServerError)
	}
	ErrApiCashierFailedRestore = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore cashier", http.StatusInternalServerError)
	}
	ErrApiCashierFailedDeletePermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete cashier", http.StatusInternalServerError)
	}
	ErrApiCashierFailedRestoreAll = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore all cashiers", http.StatusInternalServerError)
	}
	ErrApiCashierFailedDeleteAllPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete all cashiers", http.StatusInternalServerError)
	}
)
