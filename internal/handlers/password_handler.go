package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/diagnosis/luxsuv-v4/internal/auth"
	"github.com/diagnosis/luxsuv-v4/internal/email"
	"github.com/diagnosis/luxsuv-v4/internal/logger"
	"github.com/diagnosis/luxsuv-v4/internal/repository"
	"github.com/diagnosis/luxsuv-v4/internal/validation"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHandler struct {
	authService  *auth.Service
	userRepo     *repository.UserRepository
	emailService *email.Service
	logger       *logger.Logger
}

func NewPasswordHandler(authService *auth.Service, userRepo *repository.UserRepository, emailService *email.Service, logger *logger.Logger) *PasswordHandler {
	return &PasswordHandler{
		authService:  authService,
		userRepo:     userRepo,
		emailService: emailService,
		logger:       logger,
	}
}

// ChangePassword handles password change for authenticated users
func (h *PasswordHandler) ChangePassword(c echo.Context) error {
	// Get user ID from context
	userIDClaim := c.Get("user_id")
	if userIDClaim == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid token",
		})
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
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid token claims",
		})
	}

	var req struct {
		CurrentPassword string `json:"current_password" validate:"required"`
		NewPassword     string `json:"new_password" validate:"required,min=8"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	// Validate new password
	if err := validation.ValidatePassword(req.NewPassword); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Get user from database
	user, err := h.userRepo.GetUserByID(userID)
	if err != nil {
		h.logger.Err(fmt.Sprintf("Failed to get user for password change: %s", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to process request",
		})
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		h.logger.Warn(fmt.Sprintf("Password change failed: invalid current password for user %d", userID))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "current password is incorrect",
		})
	}

	// Check if new password is different from current
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.NewPassword)); err == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "new password must be different from current password",
		})
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Err(fmt.Sprintf("Failed to hash new password for user %d: %s", userID, err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to process new password",
		})
	}

	// Update password in database
	if err := h.userRepo.UpdatePassword(userID, string(hashedPassword)); err != nil {
		h.logger.Err(fmt.Sprintf("Failed to update password for user %d: %s", userID, err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to update password",
		})
	}

	h.logger.Info(fmt.Sprintf("Password changed successfully for user %d", userID))
	return c.JSON(http.StatusOK, map[string]string{
		"message": "password changed successfully",
	})
}

// ResetPasswordRequest handles password reset request (generates reset token)
func (h *PasswordHandler) ResetPasswordRequest(c echo.Context) error {
	var req struct {
		Email string `json:"email" validate:"required,email"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	if err := validation.ValidateEmail(email); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Check if user exists
	user, err := h.userRepo.GetUserByEmail(email)
	if err != nil {
		// Don't reveal if email exists or not for security
		h.logger.Warn(fmt.Sprintf("Password reset requested for non-existent email: %s", email))
		return c.JSON(http.StatusOK, map[string]string{
			"message": "if the email exists, a password reset link has been sent",
		})
	}

	// Generate reset token (in a real app, you'd send this via email)
	resetToken, err := h.authService.GenerateResetToken(user.ID)
	if err != nil {
		h.logger.Err(fmt.Sprintf("Failed to generate reset token for user %d: %s", user.ID, err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to process reset request",
		})
	}

	// Store reset token in database
	if err := h.userRepo.StoreResetToken(user.ID, resetToken); err != nil {
		h.logger.Err(fmt.Sprintf("Failed to store reset token for user %d (%s): %s", user.ID, email, err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to process reset request",
		})
	}

	h.logger.Info(fmt.Sprintf("Password reset token generated successfully for user %s (ID: %d)", email, user.ID))

	// Send email if email service is configured
	if h.emailService != nil {
		h.logger.Info(fmt.Sprintf("Attempting to send password reset email to %s", email))
		if err := h.emailService.SendPasswordResetEmail(email, resetToken); err != nil {
			h.logger.Err(fmt.Sprintf("Failed to send password reset email to %s: %s", email, err.Error()))
			// Don't fail the request if email fails, but log it
			h.logger.Warn("Email service failed, falling back to token response")
			return c.JSON(http.StatusOK, map[string]interface{}{
				"message":     "password reset token generated (email service failed)",
				"reset_token": resetToken,
			})
		}

		h.logger.Info(fmt.Sprintf("Password reset email sent successfully to %s", email))
		return c.JSON(http.StatusOK, map[string]string{
			"message": "if the email exists, a password reset link has been sent",
		})
	} else {
		// In development mode without email service, return the token
		h.logger.Warn("Email service not configured, returning reset token in response")
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":     "password reset token generated (email service disabled)",
			"reset_token": resetToken,
		})
	}
}

// ResetPassword handles password reset with token
func (h *PasswordHandler) ResetPassword(c echo.Context) error {
	var req struct {
		ResetToken  string `json:"reset_token" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=8"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	// Validate new password
	if err := validation.ValidatePassword(req.NewPassword); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Validate reset token and get user ID
	userID, err := h.authService.ValidateResetToken(req.ResetToken)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("Invalid reset token used: %s", err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid or expired reset token",
		})
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Err(fmt.Sprintf("Failed to hash password for reset: %s", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to process new password",
		})
	}

	// Update password
	if err := h.userRepo.UpdatePassword(userID, string(hashedPassword)); err != nil {
		h.logger.Err(fmt.Sprintf("Failed to update password during reset: %s", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to reset password",
		})
	}

	// Invalidate reset token
	if err := h.userRepo.InvalidateResetToken(userID); err != nil {
		h.logger.Warn(fmt.Sprintf("Failed to invalidate reset token: %s", err.Error()))
	}

	h.logger.Info(fmt.Sprintf("Password reset successfully for user %d", userID))
	return c.JSON(http.StatusOK, map[string]string{
		"message": "password reset successfully",
	})
}
