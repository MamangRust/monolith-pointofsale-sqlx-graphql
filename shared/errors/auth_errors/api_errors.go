package auth_errors

import (
	"net/http"

	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"

	"github.com/labstack/echo/v4"
)

var (
	ErrApiInvalidVerifyCode = func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid or missing verify_code")
	}

	ErrApiVerifyCode = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to verify code", http.StatusBadRequest)
	}

	ErrApiForgotPassword = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to forgot password", http.StatusBadRequest)
	}

	ErrApiResetPassword = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "failed to reset password", http.StatusBadRequest)
	}

	ErrApiLogin = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "login failed: invalid argument provided", http.StatusBadRequest)
	}

	ErrApiRefreshToken = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "refresh-token failed: invalid access token", http.StatusBadRequest)
	}

	ErrApiGetMe = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "get user info failed: unauthenticated", http.StatusUnauthorized)
	}

	ErrValidateLogin = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid login request", http.StatusBadRequest)
	}

	ErrValidateRegister = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid register request", http.StatusBadRequest)
	}

	ErrValidateRefreshToken = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid refresh-token request", http.StatusBadRequest)
	}

	ErrValidateForgotPassword = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid forgot-password request", http.StatusBadRequest)
	}

	ErrValidateResetPassword = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "validation failed: invalid reset-password request", http.StatusBadRequest)
	}

	ErrBindForgotPassword = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "binding failed: invalid forgot password request payload", http.StatusBadRequest)
	}

	ErrBindResetPassword = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "binding failed: invalid reset password request payload", http.StatusBadRequest)
	}

	ErrBindLogin = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "binding failed: invalid login request payload", http.StatusBadRequest)
	}

	ErrBindRefreshToken = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "binding failed: invalid refresh token request payload", http.StatusBadRequest)
	}

	ErrBindRegister = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "binding failed: invalid register request payload", http.StatusBadRequest)
	}

	ErrInvalidLogin = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid email or password", http.StatusUnauthorized)
	}

	ErrInvalidAccessToken = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "invalid access token", http.StatusInternalServerError)
	}

	ErrApiRegister = func(c echo.Context) error {
		return response.NewApiErrorResponse(c, "error", "register failed: invalid argument", http.StatusBadRequest)
	}
)
