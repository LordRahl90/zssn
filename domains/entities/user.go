package entities

import "gorm.io/gorm"

// User service entity for transporting between service and DB layer
type User struct {
	ID          string        `json:"id"`
	Email       string        `json:"email"` // let's keep email, zombie apocalypse shouldn't make us forget that :)
	Name        string        `json:"name"`
	Age         uint32        `json:"age"`
	Gender      string        `json:"gender"`
	Latitude    float64       `json:"latitude"`
	Longitude   float64       `json:"longitude"`
	FlagMonitor []FlagMonitor `json:"flag_monitor"`
	Infected    bool          `json:"infected"`
	gorm.Model
}

// FlagMonitor tracks user flagging
type FlagMonitor struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	InfectedUserID string `json:"infected_user_id"`
	InfectedUser   User   `json:"infected_user"`
	gorm.Model
}
