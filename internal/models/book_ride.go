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
}
