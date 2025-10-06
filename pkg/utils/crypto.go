package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
)

// GenerateECDSAKeyPair generates a new ECDSA key pair
func GenerateECDSAKeyPair() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

// PublicKeyToString encodes an ECDSA public key to a base64 string
func PublicKeyToString(pub *ecdsa.PublicKey) string {
	return base64.StdEncoding.EncodeToString(elliptic.Marshal(elliptic.P256(), pub.X, pub.Y))
}

// VerifySignature verifies an ECDSA signature for a given message
func VerifySignature(pub *ecdsa.PublicKey, message, signature string) bool {
	hash := sha256.Sum256([]byte(message))
	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil || len(sigBytes) != 64 {
		return false
	}
	r := big.NewInt(0).SetBytes(sigBytes[:32])
	s := big.NewInt(0).SetBytes(sigBytes[32:])
	return ecdsa.Verify(pub, hash[:], r, s)
}

// SignMessage signs a message using the given ECDSA private key
func SignMessage(priv *ecdsa.PrivateKey, message string) (string, error) {
	hash := sha256.Sum256([]byte(message))
	r, s, err := ecdsa.Sign(rand.Reader, priv, hash[:])
	if err != nil {
		return "", err
	}
	sig := append(r.Bytes(), s.Bytes()...)
	return base64.StdEncoding.EncodeToString(sig), nil
}

// DecodePublicKey decodes a base64-encoded ECDSA public key string
func DecodePublicKey(pubStr string) (*ecdsa.PublicKey, error) {
	bytes, err := base64.StdEncoding.DecodeString(pubStr)
	if err != nil {
		return nil, err
	}
	x, y := elliptic.Unmarshal(elliptic.P256(), bytes)
	if x == nil || y == nil {
		return nil, err
	}
	return &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}, nil
}

// TransferMessage creates a canonical message string for signing/verification
func TransferMessage(senderWalletID, receiverWalletID, amount int64) string {
	return fmt.Sprintf("%d:%d:%d", senderWalletID, receiverWalletID, amount)
}
