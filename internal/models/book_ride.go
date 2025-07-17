package models

type BookRide struct {
	ID                 int64  `json:"id" db:"id"`
	UserID             *int64 `json:"user_id,omitempty" db:"user_id"`
	DriverID           *int64 `json:"driver_id,omitempty" db:"driver_id"`
	YourName           string `json:"your_name" db:"your_name"`
	Email              string `json:"email" db:"email"`
	PhoneNumber        string `json:"phone_number" db:"phone_number"`
	RideType           string `json:"ride_type" db:"ride_type"`
	PickupLocation     string `json:"pickup_location" db:"pickup_location"`
	DropoffLocation    string `json:"dropoff_location" db:"dropoff_location"`
	Date               string `json:"date" db:"date"`
	Time               string `json:"time" db:"time"`
	NumberOfPassengers int    `json:"number_of_passengers" db:"number_of_passengers"`
	NumberOfLuggage    int    `json:"number_of_luggage" db:"number_of_luggage"`
	AdditionalNotes    string `json:"additional_notes,omitempty" db:"additional_notes"`
	BookStatus         string `json:"book_status" db:"book_status"`
	RideStatus         string `json:"ride_status" db:"ride_status"`
	CreatedAt          string `json:"created_at" db:"created_at"`
	UpdatedAt          string `json:"updated_at" db:"updated_at"`
}

// UpdateBookRideRequest represents the request payload for updating a booking
type UpdateBookRideRequest struct {
	YourName           string `json:"your_name,omitempty"`
	PhoneNumber        string `json:"phone_number,omitempty"`
	RideType           string `json:"ride_type,omitempty"`
	PickupLocation     string `json:"pickup_location,omitempty"`
	DropoffLocation    string `json:"dropoff_location,omitempty"`
	Date               string `json:"date,omitempty"`
	Time               string `json:"time,omitempty"`
	NumberOfPassengers *int   `json:"number_of_passengers,omitempty"`
	NumberOfLuggage    *int   `json:"number_of_luggage,omitempty"`
	AdditionalNotes    string `json:"additional_notes,omitempty"`
}

// BookRide status constants
const (
	BookStatusPending   = "Pending"
	BookStatusAccepted  = "Accepted"
	BookStatusCancelled = "Cancelled"
	BookStatusCompleted = "Completed"

	RideStatusPending   = "Pending"
	RideStatusAssigned  = "Assigned"
	RideStatusInProgress = "In Progress"
	RideStatusCompleted = "Completed"
	RideStatusCancelled = "Cancelled"
)

// BookingAssignmentRequest represents the request to assign a booking to a driver
type BookingAssignmentRequest struct {
	DriverID int64  `json:"driver_id" validate:"required"`
	Notes    string `json:"notes,omitempty"`
}

// BookingListResponse represents the response for booking lists with driver info
type BookingListResponse struct {
	*BookRide
	DriverName  string `json:"driver_name,omitempty"`
	DriverEmail string `json:"driver_email,omitempty"`
}