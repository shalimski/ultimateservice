package auth_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/shalimski/ultimateservice/business/sys/auth"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	const keyID = "71a5464a-1d6e-4841-b472-099aef886346"

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	authstore, err := auth.NewAuth(keyID, &keyStore{pk: privateKey})
	assert.NoError(t, err)

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "for test",
			Subject:   "d298109c-35b3-4c7f-a98d-91ca5541051a",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: []string{auth.RoleAdmin},
	}

	token, err := authstore.GenerateToken(claims)
	assert.NoError(t, err)

	parsedClaims, err := authstore.ValidateToken(token)
	assert.NoError(t, err)
	assert.ElementsMatch(t, claims.Roles, parsedClaims.Roles)
}

func TestAuthExpired(t *testing.T) {
	const keyID = "71a5464a-1d6e-4841-b472-099aef886346"

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	authstore, err := auth.NewAuth(keyID, &keyStore{pk: privateKey})
	assert.NoError(t, err)

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "for test",
			Subject:   "d298109c-35b3-4c7f-a98d-91ca5541051a",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(-2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: []string{auth.RoleAdmin},
	}

	token, err := authstore.GenerateToken(claims)
	assert.NoError(t, err)

	_, err = authstore.ValidateToken(token)
	assert.Error(t, err)
}

func TestClaims(t *testing.T) {
	claims := auth.Claims{
		Roles: []string{auth.RoleUser},
	}

	ctx := context.Background()

	ctx = auth.SetClaims(ctx, claims)

	claims, err := auth.GetClaims(ctx)
	assert.NoError(t, err)

	got := claims.Authorized("fake", auth.RoleAdmin)

	assert.False(t, got)

	got = claims.Authorized(auth.RoleUser)

	assert.True(t, got)
}

// =============================================================================

type keyStore struct {
	pk *rsa.PrivateKey
}

func (ks *keyStore) PrivateKey(kid string) (*rsa.PrivateKey, error) {
	return ks.pk, nil
}

func (ks *keyStore) PublicKey(kid string) (*rsa.PublicKey, error) {
	return &ks.pk.PublicKey, nil
}
