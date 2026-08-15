package selfupdate

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	_ "embed"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"
)

//go:embed pith-release.pub
var releasePublicKey []byte

func verifyManifestSignature(manifest, signature []byte) error {
	return verifyManifestSignatureWithKey(releasePublicKey, manifest, signature)
}

func verifyManifestSignatureWithKey(publicKeyPEM, manifest, signature []byte) error {
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		return fmt.Errorf("decode release public key")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("parse release public key: %w", err)
	}
	ecdsaKey, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("release public key is not ECDSA")
	}
	sig, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(signature)))
	if err != nil {
		return fmt.Errorf("decode release manifest signature: %w", err)
	}
	digest := sha256.Sum256(manifest)
	if !ecdsa.VerifyASN1(ecdsaKey, digest[:], sig) {
		return fmt.Errorf("release manifest signature verification failed")
	}
	return nil
}
