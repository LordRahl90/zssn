package users

import "gorm.io/gorm"

type UserServiceImpl struct {
	DB *gorm.DB
}

func New(db *gorm.DB) (UserService, error) {
	return nil, nil
}
