package interfaces

import "net/http"

type GenericGatewayPlugin interface {
	Setup() error
	Process(req *http.Request) (int, error)
	Close() error
}
