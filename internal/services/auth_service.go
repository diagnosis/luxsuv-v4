package services

import (
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
	// Input validation
	if user.Username == "" {
		s.log.Warn("Registration failed: username is required")
		return errors.New("username is required")
	}
	if user.Email == "" {
		s.log.Warn("Registration failed: email is required")
		return errors.New("email is required")
	}
	if user.Password == "" {
		s.log.Warn("Registration failed: password is required")
		return errors.New("password is required")
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(user.Email) {
		s.log.Warn("Registration failed: invalid email format for " + user.Email)
		return errors.New("invalid email format")
	}

	// Validate password strength
	if len(user.Password) < 8 {
		s.log.Warn("Registration failed: password too short for user " + user.Username)
		return errors.New("password must be at least 8 characters long")
	}
	if !regexp.MustCompile(`[a-zA-Z]`).MatchString(user.Password) || !regexp.MustCompile(`[0-9]`).MatchString(user.Password) {
		s.log.Warn("Registration failed: password must contain letters and numbers for user " + user.Username)
		return errors.New("password must contain letters and numbers")
	}

	// Validate role
	validRoles := map[string]bool{"admin": true, "driver": true, "rider": true}
	user.Role = strings.ToLower(user.Role)
	if user.Role == "" {
		user.Role = "rider"
	} else if !validRoles[user.Role] {
		s.log.Warn("Registration failed: invalid role " + user.Role + " for user " + user.Username)
		return errors.New("invalid role; must be admin, driver, or rider")
	}

	// Set super_admin for admin role
	user.SuperAdmin = user.Role == "admin"

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Err("Registration failed: error hashing password for user " + user.Username + ": " + err.Error())
		return err
	}
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()

	if err := s.repo.CreateUser(user); err != nil {
		s.log.Err("Registration failed: database error for user " + user.Username + ": " + err.Error())
		return err
	}
	s.log.Info("User registered successfully: " + user.Username)
	return nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	if email == "" || password == "" {
		s.log.Warn("Login failed: email and password are required")
		return "", errors.New("email and password are required")
	}

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		s.log.Warn("Login failed: user not found for email " + email)
		return "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.log.Warn("Login failed: invalid password for email " + email)
		return "", errors.New("invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":          user.ID,
		"role":        user.Role,
		"super_admin": user.SuperAdmin,
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		s.log.Err("Login failed: error generating JWT for email " + email + ": " + err.Error())
		return "", err
	}

	s.log.Info("User logged in successfully: " + email)
	return tokenString, nil
}

func (s *AuthService) GetUserByID(id int64) (*data.User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		s.log.Err("Failed to get user by ID " + fmt.Sprintf("%d", id) + ": " + err.Error())
		return nil, err
	}
	s.log.Info("Retrieved user by ID: " + fmt.Sprintf("%d", id))
	return user, nil
}
