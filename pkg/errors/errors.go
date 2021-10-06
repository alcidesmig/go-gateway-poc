package pkg

import (
	"errors"
	"fmt"
)

var (
	ErrValidationError  = errors.New("validation error")
	ErrPermissionDenied = errors.New("permission denied")
)

type GeneralError struct {
	Info       string
	Error      error
	StatusCode int
}

func (w GeneralError) String() string {
	return fmt.Sprintf("%s: %v", w.Error.Error(), w.Info)
}

func (w GeneralError) JSON() []byte {
	return []byte(
		fmt.Sprintf(`{"message":"%s","error_type":"%s","status_code":"%d"}`,
			w.Info, w.Error.Error(), w.StatusCode))
}

func (w *GeneralError) SetStatusCode(statusCode int) {
	w.StatusCode = statusCode
}

func Error(errorType error, info string) *GeneralError {
	err := &GeneralError{Error: errorType, Info: info}
	switch errorType {
	case ErrValidationError:
		err.StatusCode = 400
	case ErrPermissionDenied:
		err.StatusCode = 403
	default:
		err.StatusCode = 500
	}
	return err
}
