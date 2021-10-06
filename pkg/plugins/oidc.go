package plugins

import (
	"errors"
	"net/http"
	"poc-gateway/pkg/interfaces"

	"github.com/MicahParks/keyfunc"
	jwt "github.com/golang-jwt/jwt/v4"
)

type OIDCPlugin struct {
	interfaces.GenericGatewayPlugin
	JWKs        *keyfunc.JWKs
	IDPUrl      string
	AllowedAuds []string
}

const prefixTokenValue string = "Bearer "
const prefixTokenType int = len(prefixTokenValue)

func (p OIDCPlugin) Process(req *http.Request) (int, error) {
	rawToken := req.Header.Get("Authorization")

	// Verify if the prefixTokenValue is present in the header
	if len(rawToken) < prefixTokenType ||
		rawToken[:prefixTokenType] != prefixTokenValue {
		return 400, errors.New("unknown token type")
	}

	// Skip the token prefix
	rawToken = rawToken[prefixTokenType:]

	// Parse and validate the JWT
	// (including token exp)
	parsedToken, err := jwt.Parse(rawToken, p.JWKs.Keyfunc)
	if err != nil {
		return 400, err
	}

	// Get JWT "aud" field
	audience, ok := parsedToken.Claims.(jwt.MapClaims)["aud"]
	if !ok {
		return 400, err
	}
	// Verifies if it match one of the allowed auds
	for _, aud := range p.AllowedAuds {
		if audience == aud {
			return 200, nil
		}
	}

	// Didn't matched any allowed aud
	return 400, err
}

func (p *OIDCPlugin) Setup() error {
	jwks, err := keyfunc.Get(p.IDPUrl, keyfunc.Options{
		Client: &http.Client{
			// 0 timeout => no timeout
			Timeout: 0,
		},
	})
	p.JWKs = jwks
	return err
}

func (p OIDCPlugin) Close() error {
	return nil
}
