package aptos

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"
	"github.com/stretchr/testify/assert"
	"testing"

	"golang.org/x/crypto/sha3"
)

func Test_GenerateAptosAccount(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Private Key: 0x%s\n", hex.EncodeToString(privateKey))
	fmt.Printf("Public Key: 0x%s\n", hex.EncodeToString(publicKey))

	hasher := sha3.New256()
	hasher.Write(publicKey)
	address := hasher.Sum(nil)
	fmt.Printf("Aptos Address: 0x%s\n", hex.EncodeToString(address))

	//Private Key: 0xc0b79816fd85de6f645087a54ac4aa2903dead69acf9426ab010c97dd58aed798862f29d3c1f067cbe9eaba619e5dfcb269a9f059cb2c93fab362d1e93c3281c
	//Public Key: 0x8862f29d3c1f067cbe9eaba619e5dfcb269a9f059cb2c93fab362d1e93c3281c
	//Aptos Address: 0xfc38d27af874e409de8056d11cc8e10b8f8449e6f723a59251f04e62a24d7475
}

func TestChainAdaptor_GetSupportChains(t *testing.T) {
	adaptor := ChainAdaptor{}

	req := &account.SupportChainsRequest{
		Chain:   ChainName,
		Network: "mainnet",
	}

	resp, err := adaptor.GetSupportChains(req)

	if err != nil {
		t.Errorf("GetSupportChains failed with error: %v", err)
	}
	fmt.Printf("resp: %s\n", resp)

	if resp.Code != common2.ReturnCode_SUCCESS {
		t.Errorf("Expected success code, got %v", resp.Code)
	}

	if !resp.Support {
		t.Error("Expected Support to be true")
	}
}

func TestChainAdaptor_ConvertAddress(t *testing.T) {
	const (
		hexPublicKey          = "8862f29d3c1f067cbe9eaba619e5dfcb269a9f059cb2c93fab362d1e93c3281c"
		hexPublicKeyToAddress = "0xfc38d27af874e409de8056d11cc8e10b8f8449e6f723a59251f04e62a24d7475"
	)

	tests := []struct {
		name        string
		publicKey   string
		wantAddress string
		wantErr     bool
	}{
		{
			name:        "Valid public key",
			publicKey:   hexPublicKey,
			wantAddress: hexPublicKeyToAddress,
			wantErr:     false,
		},
		{
			name:        "Empty public key",
			publicKey:   "",
			wantAddress: "0xa7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a",
			wantErr:     false,
		},
	}

	adaptor := ChainAdaptor{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &account.ConvertAddressRequest{
				Chain:     ChainName,
				Network:   "mainnet",
				PublicKey: tt.publicKey,
			}

			resp, err := adaptor.ConvertAddress(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("resp: %s\n", resp)

			if resp.Code != common2.ReturnCode_SUCCESS {
				t.Errorf("Expected success code, got %v", resp.Code)
			}

			if resp.Address != tt.wantAddress {
				t.Errorf("ConvertAddress() got = %v, want %v", resp.Address, tt.wantAddress)
			}
		})
	}
}

func TestChainAdaptor_GetTxByHash(t *testing.T) {
	const (
		baseURL     = "https://api.mainnet.aptoslabs.com/"
		apiKey      = "aptoslabs_7Gd8hUMMp85_JxF2SXZCDcmeP4tjuuBXjwFwqyY6nTFup"
		withDebug   = true
		validTxHash = "0x4e76f0d0d244685e0f2d3f05dc8637cc8330baf469903d8eb497b7412e262e47"
	)

	aptosClient, err := NewAptosClientAll(baseURL, apiKey, withDebug)
	assert.NoError(t, err, "failed to initialize aptos client")

	adaptor := ChainAdaptor{
		aptosClient: aptosClient,
	}

	tests := []struct {
		name     string
		hash     string
		wantCode common2.ReturnCode
		wantErr  bool
	}{
		{
			name:     "Valid Transaction Hash",
			hash:     validTxHash,
			wantCode: common2.ReturnCode_SUCCESS,
			wantErr:  false,
		},
		{
			name:     "Invalid Transaction Hash",
			hash:     "0xinvalid_hash",
			wantCode: common2.ReturnCode_ERROR,
			wantErr:  false,
		},
		{
			name:     "Empty Transaction Hash",
			hash:     "",
			wantCode: common2.ReturnCode_ERROR,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &account.TxHashRequest{
				Hash: tt.hash,
			}
			got, err := adaptor.GetTxByHash(req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, tt.wantCode, got.Code)

			if tt.wantCode == common2.ReturnCode_SUCCESS {
				assert.NotNil(t, got.Tx)
				assert.Equal(t, tt.hash, got.Tx.Hash)
				t.Logf("Transaction Hash: %s", got.Tx.Hash)
			}

			t.Logf("Response Code: %v", got.Code)
			t.Logf("Response Message: %s", got.Msg)
		})
	}

	//t.Run("Concurrent Requests", func(t *testing.T) {
	//	var wg sync.WaitGroup
	//	for i := 0; i < 5; i++ {
	//		wg.Add(1)
	//		go func() {
	//			defer wg.Done()
	//			req := &account.TxHashRequest{
	//				Hash: validTxHash,
	//			}
	//			resp, err := adaptor.GetTxByHash(req)
	//			assert.NoError(t, err)
	//			assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
	//		}()
	//	}
	//	wg.Wait()
	//})

}
