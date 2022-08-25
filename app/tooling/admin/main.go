package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
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
