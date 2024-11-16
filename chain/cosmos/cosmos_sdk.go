package cosmos

import (
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
)

func GenerateKeyPair() (string, *secp256k1.PrivKey, cryptotypes.PubKey, types.AccAddress, error) {
	config := types.GetConfig()
	config.SetBech32PrefixForAccount("cosmos", "cosmospub")

	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", nil, nil, nil, fmt.Errorf("生成熵失败: %v", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", nil, nil, nil, fmt.Errorf("生成助记词失败: %v", err)
	}

	return GenerateWalletFromMnemonic(mnemonic)
}

func GenerateWalletFromMnemonic(mnemonic string) (string, *secp256k1.PrivKey, cryptotypes.PubKey, types.AccAddress, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return "", nil, nil, nil, fmt.Errorf("无效的助记词: %s", mnemonic)
	}

	seed := bip39.NewSeed(mnemonic, "")
	master, ch := hd.ComputeMastersFromSeed(seed)

	path := hd.NewFundraiserParams(0, 118, 0).String()
	privKey, err := hd.DerivePrivateKeyForPath(master, ch, path)
	if err != nil {
		return "", nil, nil, nil, fmt.Errorf("派生私钥失败: %v", err)
	}

	privateKey := &secp256k1.PrivKey{Key: privKey}
	publicKey := privateKey.PubKey()
	addr := types.AccAddress(publicKey.Address())

	return mnemonic, privateKey, publicKey, addr, nil
}

func FromPrivKeyHex(privKeyHex string) (*secp256k1.PrivKey, error) {
	privKeyBytes, err := hex.DecodeString(privKeyHex)
	if err != nil {
		return nil, fmt.Errorf("解析私钥失败: %v", err)
	}
	privateKey := &secp256k1.PrivKey{Key: privKeyBytes}

	return privateKey, nil
}

func FromPubKeyHex(pubKeyHex string) (cryptotypes.PubKey, error) {
	pubKeyBytes, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return nil, fmt.Errorf("解析公钥失败: %v", err)
	}
	publicKey := &secp256k1.PubKey{Key: pubKeyBytes}

	return publicKey, nil
}

func FromAddressHex(addressStr string) (types.AccAddress, error) {
	address, err := types.AccAddressFromBech32(addressStr)
	if err != nil {
		return nil, fmt.Errorf("解析地址失败: %v", err)
	}

	return address, nil
}
