package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestNewManager(t *testing.T) {
	t.Run("empty secret key returns error", func(t *testing.T) {
		if _, err := NewManager(""); err == nil {
			t.Fatal("expected error for empty secret key")
		}
	})

	t.Run("non-empty secret key succeeds", func(t *testing.T) {
		if _, err := NewManager("my-secret"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestGenerateAndValidateToken(t *testing.T) {
	m, err := NewManager("test-secret-key")
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	token, err := m.GenerateToken(42, "api")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	sub, err := m.ValidateToken(token)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if sub != "42" {
		t.Fatalf("expected subject '42', got %q", sub)
	}
}

func TestValidateTokenExpired(t *testing.T) {
	m, err := NewManager("test-secret-key")
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		Subject:   "1",
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secretKey)
	if err != nil {
		t.Fatalf("failed to sign expired token: %v", err)
	}

	_, err = m.ValidateToken(token)
	if !errors.Is(err, ErrTokenExpired) {
		t.Fatalf("expected ErrTokenExpired, got %v", err)
	}
}

func TestValidateTokenRejectsHS512(t *testing.T) {
	m, err := NewManager("test-secret-key")
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	claims := jwt.RegisteredClaims{Subject: "1"}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(m.secretKey)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	// ValidateToken is locked to HS256 only (alg-confusion hardening).
	if _, err := m.ValidateToken(token); err == nil {
		t.Fatal("expected error for HS512 token when locked to HS256")
	}
}

func TestValidateTokenRejectsNonHMAC(t *testing.T) {
	m, err := NewManager("test-secret-key")
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}

	claims := jwt.RegisteredClaims{Subject: "1"}
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	if _, err := m.ValidateToken(token); err == nil {
		t.Fatal("expected error for RSA-signed token")
	}
}

func TestValidateTokenRejectsGarbage(t *testing.T) {
	m, err := NewManager("test-secret-key")
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	if _, err := m.ValidateToken("not-a-jwt-token"); err == nil {
		t.Fatal("expected error for invalid token string")
	}
}
