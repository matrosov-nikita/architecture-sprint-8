package token_verifier

import (
	"context"
	"errors"
	"fmt"

	"os"

	"github.com/Nerzal/gocloak/v13"
	"github.com/golang-jwt/jwt/v5"
)

const (
	requiredRole = "prothetic_user"
)

var (
	ErrRoleAccessDenied error = errors.New("access denied for user role")
)

type TokenVerifier struct {
	keycloakClient *gocloak.GoCloak
	keyCloakRealm  string
}

func newTokenVerifier(
	keycloakClient *gocloak.GoCloak,
	keyCloakRealm string,
) *TokenVerifier {
	return &TokenVerifier{
		keycloakClient: keycloakClient,
		keyCloakRealm:  keyCloakRealm,
	}
}

func MustCreateTokenVerifier() *TokenVerifier {
	keyCloakURL := os.Getenv("APP_KEYCLOAK_URL")
	keyCloakRealm := os.Getenv("APP_KEYCLOAK_REALM")

	client := gocloak.NewClient(keyCloakURL)
	return newTokenVerifier(client, keyCloakRealm)
}

func (v *TokenVerifier) VerifyToken(ctx context.Context, tokenString string) error {
	decodedToken, claims, err := v.keycloakClient.DecodeAccessToken(ctx, tokenString, v.keyCloakRealm)
	if err != nil {
		return fmt.Errorf("decode access token: %v", err)
	}

	if !decodedToken.Valid || claims == nil {
		return errors.New("access token is not valid")
	}

	if !hasRoleInClaims(*claims, requiredRole) {
		return ErrRoleAccessDenied
	}

	return nil
}

func hasRoleInClaims(claims jwt.MapClaims, role string) bool {
	if realmAccess, ok := claims["realm_access"].(map[string]interface{}); ok {
		if roles, ok := realmAccess["roles"].([]interface{}); ok {
			for _, r := range roles {
				if r == role {
					return true
				}
			}
		}
	}
	return false
}
