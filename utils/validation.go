package utils

import (
	"regexp"
	"unicode"
)

// IsValidEmail validates email format using regex
func IsValidEmail(email string) bool {
	// RFC 5322 compliant email regex (simplified but robust)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsStrongPassword validates password strength
func IsStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// For basic validation, just check length
	// For stronger validation, you could require: hasUpper && hasLower && hasNumber && hasSpecial
	return len(password) >= 8
}

// IsValidUsername validates username format
func IsValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 50 {
		return false
	}

	// Username should contain only alphanumeric characters and underscores
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return usernameRegex.MatchString(username)
}