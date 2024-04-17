package entities

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/muhwyndhamhp/tigerhall-kittens/graph/model"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/config"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/timex"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

type UserUsecase interface {
	CreateUser(ctx context.Context, usr *model.NewUser) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
	RefreshToken(ctx context.Context, token string) (string, error)
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id uint) (*User, error)
}

var (
	ErrUserByCtxNotFound = errors.New("ErrUserByCtxNotFound: user not found in context")
	ErrUserAlreadyExists = errors.New("ErrUserAlreadyExists: user already exists")
)

// Password Hashing Implementation
func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (u *User) ValidatePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return err
	}

	return nil
}

func GetSecretKey() []byte {
	secretStr := config.Get(config.JWT_SECRET)
	if secretStr == "" {
		secretStr = "MuhWyndham-TigerHall-Kittens-Test"
	}
	return []byte(secretStr)
}

// JWT Implementation
func (u *User) GenerateToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	expStr := config.Get(config.JWT_EXPIRY_DURATION)
	if expStr == "" {
		expStr = "86400" // 24 hours
	}

	exp, err := strconv.Atoi(expStr)
	if err != nil {
		return "", err
	}

	claims["id"] = u.ID
	claims["username"] = u.Name
	claims["email"] = u.Email
	claims["exp"] = timex.Now().Add(time.Second * time.Duration(exp)).Unix()
	ts, err := token.SignedString(GetSecretKey())
	if err != nil {
		return "", err
	}

	return ts, nil
}

func ParseToken(tokenString string) (*User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return GetSecretKey(), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := &User{
			Model: gorm.Model{ID: uint(claims["id"].(float64))},
			Name:  claims["username"].(string),
			Email: claims["email"].(string),
		}
		return user, nil
	} else {
		return nil, err
	}
}
