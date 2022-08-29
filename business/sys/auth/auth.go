package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

// KeyLookup is a method set for JWT use
type KeyLookup interface {
	PrivateKey(kid string) (*rsa.PrivateKey, error)
	PublicKey(kid string) (*rsa.PublicKey, error)
}

type Auth struct {
	activeKID string
	keyLookup KeyLookup
	method    jwt.SigningMethod
	keyFunc   func(t *jwt.Token) (any, error)
	parser    jwt.Parser
}

func NewAuth(activeKID string, keyLookup KeyLookup) (*Auth, error) {
	_, err := keyLookup.PrivateKey(activeKID)
	if err != nil {
		return nil, fmt.Errorf("active KID not exist in store: %w", err)
	}

	method := jwt.GetSigningMethod("RS256")

	keyFunc := func(t *jwt.Token) (any, error) {
		kid, ok := t.Header["kid"]
		if !ok {
			return nil, errors.New("missing kid key in token header")
		}

		kidStr, ok := kid.(string)
		if !ok {
			return nil, errors.New("kid key must be a string")
		}

		return keyLookup.PublicKey(kidStr)
	}

	parser := jwt.Parser{
		ValidMethods: []string{method.Alg()},
	}

	auth := Auth{
		activeKID: activeKID,
		keyLookup: keyLookup,
		method:    method,
		keyFunc:   keyFunc,
		parser:    parser,
	}

	return &auth, nil
}

func (a *Auth) GenerateToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(a.method, claims)
	token.Header["kid"] = a.activeKID

	privateKey, err := a.keyLookup.PrivateKey(a.activeKID)
	if err != nil {
		return "", fmt.Errorf("private key lookup failed: %w", err)
	}

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("signing token failed: %w", err)
	}

	return tokenString, nil
}

func (a *Auth) ValidateToken(tokenStr string) (Claims, error) {
	var claims Claims

	token, err := a.parser.ParseWithClaims(tokenStr, &claims, a.keyFunc)
	if err != nil {
		return Claims{}, fmt.Errorf("parsing token: %w", err)
	}

	if !token.Valid {
		return Claims{}, errors.New("invalid token")
	}

	return claims, nil
}
