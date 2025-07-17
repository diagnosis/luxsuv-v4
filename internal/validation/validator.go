package validation

import (
	"errors"
	"regexp"
	"strings"
	"time"
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
		return nil // Will default to rider
	}
	if !models.IsValidRole(role) {
		return errors.New("invalid role; must be rider, driver, super_driver, dispatcher, or admin")
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

// ValidateBookRide validates the book ride data
func ValidateBookRide(br *models.BookRide) error {
	if br.YourName = strings.TrimSpace(br.YourName); br.YourName == "" {
		return errors.New("name is required")
	}
	if len(br.YourName) > 100 {
		return errors.New("name must be no more than 100 characters long")
	}

	if err := ValidateEmail(br.Email); err != nil {
		return err
	}

	if br.PhoneNumber = strings.TrimSpace(br.PhoneNumber); br.PhoneNumber == "" {
		return errors.New("phone number is required")
	}
	if len(br.PhoneNumber) < 7 || len(br.PhoneNumber) > 20 {
		return errors.New("phone number must be between 7 and 20 characters")
	}

	if br.RideType = strings.TrimSpace(br.RideType); br.RideType == "" {
		return errors.New("ride type is required")
	}

	if br.PickupLocation = strings.TrimSpace(br.PickupLocation); br.PickupLocation == "" {
		return errors.New("pickup location is required")
	}

	if br.DropoffLocation = strings.TrimSpace(br.DropoffLocation); br.DropoffLocation == "" {
		return errors.New("dropoff location is required")
	}

	if br.Date = strings.TrimSpace(br.Date); br.Date == "" {
		return errors.New("date is required")
	}
	// Basic date validation (assuming YYYY-MM-DD format)
	if _, err := time.Parse("2006-01-02", br.Date); err != nil {
		return errors.New("invalid date format; use YYYY-MM-DD")
	}

	if br.Time = strings.TrimSpace(br.Time); br.Time == "" {
		return errors.New("time is required")
	}
	// Basic time validation (assuming HH:MM format)
	if _, err := time.Parse("15:04", br.Time); err != nil {
		return errors.New("invalid time format; use HH:MM")
	}

	if br.NumberOfPassengers <= 0 {
		return errors.New("number of passengers must be at least 1")
	}

	if br.NumberOfLuggage < 0 {
		return errors.New("number of luggage cannot be negative")
	}

	if br.AdditionalNotes != "" && len(br.AdditionalNotes) > 500 {
		return errors.New("additional notes must be no more than 500 characters")
	}

	// Status defaults are set in DB, but validate if provided
	if br.BookStatus != "" && br.BookStatus != "Pending" {
		return errors.New("initial book status must be Pending")
	}
	if br.RideStatus != "" && br.RideStatus != "Pending" {
		return errors.New("initial ride status must be Pending")
	}

	return nil
}

// ValidateUpdateBookRide validates the update book ride data
func ValidateUpdateBookRide(updates *models.UpdateBookRideRequest) error {
	if updates.YourName != "" {
		if len(updates.YourName) > 100 {
			return errors.New("name must be no more than 100 characters long")
		}
	}

	if updates.PhoneNumber != "" {
		if len(updates.PhoneNumber) < 7 || len(updates.PhoneNumber) > 20 {
			return errors.New("phone number must be between 7 and 20 characters")
		}
	}

	if updates.Date != "" {
		if _, err := time.Parse("2006-01-02", updates.Date); err != nil {
			return errors.New("invalid date format; use YYYY-MM-DD")
		}
	}

	if updates.Time != "" {
		if _, err := time.Parse("15:04", updates.Time); err != nil {
			return errors.New("invalid time format; use HH:MM")
		}
	}

	if updates.NumberOfPassengers != nil && *updates.NumberOfPassengers <= 0 {
		return errors.New("number of passengers must be at least 1")
	}

	if updates.NumberOfLuggage != nil && *updates.NumberOfLuggage < 0 {
		return errors.New("number of luggage cannot be negative")
	}

	if updates.AdditionalNotes != "" && len(updates.AdditionalNotes) > 500 {
		return errors.New("additional notes must be no more than 500 characters")
	}

	return nil
}

// ValidateBookingDateTime validates that booking is at least 24 hours in the future
func ValidateBookingDateTime(dateStr, timeStr string) error {
	if dateStr == "" || timeStr == "" {
		return errors.New("date and time are required")
	}

	// Parse the booking date and time
	bookingDateTime, err := time.Parse("2006-01-02 15:04", dateStr+" "+timeStr)
	if err != nil {
		return errors.New("invalid date or time format")
	}

	// Check if booking is at least 24 hours in the future
	now := time.Now()
	minBookingTime := now.Add(24 * time.Hour)

	if bookingDateTime.Before(minBookingTime) {
		return errors.New("booking must be at least 24 hours in advance")
	}

	return nil
}