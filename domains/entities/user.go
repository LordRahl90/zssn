package entities

import "zssn/domains/users/store"

// User service entity for transporting between service and DB layer
type User struct {
	ID          string        `json:"id"`
	Email       string        `json:"email"` // let's keep email, zombie apocalypse shouldn't make us forget that :)
	Name        string        `json:"name"`
	Age         uint32        `json:"age"`
	Gender      string        `json:"gender"`
	Latitude    float64       `json:"latitude"`
	Longitude   float64       `json:"longitude"`
	Infected    bool          `json:"infected"`
	FlagMonitor []FlagMonitor `json:"flag_monitor"`
}

// FlagMonitor tracks user flagging details
type FlagMonitor struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	InfectedUserID string `json:"infected_user_id"`
	InfectedUser   User   `json:"infected_user"`
}

// ToUserDBEntity converts service entity to DB entity
func (u *User) ToUserDBEntity() *store.User {
	return &store.User{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Age:       u.Age,
		Gender:    store.GenderFromString(u.Gender),
		Latitude:  u.Latitude,
		Longitude: u.Longitude,
	}
}

// FromUserDBEntity returns a service entity from the db entity
func FromUserDBEntity(m *store.User) *User {
	if m == nil {
		return nil
	}
	u := &User{
		ID:        m.ID,
		Email:     m.Email,
		Name:      m.Name,
		Age:       m.Age,
		Gender:    m.Gender.String(),
		Latitude:  m.Latitude,
		Longitude: m.Longitude,
		Infected:  m.Infected,
	}
	for _, v := range m.FlagMonitor {
		u.FlagMonitor = append(u.FlagMonitor, FlagMonitor{
			ID:             v.ID,
			UserID:         v.UserID,
			InfectedUserID: v.InfectedUserID,
		})
	}

	return u
}
