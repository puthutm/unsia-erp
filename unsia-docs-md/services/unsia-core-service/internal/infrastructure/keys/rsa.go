package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
)

// EnsureRSAKeys checks if the private key exists. If not, it generates a new 2048-bit RSA keypair and saves it.
// Returns the loaded/generated private key.
func EnsureRSAKeys(privPath, pubPath string) (*rsa.PrivateKey, error) {
	if privPath == "" {
		privPath = ".keys/private.pem"
	}
	if pubPath == "" {
		pubPath = ".keys/public.pem"
	}

	// Check if private key exists
	if _, err := os.Stat(privPath); err == nil {
		// Load existing key
		privBytes, err := os.ReadFile(privPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key file: %w", err)
		}
		block, _ := pem.Decode(privBytes)
		if block == nil || (block.Type != "RSA PRIVATE KEY" && block.Type != "PRIVATE KEY") {
			return nil, fmt.Errorf("invalid PEM private key block")
		}
		
		// Handle both PKCS1 and PKCS8 private keys
		var privKey *rsa.PrivateKey
		if block.Type == "RSA PRIVATE KEY" {
			privKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse PKCS1 private key: %w", err)
			}
		} else {
			key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse PKCS8 private key: %w", err)
			}
			var ok bool
			privKey, ok = key.(*rsa.PrivateKey)
			if !ok {
				return nil, fmt.Errorf("not an RSA private key")
			}
		}
		return privKey, nil
	}

	// If it doesn't exist, create directories
	if err := os.MkdirAll(filepath.Dir(privPath), 0700); err != nil {
		return nil, fmt.Errorf("failed to create directory for keys: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(pubPath), 0700); err != nil {
		return nil, fmt.Errorf("failed to create directory for keys: %w", err)
	}

	// Generate key pair
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA keypair: %w", err)
	}

	// Save Private Key
	privFile, err := os.OpenFile(privPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open private key file for writing: %w", err)
	}
	defer privFile.Close()

	privBytes := x509.MarshalPKCS1PrivateKey(privKey)
	privBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	}
	if err := pem.Encode(privFile, privBlock); err != nil {
		return nil, fmt.Errorf("failed to encode private key to PEM: %w", err)
	}

	// Save Public Key
	pubFile, err := os.OpenFile(pubPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open public key file for writing: %w", err)
	}
	defer pubFile.Close()

	pubBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}
	pubBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}
	if err := pem.Encode(pubFile, pubBlock); err != nil {
		return nil, fmt.Errorf("failed to encode public key to PEM: %w", err)
	}

	return privKey, nil
}
