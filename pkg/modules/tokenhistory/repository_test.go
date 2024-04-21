package tokenhistory

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/muhwyndhamhp/tigerhall-kittens/db"
	"github.com/muhwyndhamhp/tigerhall-kittens/pkg/entities"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRepository_Create(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string

		token   *entities.TokenHistory
		want    *entities.TokenHistory
		wantErr error
	}{
		{
			name: "should return token history with id 2 and nil error",
			token: &entities.TokenHistory{
				Token:     "token-2",
				RevokedAt: now,
			},
			want: &entities.TokenHistory{
				Model: gorm.Model{
					ID: 2,
				},
				Token:     "token-2",
				RevokedAt: now,
			},
			wantErr: nil,
		},

		{
			name: "should return error given token already exists",
			token: &entities.TokenHistory{
				Model: gorm.Model{
					ID: 1,
				},
				Token:     "token-1",
				RevokedAt: now,
			},
			want:    nil,
			wantErr: errors.New("UNIQUE constraint failed: token_histories.id"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := db.GetTestDB()
			repo := NewTokenHistoryRepository(d)
			SeedTokenHistory(d, now)

			err := repo.Create(context.Background(), tc.token)

			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr.Error(), err.Error())
			}

			if tc.want != nil {
				// assert every field of token history
				assert.Equal(t, tc.want.ID, tc.token.ID)
				assert.Equal(t, tc.want.Token, tc.token.Token)
				assert.Equal(t, tc.want.RevokedAt.Unix(), tc.token.RevokedAt.Unix())
			}
		})
	}
}

func TestRepository_FindByToken(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string

		token   string
		want    *entities.TokenHistory
		wantErr error
	}{
		{
			name:  "should return token history and nil error",
			token: "token-1",
			want: &entities.TokenHistory{
				Model: gorm.Model{
					ID: 1,
				},
				Token:     "token-1",
				RevokedAt: now,
			},
			wantErr: nil,
		},
		{
			name:    "should return nil and error given token not found",
			token:   "token-2",
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := db.GetTestDB()
			repo := NewTokenHistoryRepository(d)
			SeedTokenHistory(d, now)

			res, err := repo.FindByToken(context.Background(), tc.token)

			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr.Error(), err.Error())
			}

			if tc.want != nil {
				// assert every field of token history
				assert.Equal(t, tc.want.ID, res.ID)
				assert.Equal(t, tc.want.Token, res.Token)
				assert.Equal(t, tc.want.RevokedAt.Unix(), res.RevokedAt.Unix())
			}
		})
	}
}

func SeedTokenHistory(d *gorm.DB, now time.Time) {
	err := d.AutoMigrate(
		&entities.Tiger{},
		&entities.Sighting{},
		&entities.User{},
		&entities.TokenHistory{},
	)
	if err != nil {
		panic(err)
	}

	err = d.Create(&entities.TokenHistory{
		Token:     "token-1",
		RevokedAt: now,
	}).Error
	if err != nil {
		panic(err)
	}
}
