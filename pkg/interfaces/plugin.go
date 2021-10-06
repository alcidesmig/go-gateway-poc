package interfaces

import (
	"net/http"
	errors "poc-gateway/pkg/errors"
)

type GenericGatewayPlugin interface {
	Setup() error
	Process(req *http.Request) *errors.GeneralError
	Close() error
}
