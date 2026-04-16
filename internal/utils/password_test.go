package utils

import "testing"

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("password123")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if hash == "" {
		t.Fatal("hash should not be empty")
	}
	if hash == "password123" {
		t.Fatal("hash should not equal plain text")
	}
}

func TestCheckPassword(t *testing.T) {
	hash, _ := HashPassword("mypassword")

	if !CheckPassword(hash, "mypassword") {
		t.Error("CheckPassword should return true for correct password")
	}
	if CheckPassword(hash, "wrongpassword") {
		t.Error("CheckPassword should return false for wrong password")
	}
}

func TestCheckPasswordEmpty(t *testing.T) {
	hash, _ := HashPassword("test")
	if CheckPassword(hash, "") {
		t.Error("CheckPassword should return false for empty password")
	}
}
