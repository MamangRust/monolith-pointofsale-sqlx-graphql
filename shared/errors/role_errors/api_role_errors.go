package role_errors

import (
	"net/http"

	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiRoleNotFound = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "Role not found", http.StatusNotFound)
	}

	ErrApiRoleInvalidId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid Role id", http.StatusNotFound)
	}

	ErrApiFailedFindAll = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch Roles", http.StatusInternalServerError)
	}

	ErrApiFailedFindActive = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch active Roles", http.StatusInternalServerError)
	}

	ErrApiFailedFindTrashed = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to fetch trashed Roles", http.StatusInternalServerError)
	}

	ErrApiFailedCreateRole = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to create Role", http.StatusInternalServerError)
	}

	ErrApiFailedUpdateRole = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to update Role", http.StatusInternalServerError)
	}

	ErrApiValidateCreateRole = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid create Role request", http.StatusBadRequest)
	}

	ErrApiValidateUpdateRole = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid update Role request", http.StatusBadRequest)
	}

	ErrInvalidRoleId = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid Role id", http.StatusBadRequest)
	}

	ErrApiBindCreateRole = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid create Role request", http.StatusBadRequest)
	}

	ErrApiBindUpdateRole = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "bind failed: invalid update Role request", http.StatusBadRequest)
	}

	ErrApiFailedTrashedRole = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to move Role to trash", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreRole = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore Role", http.StatusInternalServerError)
	}

	ErrApiFailedDeletePermanent = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to delete Role permanently", http.StatusInternalServerError)
	}

	ErrApiFailedRestoreAll = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to restore all Roles", http.StatusInternalServerError)
	}

	ErrApiFailedDeleteAll = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to delete all Roles permanently", http.StatusInternalServerError)
	}
)
