package user

import (
	"github.com/muhwyndhamhp/tigerhall-kittens/graph/model"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
)

type usecase struct {
	repo entities.UserRepository
}

// GetUserByID implements entities.UserUsecase.
func (u *usecase) GetUserByID(id uint) (*model.User, error) {
	usr, err := u.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:    usr.ID,
		Name:  usr.Name,
		Email: usr.Email,
	}, nil
}

// CreateUser implements entities.UserUsecase.
func (u *usecase) CreateUser(name, email, password string) (string, error) {
	h, err := entities.HashPassword(password)
	if err != nil {
		return "", err
	}

	usr := entities.User{
		Name:         name,
		Email:        email,
		PasswordHash: h,
	}
	err = u.repo.Create(&usr)
	if err != nil {
		return "", err
	}

	token, err := usr.GenerateToken()
	if err != nil {
		return "", err
	}

	return token, nil
}

// Login implements entities.UserUsecase.
func (u *usecase) Login(email string, password string) (string, error) {
	usr, err := u.repo.FindByEmail(email)
	if err != nil {
		return "", err
	}

	err = usr.ValidatePassword(password)
	if err != nil {
		return "", err
	}

	token, err := usr.GenerateToken()
	if err != nil {
		return "", err
	}

	return token, nil
}

func NewUserUsecase(repo entities.UserRepository) entities.UserUsecase {
	return &usecase{repo}
}
