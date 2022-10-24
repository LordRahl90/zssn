package users

import (
	"context"

	"zssn/domains/entities"
	"zssn/domains/users/store"
)

var (
	_ IUserService = (*UserService)(nil)
)

type UserService struct {
	Storage store.IUserStorage
}

// New create a new user service object
func New(storage store.IUserStorage) (IUserService, error) {
	return &UserService{
		Storage: storage,
	}, nil
}

// Create converts the entities and creates a new record inside the database
func (u *UserService) Create(ctx context.Context, user *entities.User) error {
	dbEntity := user.ToUserDBEntity()
	if err := u.Storage.Create(ctx, dbEntity); err != nil {
		return err
	}
	user.ID = dbEntity.ID
	return nil
}

// Find finds an existing record in the database
func (u *UserService) Find(ctx context.Context, id string) (*entities.User, error) {
	res, err := u.Storage.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	// if res.Infected {
	// 	return nil, fmt.Errorf("user is infected")
	// }
	return entities.FromUserDBEntity(res), nil
}

// FindUsers implements IUserService
func (u *UserService) FindUsers(ctx context.Context, ids ...string) (map[string]*entities.User, error) {
	result := make(map[string]*entities.User)
	res, err := u.Storage.FindUsers(ctx, ids...)
	if err != nil {
		return nil, err
	}
	for k, v := range res {
		// if v.Infected {
		// 	return nil, fmt.Errorf("user  %s(%s) is infected", v.Name, v.Email)
		// }
		result[k] = entities.FromUserDBEntity(v)
	}
	return result, nil
}

// FindByEmail returns a user whose email matches the given email
func (u *UserService) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	res, err := u.Storage.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return entities.FromUserDBEntity(res), nil
}

// UpdateLocation updates user's location
func (u *UserService) UpdateLocation(ctx context.Context, id string, lat, long float64) error {
	return u.Storage.UpdateLocation(ctx, id, lat, long)
}

// FlagUser flags a user using the storage service and if the infected user has been flagged 3 times or more
// then the infection status of the user is updated
func (u *UserService) FlagUser(ctx context.Context, id, infectedUserID string) error {
	usr, err := u.Storage.Find(ctx, infectedUserID)
	if err != nil {
		return err
	}

	if len(usr.FlagMonitor)+1 >= 3 && !usr.Infected {
		// update the user as infected once we get to 3 and the user hasn't been flagged already
		if err := u.Storage.UpdateInfectedStatus(ctx, infectedUserID); err != nil {
			return err
		}
	}

	return u.Storage.FlagUser(ctx, id, infectedUserID)
}

// IsInfected return if a user is infected or not
func (u *UserService) IsInfected(ctx context.Context, id string) (bool, error) {
	res, err := u.Storage.Find(ctx, id)
	if err != nil {
		return false, err
	}
	return res.Infected, nil
}
