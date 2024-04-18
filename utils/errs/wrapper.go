package errs

import (
	"errors"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

type ServiceError struct {
	ErrorCode string
	Err       error
}

func (e ServiceError) Error() string {
	return e.Err.Error()
}

func RespError(err error) error {
	errList := gqlerror.List{}

	var serviceErr ServiceError
	if errors.As(err, &serviceErr) {
		errList = append(errList, &gqlerror.Error{
			Message: serviceErr.Error(),
			Extensions: map[string]interface{}{
				"code": serviceErr.ErrorCode,
			},
		})
	} else {
		errList = append(errList, &gqlerror.Error{
			Message: err.Error(),
		})
	}

	return errList
}
