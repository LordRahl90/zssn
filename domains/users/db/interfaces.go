package db

import (
	"context"
)

// IUserStorage interface describing the expectations for storage engine
type IUserStorage interface {
	Create(ctx context.Context, user *User) error
	Find(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	UpdateLocation(ctx context.Context, id string, lat, long float64) error
	FlagUser(ctx context.Context, id, infectedUser string) error
	UpdateInfectedStatus(ctx context.Context, id string) error
}
