package utils

import "testing"

func TestValidateRegister_Valid(t *testing.T) {
	errs := ValidateRegister("test@example.com", "username", "password123", "password123")
	if errs.HasErrors() {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidateRegister_EmptyEmail(t *testing.T) {
	errs := ValidateRegister("", "username", "password123", "password123")
	if _, ok := errs["email"]; !ok {
		t.Error("expected email error")
	}
}

func TestValidateRegister_InvalidEmail(t *testing.T) {
	errs := ValidateRegister("notanemail", "username", "password123", "password123")
	if _, ok := errs["email"]; !ok {
		t.Error("expected email error for invalid email")
	}
}

func TestValidateRegister_ShortUsername(t *testing.T) {
	errs := ValidateRegister("test@test.com", "ab", "password123", "password123")
	if _, ok := errs["username"]; !ok {
		t.Error("expected username error for short name")
	}
}

func TestValidateRegister_ShortPassword(t *testing.T) {
	errs := ValidateRegister("test@test.com", "username", "123", "123")
	if _, ok := errs["password"]; !ok {
		t.Error("expected password error for short password")
	}
}

func TestValidateRegister_PasswordMismatch(t *testing.T) {
	errs := ValidateRegister("test@test.com", "username", "password123", "different")
	if _, ok := errs["confirm_password"]; !ok {
		t.Error("expected confirm_password error")
	}
}

func TestValidateLogin_Empty(t *testing.T) {
	errs := ValidateLogin("", "")
	if len(errs) != 2 {
		t.Errorf("expected 2 errors, got %d", len(errs))
	}
}

func TestValidatePost_Valid(t *testing.T) {
	errs := ValidatePost("Title", "Content", []int{1})
	if errs.HasErrors() {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidatePost_NoCategories(t *testing.T) {
	errs := ValidatePost("Title", "Content", []int{})
	if _, ok := errs["categories"]; !ok {
		t.Error("expected categories error")
	}
}

func TestValidateComment_Empty(t *testing.T) {
	errs := ValidateComment("")
	if _, ok := errs["content"]; !ok {
		t.Error("expected content error for empty comment")
	}
}

func TestValidateComment_Whitespace(t *testing.T) {
	errs := ValidateComment("   ")
	if _, ok := errs["content"]; !ok {
		t.Error("expected content error for whitespace-only comment")
	}
}
