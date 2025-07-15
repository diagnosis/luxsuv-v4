package handlers

import (
	"fmt"
	"github.com/diagnosis/luxsuv-v4/internal/logger"
	"github.com/diagnosis/luxsuv-v4/internal/models"
	"github.com/diagnosis/luxsuv-v4/internal/repository"
	"github.com/diagnosis/luxsuv-v4/internal/validation"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
	"strconv"
)

type BookRideHandler struct {
	repo   repository.BookRideRepository
	logger *logger.Logger
}

func NewBookRideHandler(repo repository.BookRideRepository, logger *logger.Logger) *BookRideHandler {
	return &BookRideHandler{
		repo:   repo,
		logger: logger,
	}
}

func (h *BookRideHandler) Create(c echo.Context) error {
	br := &models.BookRide{}
	if err := c.Bind(br); err != nil {
		h.logger.Warn(fmt.Sprintf("Invalid request body: %s", err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	h.logger.Info(fmt.Sprintf("Received booking data: %+v", br))

	//validate input
	if err := validation.ValidateBookRide(br); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Get user ID from context if user is authenticated
	userIDClaim := c.Get("user_id")
	h.logger.Info(fmt.Sprintf("User ID claim from context: %v (type: %T)", userIDClaim, userIDClaim))

	if userIDClaim != nil {
		var userID int64
		switch v := userIDClaim.(type) {
		case float64:
			userID = int64(v)
		case int64:
			userID = v
		case int:
			userID = int64(v)
		case int32:
			userID = int64(v)
		default:
			h.logger.Warn(fmt.Sprintf("Unexpected user_id type: %T, value: %v", userIDClaim, userIDClaim))
			userID = 0
		}

		if userID > 0 {
			br.UserID = &userID
			h.logger.Info(fmt.Sprintf("✅ User ID successfully set: %d", userID))
		} else {
			h.logger.Warn("❌ Failed to extract valid user ID from claims")
		}
	} else {
		h.logger.Info("No user_id in context - guest booking")
	}

	h.logger.Info(fmt.Sprintf("Final booking before DB save - UserID: %v, Name: %s, Email: %s",
		func() interface{} {
			if br.UserID != nil {
				return *br.UserID
			}
			return "nil"
		}(), br.YourName, br.Email))

	br.BookStatus = "Pending"
	br.RideStatus = "Pending"

	if err := h.repo.Create(c.Request().Context(), br); err != nil {
		h.logger.Err(fmt.Sprintf("Error creating book ride: %s", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error creating book ride"})
	}

	h.logger.Info(fmt.Sprintf("Booking created successfully: ID %d", br.ID))
	return c.JSON(http.StatusCreated, br)
}

func (h *BookRideHandler) GetByEmail(c echo.Context) error {
	encodedEmail := c.Param("email")
	email, err := url.PathUnescape(encodedEmail)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid email format"})
	}
	if email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "email is required!"})
	}
	if err := validation.ValidateEmail(email); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	bookings, err := h.repo.GetByEmail(c.Request().Context(), email)
	if err != nil {
		h.logger.Err(fmt.Sprintf("Failed to get bookings by email %s, %s", email, err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error getting bookings by email"})
	}
	if len(bookings) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "booking not found"})
	}
	h.logger.Info(fmt.Sprintf("Retrieved %d bookings for email %s", len(bookings), email))
	return c.JSON(http.StatusOK, bookings)
}

// accept
func (h *BookRideHandler) Accept(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid booking ID"})
	}

	driverID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "driver not authorized"})
	}

	//role check
	role, ok := c.Get("role").(string)
	if !ok || role != models.RoleDriver {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "access denied: driver role required"})
	}

	if err := h.repo.Accept(c.Request().Context(), id, driverID); err != nil {
		h.logger.Err(fmt.Sprintf("Failed to accept book ride: %s", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error accepting book ride"})
	}
	h.logger.Info(fmt.Sprintf("Booking accepted successfully: ID %d", id))
	return c.JSON(http.StatusOK, map[string]string{"message": "booking accepted successfully"})
}

func (h *BookRideHandler) GetByUserID(c echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not authorized"})
	}

	bookings, err := h.repo.GetByUserID(c.Request().Context(), userID)
	if err != nil {
		h.logger.Err(fmt.Sprintf("Failed to get bookings by user id %d: %s", userID, err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error getting bookings by user id"})
	}
	h.logger.Info(fmt.Sprintf("Retrieved %d bookings for user id %d", len(bookings), userID))
	return c.JSON(http.StatusOK, bookings)

}

// Update handles updating a booking (authenticated users or via secure token)
func (h *BookRideHandler) Update(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid booking ID"})
	}

	var updates models.UpdateBookRideRequest
	if err := c.Bind(&updates); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	// Validate the updates
	if err := validation.ValidateUpdateBookRide(&updates); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Check if user is authenticated or using secure token
	userID := c.Get("user_id")
	token := c.QueryParam("token")

	var booking *models.BookRide

	if userID != nil {
		// Authenticated user - verify they own the booking
		var uid int64
		switch v := userID.(type) {
		case float64:
			uid = int64(v)
		case int64:
			uid = v
		case int:
			uid = int64(v)
		default:
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user authentication"})
		}

		booking, err = h.repo.GetByID(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "booking not found"})
		}

		if booking.UserID == nil || *booking.UserID != uid {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
		}
	} else if token != "" {
		// Guest user with secure token
		bookingID, email, err := h.authService.ValidateBookingUpdateToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
		}

		if bookingID != id {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "token not valid for this booking"})
		}

		booking, err = h.repo.GetByIDAndEmail(c.Request().Context(), id, email)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "booking not found"})
		}
	} else {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "authentication required"})
	}

	// Check if booking can be updated (not cancelled or completed)
	if booking.BookStatus == models.BookStatusCancelled || booking.BookStatus == models.BookStatusCompleted {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot update cancelled or completed booking"})
	}

	// If date/time is being updated, validate 24-hour rule
	dateToCheck := booking.Date
	timeToCheck := booking.Time

	if updates.Date != "" {
		dateToCheck = updates.Date
	}
	if updates.Time != "" {
		timeToCheck = updates.Time
	}

	if err := validation.ValidateBookingDateTime(dateToCheck, timeToCheck); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Perform the update
	if err := h.repo.Update(c.Request().Context(), id, &updates); err != nil {
		h.logger.Err(fmt.Sprintf("Failed to update booking %d: %s", id, err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update booking"})
	}

	// Get updated booking
	updatedBooking, err := h.repo.GetByID(c.Request().Context(), id)
	if err != nil {
		h.logger.Err(fmt.Sprintf("Failed to get updated booking %d: %s", id, err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "booking updated but failed to retrieve updated data"})
	}

	h.logger.Info(fmt.Sprintf("Booking updated successfully: ID %d", id))
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "booking updated successfully",
		"booking": updatedBooking,
	})
}

