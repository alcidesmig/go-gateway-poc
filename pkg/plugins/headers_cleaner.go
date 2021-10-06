package plugins

import (
	"net/http"
	errors "poc-gateway/pkg/errors"
	"poc-gateway/pkg/interfaces"
	"strings"
)

type HeaderCleanerPlugin struct {
	interfaces.GenericGatewayPlugin
	AllowedHeaders []string
}

func (p HeaderCleanerPlugin) Process(req *http.Request) *errors.GeneralError {

	// Iterate over existent headers comparing
	// with allowed headers, excluding those who
	// are not in the allowed headers list
RangeHeaders:
	for header := range req.Header {
		stdHeader := strings.ToLower(header)
		for _, allowed := range p.AllowedHeaders {
			if stdHeader == allowed {
				continue RangeHeaders
			}
		}
		req.Header.Del(header)
	}
	return nil
}

func (p *HeaderCleanerPlugin) Setup() error {
	// Set all headers name to lowercase
	for index, header := range p.AllowedHeaders {
		p.AllowedHeaders[index] = strings.ToLower(header)
	}
	return nil
}

func (p HeaderCleanerPlugin) Close() error {
	return nil
}
