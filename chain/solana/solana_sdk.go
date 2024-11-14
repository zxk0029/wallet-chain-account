package solana

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/cosmos/btcutil/base58"
)

func generateKeyPair() (*ed25519.PrivateKey, *ed25519.PublicKey, string, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, "", err
	}
	address := base58.Encode(publicKey)
	return &privateKey, &publicKey, address, nil
}

func PrivateKeyHexToPrivateKey(privateKeyHex string) (*ed25519.PrivateKey, error) {
	privKeyByteList, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key hex: %w", err)
	}
	privKey := ed25519.PrivateKey(privKeyByteList)
	return &privKey, nil
}

func PrivateKeyToPubKey(privateKey *ed25519.PrivateKey) (*ed25519.PublicKey, error) {
	if privateKey == nil {
		return nil, fmt.Errorf("private key is nil")
	}

	pubKey := (*privateKey).Public().(ed25519.PublicKey)
	return &pubKey, nil
}

func PrivateKeyHexToPubKey(privateKeyHex string) (*ed25519.PublicKey, error) {
	privKey, err := PrivateKeyHexToPrivateKey(privateKeyHex)
	if err != nil {
		return nil, err
	}
	return PrivateKeyToPubKey(privKey)
}

func PubKeyHexToPubKey(publicKeyHex string) (*ed25519.PublicKey, error) {
	pubKeyByteList, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key hex: %w", err)
	}
	pubKey := ed25519.PublicKey(pubKeyByteList)
	return &pubKey, nil
}

func PubKeyToPubKeyHex(publicKey *ed25519.PublicKey) (string, error) {
	if publicKey == nil {
		return "", fmt.Errorf("public key is nil")
	}
	return hex.EncodeToString(*publicKey), nil
}

func PubKeyToAddress(publicKey *ed25519.PublicKey) (string, error) {
	if publicKey == nil {
		return "", fmt.Errorf("public key is nil")
	}
	return base58.Encode(*publicKey), nil
}

func PubKeyHexToAddress(publicKeyHex string) (string, error) {
	pubKey, err := PubKeyHexToPubKey(publicKeyHex)
	if err != nil {
		return "", err
	}
	return PubKeyToAddress(pubKey)
}
