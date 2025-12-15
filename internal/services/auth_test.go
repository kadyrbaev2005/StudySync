package services

import (
	"testing"
	"time"
)

// Проверяем, что токен генерируется и корректно парсится
func TestGenerateAndParseJWT_Success(t *testing.T) {
	token, err := GenerateJWT(42, RoleAdmin)
	if err != nil {
		t.Fatalf("GenerateJWT returned error: %v", err)
	}
	if token == "" {
		t.Fatalf("expected non-empty token")
	}

	claims, err := ParseJWT(token)
	if err != nil {
		t.Fatalf("ParseJWT returned error: %v", err)
	}

	if claims.UserID != 42 {
		t.Errorf("expected UserID 42, got %d", claims.UserID)
	}
	if claims.Role != RoleAdmin {
		t.Errorf("expected Role %q, got %q", RoleAdmin, claims.Role)
	}

	if claims.ExpiresAt == nil {
		t.Fatalf("expected ExpiresAt to be set")
	}
	if time.Until(claims.ExpiresAt.Time) <= 0 {
		t.Fatalf("expected token to be valid in the future, got expired")
	}
}

// Невалидный токен должен давать ошибку
func TestParseJWT_InvalidToken(t *testing.T) {
	_, err := ParseJWT("not.a.valid.token")
	if err == nil {
		t.Fatalf("expected error for invalid token, got nil")
	}
}

// Проверяем hash + сравнение пароля
func TestHashPasswordAndCheckPasswordHash(t *testing.T) {
	pw := "super-secret"
	hash := HashPassword(pw)

	if hash == pw {
		t.Fatalf("hash must not be equal to original password")
	}

	if ok := CheckPasswordHash(pw, hash); !ok {
		t.Fatalf("expected CheckPasswordHash to return true for correct password")
	}

	if ok := CheckPasswordHash("wrong-password", hash); ok {
		t.Fatalf("expected CheckPasswordHash to return false for wrong password")
	}
}
