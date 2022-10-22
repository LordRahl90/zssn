package db

import (
	"gorm.io/gorm"
)

// Gender using this instead of enum to represent gender
type Gender int

const (
	// Male male gender
	GenderMale Gender = iota + 1
	// Female female gender
	GenderFemale
)

// User contains the user db entities
type User struct {
	ID          string        `json:"id" gorm:"primaryKey"`
	Email       string        `json:"email" gorm:"size:50;uniqueIndex"` // let's keep email, zombie apocalypse shouldn't make us forget that :)
	Name        string        `json:"name"`
	Age         uint32        `json:"age"`
	Gender      Gender        `json:"gender"`
	Latitude    float64       `json:"latitude"`
	Longitude   float64       `json:"longitude"`
	FlagMonitor []FlagMonitor `json:"flag_monitor" gorm:"foreignKey:InfectedUserID"`
	Infected    bool          `json:"infected"`
	Token       string        `json:"token" gorm:"-"`
	gorm.Model
}

// FlagMonitor tracks user flagging
// NB bad practice as this should have it's own separate flow,
// but for time I will combine both concerns
type FlagMonitor struct {
	ID             string `json:"id,omitempty" gorm:"primaryKey"`
	UserID         string `json:"user_id,omitempty" gorm:"size:50; index:idx_flagged,unique; "`
	User           *User  `json:"user" gorm:"foreignKey:UserID"`
	InfectedUserID string `json:"infected_user_id,omitempty" gorm:"size:50; index:idx_flagged,unique"`
	InfectedUser   *User  `json:"infected_user,omitempty" gorm:"foreignKey:InfectedUserID"`
	gorm.Model
}
