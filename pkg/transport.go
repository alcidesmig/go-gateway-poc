package pkg

import (
	"log"
	"net/http"
	"time"

	"github.com/beinan/fastid"
)

type TransportLayer struct{}

func logIn(requestID int64, request *http.Request) {
	log.Printf("[-> %d] %s %s\n", requestID, request.Method, request.URL.EscapedPath())
}

func logOut(requestID int64, st time.Time, response *http.Response) {
	log.Printf("[<- %d] %d - %s\n", requestID, response.StatusCode, time.Since(st))
}

func (t *TransportLayer) RoundTrip(request *http.Request) (*http.Response, error) {

	requestID := fastid.CommonConfig.GenInt64ID()
	start := time.Now()

	logIn(requestID, request)
	response, err := http.DefaultTransport.RoundTrip(request)
	logOut(requestID, start, response)

	return response, err
}
