package auth

import (
	"testing"
)

func TestHashPasswordAndCheckPasswordHash(t *testing.T) {
	password := "password"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Errorf("Error in HashPassword(): %v", err)
	}

	t.Logf("hashedPassword: %v", hashedPassword)

	err = CheckPasswordHash(password, hashedPassword)
	if err != nil {
		t.Errorf("Password did not match: %v", err)
	}

	t.Logf("Success: Password matched, hashedPassword = %v, password = %v", hashedPassword, password)
}
