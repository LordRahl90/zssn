package user

import "context"

type UserService interface {
	Create(ctx context.Context) error
}
