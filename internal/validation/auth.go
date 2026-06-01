package validation

import (
	"errors"
	"strings"
)

func RegisterInput(email, password, name string) error {
	if email == "" {
		return errors.New("email is required")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	if !strings.Contains(email, "@") {
		return errors.New("email is invalid")
	}

	// TODO: validation for "name" field // ?: whats a good character limit here? like is there a max value for strings in the DB? is this something I should be worrying about?
	return nil
}
