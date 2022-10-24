package mocks

import (
	"context"
	"errors"

	"zssn/domains/entities"
	"zssn/domains/users"
)

var (
	_ users.IUserService = (*MockUserService)(nil)

	errMockNotDefined = errors.New("mock not initialized")
)

// MockUserService mock user service
type MockUserService struct {
	CreateFunc         func(ctx context.Context, user *entities.User) error
	FindFunc           func(ctx context.Context, id string) (*entities.User, error)
	FindByEmailFunc    func(ctx context.Context, email string) (*entities.User, error)
	FindUsersFunc      func(ctx context.Context, ids ...string) (map[string]*entities.User, error)
	FlagUserFunc       func(ctx context.Context, id string, infectedUser string) error
	IsInfectedFunc     func(ctx context.Context, id string) (bool, error)
	UpdateLocationFunc func(ctx context.Context, id string, lat float64, long float64) error
}

// NewUserMock returns a legit user service using mocked db
func NewUserMock() (users.IUserService, error) {
	return users.New(users.NewMockStore())
}

// Create implements users.IUserService
func (m *MockUserService) Create(ctx context.Context, user *entities.User) error {
	if m.CreateFunc == nil {
		return errMockNotDefined
	}
	return m.CreateFunc(ctx, user)
}

// Find implements users.IUserService
func (m *MockUserService) Find(ctx context.Context, id string) (*entities.User, error) {
	if m.FindFunc == nil {
		return nil, errMockNotDefined
	}
	return m.FindFunc(ctx, id)
}

// FindByEmail implements users.IUserService
func (m *MockUserService) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	if m.FindByEmailFunc == nil {
		return nil, errMockNotDefined
	}

	return m.FindByEmailFunc(ctx, email)
}

// FindUsers implements users.IUserService
func (m *MockUserService) FindUsers(ctx context.Context, ids ...string) (map[string]*entities.User, error) {
	if m.FindUsersFunc == nil {
		return nil, errMockNotDefined
	}
	return m.FindUsersFunc(ctx, ids...)
}

// FlagUser implements users.IUserService
func (m *MockUserService) FlagUser(ctx context.Context, id string, infectedUser string) error {
	if m.FlagUserFunc == nil {
		return errMockNotDefined
	}
	return m.FlagUserFunc(ctx, id, infectedUser)
}

// IsInfected implements users.IUserService
func (m *MockUserService) IsInfected(ctx context.Context, id string) (bool, error) {
	if m.IsInfectedFunc == nil {
		return false, errMockNotDefined
	}
	return m.IsInfectedFunc(ctx, id)
}

// UpdateLocation implements users.IUserService
func (m *MockUserService) UpdateLocation(ctx context.Context, id string, lat float64, long float64) error {
	if m.UpdateLocationFunc == nil {
		return errMockNotDefined
	}
	return m.UpdateLocationFunc(ctx, id, lat, long)
}
