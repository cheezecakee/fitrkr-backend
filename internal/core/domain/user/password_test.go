package user_test

import (
	"testing"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

func TestNewPassword_ValidCases(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "valid password - meet all requirements",
			password: "SecurePass123!",
		},
		{
			name:     "valid password - longer password",
			password: "SecurePass123!AbCdEfGhIjKlMnOpQrStUvWxYz0123456789!@#$%^&*()",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := user.NewPassword(tt.password)
			if err != nil {
				t.Errorf("NewPassword() unexpected error = %v", err)
			}
		})
	}
}

func TestNewPassword_InvalidCases(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "invalid password - too short",
			password: "Short1!",
			wantErr:  true,
		},
		{
			name:     "invalid password - too long",
			password: "SecurePass123!AbCdEfGhIjKlMnOpQrStUvWxYz0123456789!@#$%^&*()XYSCS",
			wantErr:  true,
		},
		{
			name:     "invalid password - no uppercase",
			password: "lowercase123!",
			wantErr:  true,
		},
		{
			name:     "invalid password - no lowercase",
			password: "UPPERCASE123!",
			wantErr:  true,
		},
		{
			name:     "invalid password - no digit",
			password: "NoDigitsHere!",
			wantErr:  true,
		},
		{
			name:     "invalid password - too special character",
			password: "NoSpecial123",
			wantErr:  true,
		},
		{
			name:     "invalid password - empty",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := user.NewPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPassword_HashingAndComparison(t *testing.T) {
	password := "SecurePass123!"

	hashedPass, err := user.NewPassword(password)
	if err != nil {
		t.Fatalf("NewPassword() unexpected error = %v", err)
	}

	// Test 1: Password should be hashed (not plaintext)
	if string(hashedPass) == password {
		t.Errorf("Password was not hashed")
	}

	// Test 2: Correct password should match
	if !hashedPass.Verify(password) {
		t.Error("Verify() failed for correct password")
	}

	// Test 3: Incorrect password should not match
	if hashedPass.Verify("WrongPassword123!") {
		t.Error("Verify() succeeded for incorrect password")
	}

	// Test 4: Similar but different password should not match
	if hashedPass.Verify("SecurePass123") {
		t.Error("Verify() succeeded for similar but incorrect password")
	}
}

func TestPassword_UniqueHashes(t *testing.T) {
	password := "SecurePass123!"

	// Hash the same password twice
	hash1, err := user.NewPassword(password)
	if err != nil {
		t.Fatalf("NewPassword() error = %v", err)
	}

	hash2, err := user.NewPassword(password)
	if err != nil {
		t.Fatalf("NewPassword() error = %v", err)
	}

	// Hashes should be different (bcrypt uses salt)
	if string(hash1) == string(hash2) {
		t.Error("Same password produced identical hashes (salt not working)")
	}

	// Both passwords should still match the original password
	if !hash1.Verify(password) {
		t.Error("hash1 doesn't match original password")
	}
	if !hash2.Verify(password) {
		t.Error("hash2 doesn't match original password")
	}
}
