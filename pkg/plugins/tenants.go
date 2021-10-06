package plugins

import (
	"net/http"
	"poc-gateway/pkg/interfaces"

	errors "poc-gateway/pkg/errors"

	jwt "github.com/dgrijalva/jwt-go"
)

type TenantsPlugin struct {
	interfaces.GenericGatewayPlugin
}

type MagaluToken struct {
	jwt.StandardClaims
	Audience    interface{} `json:"aud,omitempty"` // keycloak aud can be []string or string
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	Tenants []struct {
		UUID       string `json:"uuid"`
		Type       string `json:"type"`
		InternalID string `json:"internal_id"`
	} `json:"tenants"`
}

const prefixTokenTypeBearer int = len("Bearer ")

func (p TenantsPlugin) Process(req *http.Request) *errors.GeneralError {
	rawToken := req.Header.Get("Authorization")
	specifiedTenant := req.Header.Get("x-tenant-id")

	// Verify if one x-tenant-id header was sent
	if specifiedTenant == "" {
		return errors.Error(errors.ErrValidationError, "one tenant id header need to be specified using x-tenant-id header")
	}

	// Assume one well formed JWT in Authorization header after "Bearer" prefix
	// Because of oidc.go plugin validation
	rawToken = rawToken[prefixTokenTypeBearer:]

	// Parse the JWK without verifying it
	processedToken, _, err := new(jwt.Parser).ParseUnverified(rawToken, &MagaluToken{})
	if err != nil {
		return errors.Error(errors.ErrValidationError, "invalid token")
	}
	// Cast to MagaluToken
	parsedToken := processedToken.Claims.(*MagaluToken)

	for _, tenant := range parsedToken.Tenants {
		if tenant.UUID == specifiedTenant {

			// Guarantee that only the allowed values are
			// present in the headers
			req.Header.Del("x-tenant-internal-id")
			req.Header.Del("x-tenant-type")

			// Insert new headers for the back-end
			req.Header.Add("x-tenant-internal-id", tenant.InternalID)
			req.Header.Add("x-tenant-type", tenant.Type)
			return nil
		}
	}

	// All the existent tenants in the JWT were checked
	// and none of them is valid
	return errors.Error(errors.ErrPermissionDenied, "you don't have access to the specified tenant")
}

func (p TenantsPlugin) Setup() error {
	return nil
}

func (p TenantsPlugin) Close() error {
	return nil
}
