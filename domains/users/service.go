package users

import (
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func New(db *gorm.DB) (*IUserService, error) {
	return nil, nil
}
