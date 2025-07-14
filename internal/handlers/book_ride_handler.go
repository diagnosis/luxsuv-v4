package handlers

import (
	"fmt"
	"github.com/diagnosis/luxsuv-v4/internal/logger"
	"github.com/diagnosis/luxsuv-v4/internal/models"
	"github.com/diagnosis/luxsuv-v4/internal/repository"
	"github.com/diagnosis/luxsuv-v4/internal/validation"
	"github.com/labstack/echo/v4"
	"net/http"
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

	//validate input
	if err := validation.ValidateBookRide(br); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	userID, ok := c.Get("user_id").(int64)
	if ok {
		br.UserID = &userID
	}
	if err := h.repo.Create(c.Request().Context(), br); err != nil {
		h.logger.Err(fmt.Sprintf("Error creating book ride: %s", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error creating book ride"})
	}
	h.logger.Info(fmt.Sprintf("Booking created successfully: ID %d", br.ID))
	return c.JSON(http.StatusCreated, br)
}

func (h *BookRideHandler) GetByEmail(c echo.Context) error {
	email := c.Param("email")
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
