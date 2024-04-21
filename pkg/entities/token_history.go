package entities

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type TokenHistory struct {
	gorm.Model
	Token     string `gorm:"index"`
	RevokedAt time.Time
}

type TokenHistoryRepository interface {
	Create(ctx context.Context, token *TokenHistory) error
	FindByToken(ctx context.Context, token string) (*TokenHistory, error)
}
