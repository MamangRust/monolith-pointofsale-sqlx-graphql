package hash

import (
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	h := NewHashingPassword()

	hashed, err := h.HashPassword("secret123")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if hashed == "secret123" {
		t.Fatal("hash must not equal plaintext password")
	}

	if !strings.HasPrefix(hashed, "$2") {
		t.Fatalf("expected bcrypt hash prefix '$2', got %q", hashed)
	}
}

func TestHashPasswordProducesUniqueSalt(t *testing.T) {
	h := NewHashingPassword()

	hash1, err := h.HashPassword("same-password")
	if err != nil {
		t.Fatalf("failed to hash: %v", err)
	}
	hash2, err := h.HashPassword("same-password")
	if err != nil {
		t.Fatalf("failed to hash: %v", err)
	}

	if hash1 == hash2 {
		t.Fatal("expected different hashes for same password (random salt)")
	}
}

func TestComparePassword(t *testing.T) {
	h := NewHashingPassword()

	hashed, err := h.HashPassword("secret123")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	t.Run("correct password matches", func(t *testing.T) {
		if err := h.ComparePassword(hashed, "secret123"); err != nil {
			t.Fatalf("expected match, got %v", err)
		}
	})

	t.Run("wrong password rejected", func(t *testing.T) {
		if err := h.ComparePassword(hashed, "wrong-password"); err != bcrypt.ErrMismatchedHashAndPassword {
			t.Fatalf("expected bcrypt mismatch error, got %v", err)
		}
	})
}
