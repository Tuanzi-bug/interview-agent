package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "secret123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == password {
		t.Error("HashPassword returned the plain password")
	}

	if len(hash) == 0 {
		t.Error("HashPassword returned empty string")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "secret123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if !CheckPasswordHash(password, hash) {
		t.Error("CheckPasswordHash failed for correct password")
	}

	if CheckPasswordHash("wrongpassword", hash) {
		t.Error("CheckPasswordHash passed for wrong password")
	}
}
