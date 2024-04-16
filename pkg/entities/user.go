package entities

import (
	"errors"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/muhwyndhamhp/tigerhall-kittens/graph/model"
	"github.com/muhwyndhamhp/tigerhall-kittens/utils/config"
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
	CreateUser(name, email, password string) (string, error)
	Login(email, password string) (string, error)
	GetUserByID(id uint) (*model.User, error)
}

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
}

var ErrUserByCtxNotFound = errors.New("ErrUserByCtxNotFound: user not found in context")

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

// JWT Implementation
var secretKey = []byte(config.Get(config.JWT_SECRET))

func (u *User) GenerateToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	exp, err := strconv.Atoi(config.Get(config.JWT_EXPIRY_DURATION))
	if err != nil {
		return "", err
	}

	claims["id"] = u.ID
	claims["username"] = u.Name
	claims["email"] = u.Email
	claims["exp"] = time.Now().Add(time.Second * time.Duration(exp)).Unix()
	ts, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return ts, nil
}

func ParseToken(tokenString string) (*User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
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

func RefreshToken(tokenString string) (string, error) {
	u, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	return u.GenerateToken()
}
