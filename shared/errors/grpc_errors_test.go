package errors

import (
	"net/http"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHTTPToGrpcCode(t *testing.T) {
	cases := []struct {
		name string
		code int
		want codes.Code
	}{
		{"bad request", http.StatusBadRequest, codes.InvalidArgument},
		{"unauthorized", http.StatusUnauthorized, codes.Unauthenticated},
		{"forbidden", http.StatusForbidden, codes.PermissionDenied},
		{"not found", http.StatusNotFound, codes.NotFound},
		{"conflict", http.StatusConflict, codes.AlreadyExists},
		{"too many requests", http.StatusTooManyRequests, codes.ResourceExhausted},
		{"gateway timeout", http.StatusGatewayTimeout, codes.DeadlineExceeded},
		{"service unavailable", http.StatusServiceUnavailable, codes.Unavailable},
		{"unmapped falls back to internal", http.StatusNotAcceptable, codes.Internal},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := httpToGrpcCode(c.code); got != c.want {
				t.Errorf("httpToGrpcCode(%d) = %v, want %v", c.code, got, c.want)
			}
		})
	}
}

func TestGrpcToHTTPCode(t *testing.T) {
	cases := []struct {
		name string
		code codes.Code
		want int
	}{
		{"invalid argument", codes.InvalidArgument, http.StatusBadRequest},
		{"unauthenticated", codes.Unauthenticated, http.StatusUnauthorized},
		{"permission denied", codes.PermissionDenied, http.StatusForbidden},
		{"not found", codes.NotFound, http.StatusNotFound},
		{"already exists", codes.AlreadyExists, http.StatusConflict},
		{"resource exhausted", codes.ResourceExhausted, http.StatusTooManyRequests},
		{"deadline exceeded", codes.DeadlineExceeded, http.StatusGatewayTimeout},
		{"unavailable", codes.Unavailable, http.StatusServiceUnavailable},
		{"unmapped falls back to internal", codes.Aborted, http.StatusInternalServerError},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := grpcToHttpCode(c.code); got != c.want {
				t.Errorf("grpcToHttpCode(%v) = %d, want %d", c.code, got, c.want)
			}
		})
	}
}

func TestGrpcToErrorType(t *testing.T) {
	cases := []struct {
		name string
		code codes.Code
		want ErrorType
	}{
		{"not found", codes.NotFound, ErrorTypeNotFound},
		{"invalid argument", codes.InvalidArgument, ErrorTypeBadRequest},
		{"already exists", codes.AlreadyExists, ErrorTypeConflict},
		{"permission denied", codes.PermissionDenied, ErrorTypeForbidden},
		{"unauthenticated", codes.Unauthenticated, ErrorTypeUnauthorized},
		{"deadline exceeded", codes.DeadlineExceeded, ErrorTypeTimeout},
		{"resource exhausted", codes.ResourceExhausted, ErrorTypeBadRequest},
		{"unavailable", codes.Unavailable, ErrorTypeUnavailable},
		{"unmapped falls back to internal", codes.Aborted, ErrorTypeInternal},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := grpcToErrorType(c.code); got != c.want {
				t.Errorf("grpcToErrorType(%v) = %v, want %v", c.code, got, c.want)
			}
		})
	}
}

func TestToGrpcErrorServiceUnavailableRoundTrip(t *testing.T) {
	// ErrServiceUnavailable (503) must map to codes.Unavailable and back to
	// 503 with ErrorTypeUnavailable (gap G-01 two-way mapping).
	appErr := ErrServiceUnavailable.WithMessage("order service is temporarily unavailable")

	grpcErr := ToGrpcError(appErr)
	st, ok := status.FromError(grpcErr)
	if !ok {
		t.Fatal("ToGrpcError did not produce a gRPC status")
	}
	if st.Code() != codes.Unavailable {
		t.Errorf("gRPC code = %v, want %v", st.Code(), codes.Unavailable)
	}

	parsed := ParseGrpcError(grpcErr)
	if parsed == nil {
		t.Fatal("ParseGrpcError returned nil")
	}
	if parsed.Code != http.StatusServiceUnavailable {
		t.Errorf("parsed HTTP code = %d, want %d", parsed.Code, http.StatusServiceUnavailable)
	}
	if parsed.Type != ErrorTypeUnavailable {
		t.Errorf("parsed type = %v, want %v", parsed.Type, ErrorTypeUnavailable)
	}
}

func TestToGrpcErrorBadRequestRoundTrip(t *testing.T) {
	appErr := NewBadRequestError("order item is empty")

	grpcErr := ToGrpcError(appErr)
	st, ok := status.FromError(grpcErr)
	if !ok {
		t.Fatal("ToGrpcError did not produce a gRPC status")
	}
	if st.Code() != codes.InvalidArgument {
		t.Errorf("gRPC code = %v, want %v", st.Code(), codes.InvalidArgument)
	}

	parsed := ParseGrpcError(grpcErr)
	if parsed == nil {
		t.Fatal("ParseGrpcError returned nil")
	}
	if parsed.Code != http.StatusBadRequest {
		t.Errorf("parsed HTTP code = %d, want %d", parsed.Code, http.StatusBadRequest)
	}
	if parsed.Type != ErrorTypeBadRequest {
		t.Errorf("parsed type = %v, want %v", parsed.Type, ErrorTypeBadRequest)
	}
}

func TestParseGrpcErrorInvalidDetailCodeFallsBackToStatusCode(t *testing.T) {
	grpcErr := status.New(codes.Unavailable, "dependency unavailable").Err()
	parsed := ParseGrpcError(grpcErr)
	if parsed == nil {
		t.Fatal("ParseGrpcError returned nil")
	}
	if parsed.Code != http.StatusServiceUnavailable {
		t.Errorf("parsed HTTP code = %d, want %d", parsed.Code, http.StatusServiceUnavailable)
	}
	if parsed.Type != ErrorTypeUnavailable {
		t.Errorf("parsed type = %v, want %v", parsed.Type, ErrorTypeUnavailable)
	}
}

func TestParseGrpcErrorZeroDetailCodeDoesNotBecomeSuccess(t *testing.T) {
	grpcErr := NewGrpcError("bad detail", http.StatusOK)
	parsed := ParseGrpcError(grpcErr)
	if parsed == nil {
		t.Fatal("ParseGrpcError returned nil")
	}
	if parsed.Code == http.StatusOK || parsed.Code < 100 || parsed.Code > 599 {
		t.Errorf("parsed HTTP code = %d, want non-success valid error code", parsed.Code)
	}
}

func TestToGrpcErrorNonAppErrorFallsBackToInternal(t *testing.T) {
	rawErr := NewStdlibError("no rows in result set")

	grpcErr := ToGrpcError(rawErr)
	st, ok := status.FromError(grpcErr)
	if !ok {
		t.Fatal("ToGrpcError did not produce a gRPC status")
	}
	if st.Code() != codes.Internal {
		t.Errorf("gRPC code = %v, want %v", st.Code(), codes.Internal)
	}
}
