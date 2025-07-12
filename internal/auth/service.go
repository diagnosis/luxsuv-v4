package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/diagnosis/luxsuv-v4/internal/logger"
	"github.com/diagnosis/luxsuv-v4/internal/models"
	"github.com/diagnosis/luxsuv-v4/internal/repository"
	"github.com/diagnosis/luxsuv-v4/internal/validation"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo  *repository.UserRepository
	jwtSecret string
	logger    *logger.Logger
}

func NewService(userRepo *repository.UserRepository, jwtSecret string, logger *logger.Logger) *Service {
	return &Service{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

// Register creates a new user account
func (s *Service) Register(req *models.CreateUserRequest) (*models.User, error) {
	// Normalize inputs
	username := strings.TrimSpace(req.Username)
	email := strings.TrimSpace(strings.ToLower(req.Email))
	password := req.Password
	role := strings.TrimSpace(strings.ToLower(req.Role))

	// Set default role if empty
	if role == "" {
		role = models.RoleRider
	}

	s.logger.Info(fmt.Sprintf("Registration attempt - Username: %s, Email: %s, Role: %s", username, email, role))

	// Validate input
	if err := validation.ValidateUserRegistration(username, email, password, role); err != nil {
		s.logger.Warn(fmt.Sprintf("Registration validation failed: %s", err.Error()))
		return nil, err
	}

	// Check if user already exists by email
	if existingUser, err := s.userRepo.GetUserByEmail(email); err == nil && existingUser != nil {
		s.logger.Warn(fmt.Sprintf("Registration failed: email already exists - %s", email))
		return nil, errors.New("email already exists")
	} else if err != nil && err != sql.ErrNoRows {
		s.logger.Err(fmt.Sprintf("Database error checking email: %s", err.Error()))
		return nil, errors.New("registration failed")
	}

	// Check if user already exists by username
	if existingUser, err := s.userRepo.GetUserByUsername(username); err == nil && existingUser != nil {
		s.logger.Warn(fmt.Sprintf("Registration failed: username already exists - %s", username))
		return nil, errors.New("username already exists")
	} else if err != nil && err != sql.ErrNoRows {
		s.logger.Err(fmt.Sprintf("Database error checking username: %s", err.Error()))
		return nil, errors.New("registration failed")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Err(fmt.Sprintf("Failed to hash password for user %s: %s", username, err.Error()))
		return nil, errors.New("failed to process password")
	}

	// Create user
	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
		IsAdmin:  role == models.RoleAdmin,
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		s.logger.Err(fmt.Sprintf("Failed to create user %s: %s", username, err.Error()))
		return nil, errors.New("failed to create user")
	}

	s.logger.Info(fmt.Sprintf("User registered successfully: %s (%s)", username, email))
	
	// Remove password from response
	user.Password = ""
	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *Service) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Normalize inputs
	email := strings.TrimSpace(strings.ToLower(req.Email))
	password := req.Password

	s.logger.Info(fmt.Sprintf("Login attempt for email: %s", email))

	// Validate input
	if err := validation.ValidateLoginInput(email, password); err != nil {
		s.logger.Warn(fmt.Sprintf("Login validation failed: %s", err.Error()))
		return nil, err
	}

	// Get user from database
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Warn(fmt.Sprintf("Login failed: user not found - %s", email))
			return nil, errors.New("invalid email or password")
		}
		s.logger.Err(fmt.Sprintf("Database error during login for %s: %s", email, err.Error()))
		return nil, errors.New("login failed")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.logger.Warn(fmt.Sprintf("Login failed: invalid password for %s", email))
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := s.generateJWT(user)
	if err != nil {
		s.logger.Err(fmt.Sprintf("Failed to generate JWT for %s: %s", email, err.Error()))
		return nil, errors.New("failed to generate token")
	}

	s.logger.Info(fmt.Sprintf("User logged in successfully: %s (ID: %d)", email, user.ID))

	// Remove password from response
	user.Password = ""

	return &models.LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *Service) GetUserByID(id int64) (*models.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		s.logger.Err(fmt.Sprintf("Failed to get user by ID %d: %s", id, err.Error()))
		return nil, errors.New("failed to get user")
	}

	// Remove password from response
	user.Password = ""
	return user, nil
}

// DeleteUser deletes a user (admin only)
func (s *Service) DeleteUser(userID, adminID int64) error {
	// Get admin user to verify permissions
	admin, err := s.userRepo.GetUserByID(adminID)
	if err != nil {
		return errors.New("admin user not found")
	}

	if !admin.IsAdmin {
		return errors.New("only admins can delete users")
	}

	// Prevent self-deletion
	if userID == adminID {
		return errors.New("cannot delete your own account")
	}

	// Check if target user exists
	targetUser, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return errors.New("failed to find user")
	}

	// Delete the user
	if err := s.userRepo.DeleteUser(userID); err != nil {
		s.logger.Err(fmt.Sprintf("Failed to delete user %d: %s", userID, err.Error()))
		return errors.New("failed to delete user")
	}

	s.logger.Info(fmt.Sprintf("User deleted by admin %d: %s (ID: %d)", adminID, targetUser.Email, userID))
	return nil
}

// generateJWT creates a JWT token for the user
func (s *Service) generateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// ValidateJWT validates a JWT token and returns the claims
func (s *Service) ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// GenerateResetToken generates a password reset token
func (s *Service) GenerateResetToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "password_reset",
		"exp":     time.Now().Add(1 * time.Hour).Unix(), // 1 hour expiry
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// ValidateResetToken validates a password reset token and returns user ID
func (s *Service) ValidateResetToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	// Check token type
	tokenType, ok := claims["type"]
	if !ok || tokenType != "password_reset" {
		return 0, errors.New("invalid token type")
	}

	// Get user ID
	userIDClaim, ok := claims["user_id"]
	if !ok {
		return 0, errors.New("missing user ID in token")
	}

	var userID int64
	switch v := userIDClaim.(type) {
	case float64:
		userID = int64(v)
	case int64:
		userID = v
	case int:
		userID = int64(v)
	default:
		return 0, errors.New("invalid user ID type")
	}

	return userID, nil
}