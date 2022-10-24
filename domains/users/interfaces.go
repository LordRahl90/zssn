package users

import (
	"context"
	"zssn/domains/entities"
)

// IUserService interface describing the contracts between the services
type IUserService interface {
	Create(ctx context.Context, user *entities.User) error
	Find(ctx context.Context, id string) (*entities.User, error)
	FindUsers(ctx context.Context, ids ...string) (map[string]*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	UpdateLocation(ctx context.Context, id string, lat, long float64) error
	FlagUser(ctx context.Context, id, infectedUser string) error
	IsInfected(ctx context.Context, id string) (bool, error)
}
