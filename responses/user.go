package responses

import "zssn/domains/entities"

// User sample survivor request format
type User struct {
	ID        string       `json:"id"`
	Email     string       `json:"email"`
	Name      string       `json:"name"`
	Age       uint32       `json:"age"`
	Gender    string       `json:"gender"`
	Latitude  float64      `json:"latitude"`
	Longitude float64      `json:"longitude"`
	Inventory []*Inventory `json:"inventories,omitempty"`
	Token     string       `json:"token,omitempty"`
}

// FromUserEntity converts user entity to response user object
func FromUserEntity(u *entities.User, token string) *User {
	return &User{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Age:       u.Age,
		Gender:    u.Gender,
		Latitude:  u.Latitude,
		Longitude: u.Longitude,
		Token:     token,
	}
}
