package users

import (
	"time"
	"github.com/xesina/golang-realworld/pkg/types"
)

type User struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt types.NullTime `sql:"index"`
	Username  string         `gorm:"unique_index"`
	Email     string         `gorm:"unique_index"`
	Password  string
	Bio       types.NullString
	Image     types.NullString
}

type UserRepository interface {
	Find(id uint) (*User, error)
	Create(user *User) (*User, error)
	Update(user *User) error
}

type UserInteractor interface {
	Find(id uint) (*User, error)
	Register(username, email, password string) (*User, error)
	Update(user *User) error
}

type interactor struct {
	repo UserRepository
}

func NewUserInteractor(r UserRepository) UserInteractor {
	return &interactor{
		repo: r,
	}
}

func (i *interactor) Find(id uint) (*User, error) {
	u, err := i.repo.Find(id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (i *interactor) Register(username, email, password string) (*User, error) {
	// TODO: Encrypt password
	u, err := i.repo.Create(&User{
		Username: username,
		Email:    email,
		Password: password,
	})
	// TODO: Return appropriate error (maybe wrapped)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (i *interactor) Update(user *User) error {
	err := i.repo.Update(user)
	if err != nil {
		return err
	}
	return nil
}
