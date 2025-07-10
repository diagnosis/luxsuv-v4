package services

import (
	"database/sql"
	"errors"
	"fmt"
	"luxsuv-v4/models"
	"luxsuv-v4/utils"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db        *sql.DB
	jwtSecret string
}

func NewAuthService(db *sql.DB, jwtSecret string) *AuthService {
	return &AuthService{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

// ValidateEmail validates email format using regex
func (s *AuthService) ValidateEmail(email string) error {
	// RFC 5322 compliant email regex (simplified but robust)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

// ValidatePassword validates password strength
func (s *AuthService) ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}

// ValidateUsername validates username requirements
func (s *AuthService) ValidateUsername(username string) error {
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	if len(username) > 50 {
		return errors.New("username must be less than 50 characters")
	}
	return nil
}

// RegisterUser creates a new user account
func (s *AuthService) RegisterUser(req models.RegisterRequest) (*models.User, error) {
	// Validate input
	if err := s.ValidateEmail(req.Email); err != nil {
		return nil, err
	}
	
	if err := s.ValidatePassword(req.Password); err != nil {
		return nil, err
	}
	
	if err := s.ValidateUsername(req.Username); err != nil {
		return nil, err
	}

	// Set default role if not provided
	if req.Role == "" {
		req.Role = "rider"
	}

	// Validate role
	validRoles := []string{"rider", "driver", "admin"}
	roleValid := false
	for _, role := range validRoles {
		if req.Role == role {
			roleValid = true
			break
		}
	}
	if !roleValid {
		return nil, errors.New("invalid role")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Insert user into database
	query := `
		INSERT INTO users (username, email, password, role, super_admin, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, username, email, role, super_admin, created_at
	`
	
	user := &models.User{}
	err = s.db.QueryRow(
		query,
		req.Username,
		req.Email,
		string(hashedPassword),
		req.Role,
		false, // super_admin defaults to false
		time.Now(),
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.SuperAdmin,
		&user.CreatedAt,
	)

	if err != nil {
		// Check for PostgreSQL unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				if strings.Contains(pqErr.Detail, "email") {
					return nil, errors.New("email already exists")
				}
				if strings.Contains(pqErr.Detail, "username") {
					return nil, errors.New("username already exists")
				}
			}
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// LoginUser authenticates a user and returns a JWT token
func (s *AuthService) LoginUser(req models.LoginRequest) (string, *models.User, error) {
	// Validate email format
	if err := s.ValidateEmail(req.Email); err != nil {
		return "", nil, err
	}

	// Get user from database
	user := &models.User{}
	query := `
		SELECT id, username, email, password, role, super_admin, created_at
		FROM users
		WHERE email = $1
	`
	
	err := s.db.QueryRow(query, req.Email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.SuperAdmin,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil, errors.New("invalid credentials")
		}
		return "", nil, fmt.Errorf("database error: %w", err)
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.GenerateJWT(user)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Clear password from user object before returning
	user.Password = ""

	return token, user, nil
}

// GenerateJWT creates a JWT token for the user
func (s *AuthService) GenerateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":     user.ID,
		"username":    user.Username,
		"email":       user.Email,
		"role":        user.Role,
		"super_admin": user.SuperAdmin,
		"exp":         time.Now().Add(time.Hour * 24).Unix(), // 24 hours
		"iat":         time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// ValidateJWT validates and parses a JWT token
func (s *AuthService) ValidateJWT(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, errors.New("invalid token")
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(userID int) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, username, email, role, super_admin, created_at
		FROM users
		WHERE id = $1
	`
	
	err := s.db.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.SuperAdmin,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return user, nil
}