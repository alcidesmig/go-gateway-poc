package plugins

import (
	"net/http"
	"poc-gateway/pkg/interfaces"

	errors "poc-gateway/pkg/errors"

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

func (p OIDCPlugin) Process(req *http.Request) *errors.GeneralError {
	rawToken := req.Header.Get("Authorization")

	// Verify if the prefixTokenValue is present in the header
	if len(rawToken) < prefixTokenType ||
		rawToken[:prefixTokenType] != prefixTokenValue {
		return errors.Error(errors.ErrValidationError, "unknown token type")
	}

	// Skip the token prefix
	rawToken = rawToken[prefixTokenType:]

	// Parse and validate the JWT
	// (including token exp)
	parsedToken, err := jwt.Parse(rawToken, p.JWKs.Keyfunc)
	if err != nil {
		return errors.Error(errors.ErrValidationError, err.Error())
	}

	// Get JWT "aud" field
	audience, ok := parsedToken.Claims.(jwt.MapClaims)["aud"]
	if !ok {
		return errors.Error(errors.ErrValidationError, err.Error())
	}
	// Verifies if it match one of the allowed auds
	for _, aud := range p.AllowedAuds {
		if audience == aud {
			return nil
		}
	}

	// Didn't matched any allowed aud
	return errors.Error(errors.ErrPermissionDenied, err.Error())
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
