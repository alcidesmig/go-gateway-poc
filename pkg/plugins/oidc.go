package plugins

import (
	"errors"
	"net/http"
	"poc-gateway/pkg/interfaces"

	"github.com/MicahParks/keyfunc"
	jwt "github.com/dgrijalva/jwt-go"
)

type OIDCPlugin struct {
	interfaces.GenericGatewayPlugin
	JWKS *keyfunc.JWKs
}

const prefixTokenType int = len("Bearer")

func (p OIDCPlugin) Process(req *http.Request) (int, error) {
	rawToken := req.Header.Get("Authorization")

	if len(rawToken) < prefixTokenType ||
		rawToken[:prefixTokenType] != "Bearer" {
		return 400, errors.New("unknown token type")
	}

	rawToken = rawToken[prefixTokenType:]

	_, err := jwt.Parse(rawToken, p.JWKS.KeyfuncLegacy)
	if err != nil {
		return 400, err
	}
	return 200, nil
}

func (p OIDCPlugin) Setup() error {
	jwks, err := keyfunc.Get("https://id.magalu.com/oauth/certs", keyfunc.Options{
		Client: &http.Client{
			Timeout: 0,
		},
	})
	p.JWKS = jwks
	return err
}

func (p OIDCPlugin) Close() error {
	return nil
}
