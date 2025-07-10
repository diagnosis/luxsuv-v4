package validation

import (
	"errors"
	"regexp"
	"strings"
	"unicode"

	"github.com/diagnosis/luxsuv-v4/internal/models"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// ValidateUserRegistration validates user registration data
func ValidateUserRegistration(username, email, password, role string) error {
	if err := ValidateUsername(username); err != nil {
		return err
	}
	if err := ValidateEmail(email); err != nil {
		return err
	}
	if err := ValidatePassword(password); err != nil {
		return err
	}
	if err := ValidateRole(role); err != nil {
		return err
	}
	return nil
}

// ValidateUsername validates username
func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return errors.New("username is required")
	}
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	if len(username) > 50 {
		return errors.New("username must be no more than 50 characters long")
	}
	
	// Check for valid characters (alphanumeric and underscore only)
	for _, char := range username {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' {
			return errors.New("username can only contain letters, numbers, and underscores")
		}
	}
	return nil
}

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return errors.New("email is required")
	}
	if len(email) > 254 {
		return errors.New("email address is too long")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password is required")
	}
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if len(password) > 128 {
		return errors.New("password must be no more than 128 characters long")
	}

	var hasLetter, hasDigit bool
	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
	}

	if !hasLetter {
		return errors.New("password must contain at least one letter")
	}
	if !hasDigit {
		return errors.New("password must contain at least one number")
	}

	return nil
}

// ValidateRole validates user role
func ValidateRole(role string) error {
	role = strings.TrimSpace(strings.ToLower(role))
	if role == "" {
		return nil // Empty role will default to rider
	}
	if !models.IsValidRole(role) {
		return errors.New("role must be one of: rider, driver, admin")
	}
	return nil
}

// ValidateLoginInput validates login credentials
func ValidateLoginInput(email, password string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email is required")
	}
	if password == "" {
		return errors.New("password is required")
	}
	return nil
}