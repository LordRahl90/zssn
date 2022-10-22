package store

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	_ IUserStorage = (*MockUserStorage)(nil)

	errMockNotDefined = errors.New("mock function not defined")
	mockdDB           = make(map[string]*User)
)

// MockUserStorage returns a mocked storage object
type MockUserStorage struct {
	CreateFunc               func(ctx context.Context, user *User) error
	FlagUserFunc             func(ctx context.Context, id, infectedUser string) error
	UpdateInfectedStatusFunc func(ctx context.Context, id string) error
	FindFunc                 func(ctx context.Context, id string) (*User, error)
	FindByEmailFunc          func(ctx context.Context, email string) (*User, error)
	UpdateLocationFunc       func(ctx context.Context, id string, lat float64, long float64) error
}

// NewMockStore returns a new mock implementation of the functions
func NewMockStore() *MockUserStorage {
	return &MockUserStorage{
		CreateFunc: func(ctx context.Context, user *User) error {
			user.ID = uuid.NewString()
			mockdDB[user.ID] = user
			return nil
		},
		FlagUserFunc: func(ctx context.Context, id, infectedUserID string) error {
			v, ok := mockdDB[id]
			if !ok {
				return gorm.ErrRecordNotFound
			}
			infectedUser, ok := mockdDB[infectedUserID]
			if !ok {
				return gorm.ErrRecordNotFound
			}
			infectedUser.FlagMonitor = append(v.FlagMonitor, FlagMonitor{
				ID:             uuid.NewString(),
				UserID:         id,
				InfectedUserID: infectedUserID,
			})

			mockdDB[infectedUserID] = infectedUser
			return nil
		},
		FindFunc: func(ctx context.Context, id string) (*User, error) {
			v, ok := mockdDB[id]
			if !ok {
				return nil, gorm.ErrRecordNotFound
			}
			return v, nil
		},
		FindByEmailFunc: func(ctx context.Context, email string) (*User, error) {
			for _, v := range mockdDB {
				if v.Email == email {
					return v, nil
				}
			}
			return nil, gorm.ErrRecordNotFound
		},
		UpdateInfectedStatusFunc: func(ctx context.Context, id string) error {
			v, ok := mockdDB[id]
			if !ok {
				return gorm.ErrRecordNotFound
			}
			v.Infected = true
			mockdDB[id] = v
			return nil
		},
		UpdateLocationFunc: func(ctx context.Context, id string, lat, long float64) error {
			v, ok := mockdDB[id]
			if !ok {
				return gorm.ErrRecordNotFound
			}
			v.Latitude = lat
			v.Longitude = long
			mockdDB[id] = v
			return nil
		},
	}
}

// FlagUser implements IUserStorage
func (m *MockUserStorage) FlagUser(ctx context.Context, id string, infectedUser string) error {
	if m.FlagUserFunc == nil {
		return errMockNotDefined
	}
	return m.FlagUserFunc(ctx, id, infectedUser)
}

// UpdateInfectedStatus implements IUserStorage
func (m *MockUserStorage) UpdateInfectedStatus(ctx context.Context, id string) error {
	if m.UpdateInfectedStatusFunc == nil {
		return errMockNotDefined
	}
	return m.UpdateInfectedStatusFunc(ctx, id)
}

// Find implements IUserStorage
func (m *MockUserStorage) Find(ctx context.Context, id string) (*User, error) {
	if m.FindFunc == nil {
		return nil, errMockNotDefined
	}
	return m.FindFunc(ctx, id)
}

// FindByEmail implements IUserStorage
func (m *MockUserStorage) FindByEmail(ctx context.Context, email string) (*User, error) {
	if m.FindByEmailFunc == nil {
		return nil, errMockNotDefined
	}
	return m.FindByEmailFunc(ctx, email)
}

// UpdateLocation implements IUserStorage
func (m *MockUserStorage) UpdateLocation(ctx context.Context, id string, lat float64, long float64) error {
	if m.UpdateLocationFunc == nil {
		return errMockNotDefined
	}
	return m.UpdateLocationFunc(ctx, id, lat, long)
}

// Create mocked the create function
func (m *MockUserStorage) Create(ctx context.Context, user *User) error {
	if m.CreateFunc == nil {
		return errMockNotDefined
	}
	return m.CreateFunc(ctx, user)
}
