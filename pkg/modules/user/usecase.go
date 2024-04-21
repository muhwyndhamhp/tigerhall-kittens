package user

import (
	"context"
	"errors"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/muhwyndhamhp/tigerhall-kittens/graph/model"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
)

type usecase struct {
	repo      entities.UserRepository
	tokenRepo entities.TokenHistoryRepository
}

// RefreshToken implements entities.UserUsecase.
func (u *usecase) RefreshToken(ctx context.Context, token string) (string, error) {
	t, _ := u.tokenRepo.FindByToken(ctx, token)
	if t != nil {
		return "", entities.ErrTokenAlreadyInvalidated
	}

	tu, err := entities.ParseToken(token)
	if err != nil {
		return "", err
	}

	usr, err := u.repo.FindByID(ctx, tu.ID)
	if err != nil || usr == nil {
		return "", entities.ErrUserNotFound
	}

	newToken, err := usr.GenerateToken()
	if err != nil {
		return "", err
	}

	go func() {
		err = u.tokenRepo.Create(context.Background(), &entities.TokenHistory{Token: token, RevokedAt: time.Now()})
		if err != nil {
			log.Error(errors.New("error creating token history"))
		}
	}()

	return newToken, nil
}

// GetUserByID implements entities.UserUsecase.
func (u *usecase) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	usr, err := u.repo.FindByID(ctx, id)
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
func (u *usecase) CreateUser(ctx context.Context, usr *model.NewUser) (string, error) {
	h, err := entities.HashPassword(usr.Password)
	if err != nil {
		return "", err
	}

	existingUser, _ := u.repo.FindByEmail(ctx, usr.Email)
	if existingUser != nil {
		return "", entities.ErrUserAlreadyExists
	}

	newUsr := entities.User{
		Name:         usr.Name,
		Email:        usr.Email,
		PasswordHash: h,
	}
	err = u.repo.Create(ctx, &newUsr)
	if err != nil {
		return "", err
	}

	token, err := newUsr.GenerateToken()
	if err != nil {
		return "", err
	}

	return token, nil
}

// Login implements entities.UserUsecase.
func (u *usecase) Login(ctx context.Context, email string, password string) (string, error) {
	usr, err := u.repo.FindByEmail(ctx, email)
	if err != nil || usr == nil {
		return "", entities.ErrUserNotFound
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

func NewUserUsecase(r entities.UserRepository, tr entities.TokenHistoryRepository) entities.UserUsecase {
	return &usecase{r, tr}
}
