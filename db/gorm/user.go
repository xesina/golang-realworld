package gorm

import (
	"github.com/jinzhu/gorm"
	"github.com/xesina/golang-realworld/users"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) users.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Find(id uint) (*users.User, error) {
	user := new(users.User)
	if err := r.db.First(user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Create(u *users.User) (*users.User, error) {
	if err := r.db.Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepository) Update(u *users.User) error {
	if err := r.db.Model(u).Update(u).Error; err != nil {
		return err
	}
	return nil
}
