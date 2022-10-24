package store

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserStorage user storage implementation
type UserStorage struct {
	DB *gorm.DB
}

// New creates a new instance of user storage with the given db connection
func New(db *gorm.DB) (IUserStorage, error) {
	if err := db.AutoMigrate(&User{}, &FlagMonitor{}); err != nil {
		return nil, err
	}
	return &UserStorage{
		DB: db,
	}, nil
}

// Create creates a new user record
func (u *UserStorage) Create(ctx context.Context, user *User) error {
	user.ID = uuid.NewString()
	return u.DB.Create(&user).Error
}

// FlagUser creates a new flag record against the infected user
func (u *UserStorage) FlagUser(ctx context.Context, id string, infectedUser string) error {
	// Not sure if user can flag themselves as infected
	if id == infectedUser {
		return nil
	}
	f := FlagMonitor{
		ID:             uuid.NewString(),
		UserID:         id,
		InfectedUserID: infectedUser,
	}
	return u.DB.Create(&f).Error
}

// UpdateInfectedStatus implements IUserStorage
func (u *UserStorage) UpdateInfectedStatus(ctx context.Context, id string) error {
	return u.DB.Model(&User{}).Where("id = ?", id).Update("infected", true).Error
}

// Find implements IUserStorage
func (u *UserStorage) Find(ctx context.Context, id string) (*User, error) {
	var user *User
	err := u.DB.Model(&User{}).Preload("FlagMonitor").Where("id = ?", id).First(&user).Error
	return user, err
}

// FindUsers finds array of users with the given ID
func (u *UserStorage) FindUsers(ctx context.Context, ids ...string) (map[string]*User, error) {
	var (
		users  []*User
		result = make(map[string]*User)
	)
	err := u.DB.Model(&User{}).Preload("FlagMonitor").Where("id IN (?)", ids).Find(&users).Error
	if err != nil {
		return nil, err
	}

	for _, v := range users {
		result[v.ID] = v
	}
	return result, nil
}

// FindByEmail implements IUserStorage
func (u *UserStorage) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user *User
	err := u.DB.Preload("FlagMonitor").Where("email = ?", email).Debug().First(&user).Error
	return user, err
}

// UpdateLocation implements IUserStorage
func (u *UserStorage) UpdateLocation(ctx context.Context, id string, lat float64, long float64) error {
	d := map[string]interface{}{
		"latitude":  lat,
		"longitude": long,
	}
	return u.DB.Model(&User{}).Where("id = ?", id).Updates(d).Error
}
