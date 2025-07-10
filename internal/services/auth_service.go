package services

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/diagnosis/luxsuv-v4/internal/data"
	"github.com/diagnosis/luxsuv-v4/internal/logger"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      *data.Repository
	jwtSecret string
	log       *logger.Logger
}

func NewAuthService(repo *data.Repository, jwtSecret string, log *logger.Logger) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret, log: log}
}

func (s *AuthService) Register(user *data.User) error {
	s.log.Info("Starting registration process for user: " + user.Username)
	
	// Input validation with detailed error messages
	if err := s.validateUserInput(user); err != nil {
		s.log.Warn("User input validation failed: " + err.Error())
		return err
	}

	// Normalize and validate email
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if err := s.validateEmail(user.Email); err != nil {
		s.log.Warn("Email validation failed for: " + user.Email + " - " + err.Error())
		return err
	}

	// Validate password strength
	if err := s.validatePassword(user.Password); err != nil {
		s.log.Warn("Password validation failed: " + err.Error())
		return err
	}

	// Validate and set role
	if err := s.validateAndSetRole(user); err != nil {
		s.log.Warn("Role validation failed: " + err.Error())
		return err
	}

	s.log.Info("All validations passed, proceeding to hash password")

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Err("Registration failed: error hashing password for user " + user.Username + ": " + err.Error())
		return errors.New("failed to process password")
	}
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()

	// Create user in database
	if err := s.repo.CreateUser(user); err != nil {
		s.log.Err("Registration failed: database error for user " + user.Username + ": " + err.Error())
		return err // Let the handler deal with specific database errors
	}

	s.log.Info("User registered successfully: " + user.Username + " (" + user.Email + ")")
	return nil
}

func (s *AuthService) validateUserInput(user *data.User) error {
	if user == nil {
		return errors.New("user data is required")
	}

	user.Username = strings.TrimSpace(user.Username)
	if user.Username == "" {
		return errors.New("username is required")
	}

	if len(user.Username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}

	if len(user.Username) > 50 {
		return errors.New("username must be no more than 50 characters long")
	}

	if user.Email == "" {
		return errors.New("email is required")
	}

	if user.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

func (s *AuthService) validateEmail(email string) error {
	s.log.Info("Validating email: " + email)
	
	if email == "" {
		return errors.New("email is required")
	}
	
	// More comprehensive email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	
	if !emailRegex.MatchString(email) {
		s.log.Warn("Email regex failed for: " + email)
		return errors.New("invalid email format")
	}

	if len(email) > 254 {
		return errors.New("email address is too long")
	}

	s.log.Info("Email validation passed for: " + email)
	return nil
}

func (s *AuthService) validatePassword(password string) error {
	s.log.Info("Validating password strength")
	
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if len(password) > 128 {
		return errors.New("password must be no more than 128 characters long")
	}

	// Check for at least one letter and one number
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasLetter || !hasNumber {
		return errors.New("password must contain both letters and numbers")
	}

	s.log.Info("Password validation passed")
	return nil
}

func (s *AuthService) validateAndSetRole(user *data.User) error {
	s.log.Info("Validating role: " + user.Role)
	
	validRoles := map[string]bool{"admin": true, "driver": true, "rider": true}
	
	user.Role = strings.TrimSpace(strings.ToLower(user.Role))
	if user.Role == "" {
		user.Role = "rider" // Default role
		s.log.Info("Setting default role: rider")
	} else {
		if !validRoles[user.Role] {
			s.log.Warn("Invalid role provided: " + user.Role)
			return errors.New("invalid role; must be admin, driver, or rider")
		}
	}

	user.SuperAdmin = user.Role == "admin"
	s.log.Info("Role validation passed. Role: " + user.Role + ", SuperAdmin: " + fmt.Sprintf("%t", user.SuperAdmin))

	return nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	// Input validation
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" {
		s.log.Warn("Login failed: email and password are required")
		return "", errors.New("email and password are required")
	}

	// Get user from database
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Warn("Login failed: user not found for email " + email)
			return "", errors.New("invalid email or password")
		}
		s.log.Err("Login failed: database error for email " + email + ": " + err.Error())
		return "", errors.New("login failed")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.log.Warn("Login failed: invalid password for email " + email)
		return "", errors.New("invalid email or password")
	}

	// Generate JWT token with all required claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":          user.ID,
		"username":    user.Username,
		"email":       user.Email,
		"role":        user.Role,
		"super_admin": user.SuperAdmin,
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
		"iat":         time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		s.log.Err("Login failed: error generating JWT for email " + email + ": " + err.Error())
		return "", errors.New("failed to generate token")
	}

	s.log.Info("User logged in successfully: " + email + " (ID: " + fmt.Sprintf("%d", user.ID) + ")")
	return tokenString, nil
}

func (s *AuthService) GetUserByID(id int64) (*data.User, error) {
	if id <= 0 {
		s.log.Warn(fmt.Sprintf("Invalid user ID: %d", id))
		return nil, errors.New("invalid user ID")
	}

	user, err := s.repo.GetUserByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Warn(fmt.Sprintf("User not found for ID: %d", id))
			return nil, sql.ErrNoRows
		}
		s.log.Err("Failed to get user by ID " + fmt.Sprintf("%d", id) + ": " + err.Error())
		return nil, err
	}

	s.log.Info("Retrieved user by ID: " + fmt.Sprintf("%d", id) + " (" + user.Email + ")")
	return user, nil
}