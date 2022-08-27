package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var bits = *flag.Uint("bits", 2048, "bit size for private key generation")

func main() {
	flag.Parse()

	err := genKey()
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

	return nil
}
