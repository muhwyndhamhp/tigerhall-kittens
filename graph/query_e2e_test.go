package graph

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/muhwyndhamhp/tigerhall-kittens/graph/model"
	"github.com/stretchr/testify/assert"
)

func TestQuery_Tigers(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name string

		want    *model.TigerPagination
		wantErr error
	}{
		{
			name: "should return tigers and nil error",
			want: &model.TigerPagination{
				Tigers: []*model.Tiger{
					{
						ID:            1,
						Name:          "tiger-1",
						DateOfBirth:   now,
						LastSeen:      now,
						LastLatitude:  -7.550676,
						LastLongitude: 110.828316,
					},
				},
				Total: 1,
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, _, _ := Setup(t, now, false)

			res, err := r.Query().Tigers(context.Background(), 1, 10)

			wantJS, _ := json.Marshal(tc.want)
			resJS, _ := json.Marshal(res)

			assert.Equal(t, string(wantJS), string(resJS))
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestQuery_SightingByTiger(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name string

		want    *model.SightingsPagination
		wantErr error
	}{
		{
			name: "should return sightings and nil error",
			want: &model.SightingsPagination{
				Sightings: []*model.Sighting{
					{
						ID:        1,
						Date:      now,
						Latitude:  -7.550676,
						Longitude: 110.828316,
						TigerID:   1,
						UserID:    1,
						ImageURL:  new(string),
					},
				},
				Total: 1,
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, _, _ := Setup(t, now, false)

			res, err := r.Query().SightingByTiger(context.Background(), 1, 1, 10)

			wantJS, _ := json.Marshal(tc.want)
			resJS, _ := json.Marshal(res)

			assert.Equal(t, string(wantJS), string(resJS))
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestSighting_Tiger(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name string

		want    *model.Tiger
		wantErr error
	}{
		{
			name: "should return tiger and nil error",
			want: &model.Tiger{
				ID:            1,
				Name:          "tiger-1",
				DateOfBirth:   now,
				LastSeen:      now,
				LastLatitude:  -7.550676,
				LastLongitude: 110.828316,
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, _, _ := Setup(t, now, false)

			res, err := r.Sighting().Tiger(context.Background(), &model.Sighting{
				TigerID: 1,
			})

			wantJS, _ := json.Marshal(tc.want)
			resJS, _ := json.Marshal(res)

			assert.Equal(t, string(wantJS), string(resJS))
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestSighting_User(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name string

		want    *model.User
		wantErr error
	}{
		{
			name: "should return user and nil error",
			want: &model.User{
				ID:    1,
				Name:  "user-1",
				Email: "email-1@example.com",
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, _, _ := Setup(t, now, false)

			res, err := r.Sighting().User(context.Background(), &model.Sighting{
				UserID: 1,
			})

			wantJS, _ := json.Marshal(tc.want)
			resJS, _ := json.Marshal(res)

			assert.Equal(t, string(wantJS), string(resJS))
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestTiger_Sightings(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name string

		want    []*model.Sighting
		wantErr error
	}{
		{
			name: "should return sightings and nil error",
			want: []*model.Sighting{
				{
					ID:        1,
					Date:      now,
					Latitude:  -7.550676,
					Longitude: 110.828316,
					TigerID:   1,
					UserID:    1,
					ImageURL:  new(string),
				},
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, _, _ := Setup(t, now, false)

			res, err := r.Tiger().Sightings(context.Background(), &model.Tiger{
				ID: 1,
			})

			wantJS, _ := json.Marshal(tc.want)
			resJS, _ := json.Marshal(res)

			assert.Equal(t, string(wantJS), string(resJS))
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
