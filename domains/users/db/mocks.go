package db

import (
	"context"
	"errors"
)

var (
	_ IUserStorage = (*MockUserStorage)(nil)

	errMockNotDefined = errors.New("mock function not defined")
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
