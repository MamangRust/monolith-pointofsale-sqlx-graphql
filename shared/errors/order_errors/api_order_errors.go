package order_errors

import (
	"net/http"

	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiOrderInvalidYear = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid year", http.StatusBadRequest)
	}

	ErrApiOrderInvalidMonth = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid month", http.StatusBadRequest)
	}

	ErrApiOrderInvalidMerchantId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid merchant id", http.StatusBadRequest)
	}

	ErrApiOrderFailedFindAll = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find all orders", http.StatusInternalServerError)
	}
	ErrApiOrderFailedFindById = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find order by ID", http.StatusInternalServerError)
	}
	ErrApiOrderFailedFindByActive = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find active orders", http.StatusInternalServerError)
	}
	ErrApiOrderFailedFindByTrashed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find trashed orders", http.StatusInternalServerError)
	}

	ErrApiOrderFailedFindMonthlyTotalRevenue = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find monthly total revenue", http.StatusInternalServerError)
	}
	ErrApiOrderFailedFindYearlyTotalRevenue = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find yearly total revenue", http.StatusInternalServerError)
	}
	ErrApiOrderFailedFindMonthlyTotalRevenueByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find monthly total revenue by merchant", http.StatusInternalServerError)
	}
	ErrApiOrderFailedFindYearlyTotalRevenueByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find yearly total revenue by merchant", http.StatusInternalServerError)
	}

	ErrApiOrderFailedFindMonthlyRevenue = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find monthly revenue", http.StatusInternalServerError)
	}
	ErrApiOrderFailedFindYearlyRevenue = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find yearly revenue", http.StatusInternalServerError)
	}
	ErrApiOrderFailedFindMonthlyRevenueByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find monthly revenue by merchant", http.StatusInternalServerError)
	}
	ErrApiOrderFailedFindYearlyRevenueByMerchant = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to find yearly revenue by merchant", http.StatusInternalServerError)
	}

	ErrApiOrderFailedCreate = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to create order", http.StatusInternalServerError)
	}
	ErrApiOrderFailedUpdate = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to update order", http.StatusInternalServerError)
	}

	ErrApiValidateCreateOrder = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid create bank request", http.StatusBadRequest)
	}

	ErrApiValidateUpdateOrder = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid update bank request", http.StatusBadRequest)
	}

	ErrApiBindCreateOrder = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid create bank request", http.StatusBadRequest)
	}

	ErrApiBindUpdateOrder = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid update bank request", http.StatusBadRequest)
	}

	ErrApiOrderFailedTrashed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to trashed order", http.StatusInternalServerError)
	}
	ErrApiOrderFailedRestore = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore order", http.StatusInternalServerError)
	}
	ErrApiOrderFailedDeletePermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete order", http.StatusInternalServerError)
	}
	ErrApiOrderFailedRestoreAll = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore all orders", http.StatusInternalServerError)
	}
	ErrApiOrderFailedDeleteAllPermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to permanently delete all orders", http.StatusInternalServerError)
	}

	ErrApiOrderNotFound = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "order not found", http.StatusNotFound)
	}
	ErrApiOrderInvalidId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid order ID", http.StatusBadRequest)
	}
)
