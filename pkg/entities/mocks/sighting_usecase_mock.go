// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/muhwyndhamhp/tigerhall-kittens/graph/model"
)

// SightingUsecase is an autogenerated mock type for the SightingUsecase type
type SightingUsecase struct {
	mock.Mock
}

// CreateSighting provides a mock function with given fields: ctx, sighting, userID
func (_m *SightingUsecase) CreateSighting(ctx context.Context, sighting *model.NewSighting, userID uint) (*model.Sighting, error) {
	ret := _m.Called(ctx, sighting, userID)

	var r0 *model.Sighting
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.NewSighting, uint) (*model.Sighting, error)); ok {
		return rf(ctx, sighting, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.NewSighting, uint) *model.Sighting); ok {
		r0 = rf(ctx, sighting, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Sighting)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.NewSighting, uint) error); ok {
		r1 = rf(ctx, sighting, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSightingsByTigerID provides a mock function with given fields: ctx, tigerID, page, pageSize
func (_m *SightingUsecase) GetSightingsByTigerID(ctx context.Context, tigerID uint, page int, pageSize int) ([]*model.Sighting, int, error) {
	ret := _m.Called(ctx, tigerID, page, pageSize)

	var r0 []*model.Sighting
	var r1 int
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, int, int) ([]*model.Sighting, int, error)); ok {
		return rf(ctx, tigerID, page, pageSize)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint, int, int) []*model.Sighting); ok {
		r0 = rf(ctx, tigerID, page, pageSize)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Sighting)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint, int, int) int); ok {
		r1 = rf(ctx, tigerID, page, pageSize)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(context.Context, uint, int, int) error); ok {
		r2 = rf(ctx, tigerID, page, pageSize)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// NewSightingUsecase creates a new instance of SightingUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSightingUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *SightingUsecase {
	mock := &SightingUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