// Cancel handles cancelling a booking (authenticated users or via secure token)
func (h *BookRideHandler) Cancel(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid booking ID"})
	}

	var req struct {
		Reason string `json:"reason,omitempty"`
	}
	if err := c.Bind(&req); err != nil {
		// Reason is optional, so binding errors are not critical
		req.Reason = "Cancelled by user"
	}

	if req.Reason == "" {
		req.Reason = "Cancelled by user"
	}

	// Check if user is authenticated or using secure token
	userID := c.Get("user_id")
	token := c.QueryParam("token")

	var booking *models.BookRide

	if userID != nil {
		// Authenticated user - verify they own the booking
		var uid int64
		switch v := userID.(type) {
		case float64:
			uid = int64(v)
		case int64:
			uid = v
		case int:
			uid = int64(v)
		default:
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user authentication"})
		}

		booking, err = h.repo.GetByID(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "booking not found"})
		}

		if booking.UserID == nil || *booking.UserID != uid {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
		}
	} else if token != "" {
		// Guest user with secure token
		bookingID, email, err := h.authService.ValidateBookingUpdateToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
		}

		if bookingID != id {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "token not valid for this booking"})
		}

		booking, err = h.repo.GetByIDAndEmail(c.Request().Context(), id, email)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "booking not found"})
		}
	} else {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "authentication required"})
	}

	// Check if booking can be cancelled
	if booking.BookStatus == models.BookStatusCancelled {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "booking is already cancelled"})
	}

	if booking.BookStatus == models.BookStatusCompleted {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot cancel completed booking"})
	}

	// Validate 24-hour cancellation rule
	if err := validation.ValidateBookingDateTime(booking.Date, booking.Time); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot cancel booking less than 24 hours before scheduled time"})
	}

	// Perform the cancellation
	if err := h.repo.Cancel(c.Request().Context(), id, req.Reason); err != nil {
		h.logger.Err(fmt.Sprintf("Failed to cancel booking %d: %s", id, err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to cancel booking"})
	}

	h.logger.Info(fmt.Sprintf("Booking cancelled successfully: ID %d, Reason: %s", id, req.Reason))
	return c.JSON(http.StatusOK, map[string]string{
		"message": "booking cancelled successfully",
	})
}

// GenerateUpdateLink generates a secure update link for guest users
func (h *BookRideHandler) GenerateUpdateLink(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid booking ID"})
	}

	var req struct {
		Email string `json:"email" validate:"required,email"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	if err := validation.ValidateEmail(email); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Verify booking exists and belongs to this email
	booking, err := h.repo.GetByIDAndEmail(c.Request().Context(), id, email)
	if err != nil {
		// Don't reveal if booking exists for security
		return c.JSON(http.StatusOK, map[string]string{
			"message": "if the booking exists for this email, an update link has been sent",
		})
	}

	// Check if booking can be updated
	if booking.BookStatus == models.BookStatusCancelled || booking.BookStatus == models.BookStatusCompleted {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot generate update link for cancelled or completed booking"})
	}

	// Generate secure token
	token, err := h.authService.GenerateBookingUpdateToken(id, email)
	if err != nil {
		h.logger.Err(fmt.Sprintf("Failed to generate update token for booking %d: %s", id, err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate update link"})
	}

	h.logger.Info(fmt.Sprintf("Update token generated for booking %d, email %s", id, email))

	// Send email if email service is configured
	if h.emailService != nil {
		if err := h.emailService.SendBookingUpdateEmail(email, token, booking); err != nil {
			h.logger.Err(fmt.Sprintf("Failed to send update email to %s: %s", email, err.Error()))
			// Don't fail the request if email fails
			return c.JSON(http.StatusOK, map[string]interface{}{
				"message":      "update link generated (email service failed)",
				"update_token": token,
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"message": "if the booking exists for this email, an update link has been sent",
		})
	} else {
		// In development mode without email service, return the token
		h.logger.Warn("Email service not configured, returning update token in response")
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":      "update link generated (email service disabled)",
			"update_token": token,
		})
	}
}