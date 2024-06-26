// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	bytes "bytes"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// S3ClientInterface is an autogenerated mock type for the S3ClientInterface type
type S3ClientInterface struct {
	mock.Mock
}

// UploadImage provides a mock function with given fields: ctx, r, filename, contentType, size
func (_m *S3ClientInterface) UploadImage(ctx context.Context, r *bytes.Reader, filename string, contentType string, size int64) (string, error) {
	ret := _m.Called(ctx, r, filename, contentType, size)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *bytes.Reader, string, string, int64) (string, error)); ok {
		return rf(ctx, r, filename, contentType, size)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *bytes.Reader, string, string, int64) string); ok {
		r0 = rf(ctx, r, filename, contentType, size)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *bytes.Reader, string, string, int64) error); ok {
		r1 = rf(ctx, r, filename, contentType, size)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewS3ClientInterface creates a new instance of S3ClientInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewS3ClientInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *S3ClientInterface {
	mock := &S3ClientInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
