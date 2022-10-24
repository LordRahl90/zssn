package requests

import "fmt"

var (
	errInvalidName      = fmt.Errorf("invalid name")
	errInvalidEmail     = fmt.Errorf("invalid email")
	errInvalidAge       = fmt.Errorf("invalid age")
	errInvalidGender    = fmt.Errorf("invalid gender")
	errInvalidInventory = fmt.Errorf("invalid inventory")
)

// Survivor sample survivor request format
type Survivor struct {
	Email     string      `json:"email" form:"email"`
	Name      string      `json:"name" form:"name"`
	Age       uint32      `json:"age" form:"age"`
	Gender    string      `json:"gender" form:"gender"`
	Latitude  float64     `json:"latitude" form:"latitude"`
	Longitude float64     `json:"longitude" form:"longitude"`
	Inventory []Inventory `json:"inventories" validate:"required"`
	Token     string      `json:"token"`
}

// FlagUser request format for flagging infected users
type FlagUser struct {
	InfectedUserID string `json:"infected_user_id"`
}

// NewTokenRequest request format for requesting new tokens since no auths
type NewToken struct {
	Email string `json:"email"`
}

// UpdateLocation request format for updating user's location
type UpdateLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Validate makes sure that all important fields are provided
func (s *Survivor) Validate() error {
	switch {
	case s.Email == "":
		return errInvalidEmail
	case s.Name == "":
		return errInvalidName
	case s.Age <= 0:
		return errInvalidAge
	case s.Gender == "":
		return errInvalidGender
	case len(s.Inventory) == 0:
		return errInvalidInventory

	default:
		return nil
	}
}
