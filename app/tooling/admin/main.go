package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/shalimski/ultimateservice/business/data/schema"
	"github.com/shalimski/ultimateservice/business/sys/database"
	"github.com/shalimski/ultimateservice/internal/config"
)

var bits = *flag.Uint("bits", 2048, "bit size for private key generation")

func main() {
	flag.Parse()

	err := migrate()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func genKey() error {
	flag.Parse()

	privateKey, err := rsa.GenerateKey(rand.Reader, int(bits))
	if err != nil {
		return fmt.Errorf("failed to generate privatekey: %w", err)
	}

	privateFile, err := os.Create("private.pem")
	if err != nil {
		return fmt.Errorf("failed to create private key file: %w", err)
	}

	defer privateFile.Close()

	privateBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	if err := pem.Encode(privateFile, &privateBlock); err != nil {
		return fmt.Errorf("failed to encode private: %w", err)
	}

	publicBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to generate publickey: %w", err)
	}

	publicFile, err := os.Create("public.pem")
	if err != nil {
		return fmt.Errorf("failed to create public key file: %w", err)
	}

	defer publicFile.Close()

	publicBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicBytes,
	}

	if err := pem.Encode(publicFile, &publicBlock); err != nil {
		return fmt.Errorf("failed to encode public: %w", err)
	}

	return nil
}

func genToken() error {
	// ======== read private key ========

	keyName := "infra/keys/c2e055bb-f637-4cc3-9b4a-916a8b31304a.pem"

	file, err := os.Open(keyName)
	if err != nil {
		return fmt.Errorf("opening private key: %w", err)
	}
	defer file.Close()

	privatePEM, err := io.ReadAll(io.LimitReader(file, 1024*1024))
	if err != nil {
		return fmt.Errorf("reading private key: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return fmt.Errorf("parsing private key: %w", err)
	}

	// ======== generate token ========

	claims := struct {
		jwt.RegisteredClaims
		Roles []string
	}{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "sales-api",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			Subject:   "token",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(720 * time.Hour)),
		},
		Roles: []string{"admin"},
	}

	method := jwt.GetSigningMethod("RS256")
	token := jwt.NewWithClaims(method, claims)
	token.Header["kid"] = "c2e055bb-f637-4cc3-9b4a-916a8b31304a"

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return fmt.Errorf("generating token: %w", err)
	}

	fmt.Println("========= TOKEN BEGIN =========")
	fmt.Println(tokenString)
	fmt.Println("========= TOKEN END ===========")
	fmt.Println()

	// ======== token validation ========

	parser := jwt.NewParser(jwt.WithValidMethods([]string{"RS256"}))

	var parsedClaims struct {
		jwt.RegisteredClaims
		Roles []string
	}

	keyFunc := func(t *jwt.Token) (any, error) {
		kid, ok := t.Header["kid"]
		if !ok {
			return nil, errors.New("missing kid key in token header")
		}

		_, ok = kid.(string)
		if !ok {
			return nil, errors.New("kid key must be a string")
		}

		return &privateKey.PublicKey, nil
	}

	parsedToken, err := parser.ParseWithClaims(tokenString, &parsedClaims, keyFunc)
	if err != nil {
		return fmt.Errorf("parsing token: %w", err)
	}

	if !parsedToken.Valid {
		return errors.New("token is invalid")
	}

	fmt.Println("Token is valid")

	return nil
}

func migrate() error {
	cfg := config.New()

	dbcfg := database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	}

	db, err := database.Open(dbcfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := schema.Migrate(ctx, db); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	if err := schema.Seed(ctx, db); err != nil {
		return fmt.Errorf("seed database: %w", err)
	}

	fmt.Println("migrations and seeding complete")
	return nil
}
