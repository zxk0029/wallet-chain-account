package solana

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"
)

func Test_GetSupportChains(t *testing.T) {
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
		validPublicKey        = "7e376c64c64e88054b7a2d25dc716f45551d2f796ddc9e7be405e49c522b887c"
		validPublicKeyAddress = "9VhPRjzizPY95TyBrve7heeJTZnofgkQYJpLxRSZGZ3H"
		invalidPublicKey      = "invalid_hex"
	)

	adaptor := &ChainAdaptor{}

	t.Run("Valid Public Key", func(t *testing.T) {
		req := &account.ConvertAddressRequest{
			Chain:     ChainName,
			Network:   "mainnet",
			PublicKey: validPublicKey,
		}
		resp, err := adaptor.ConvertAddress(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.Equal(t, validPublicKeyAddress, resp.Address)
		assert.Equal(t, "convert address success", resp.Msg)

		t.Logf("Response Code: %v", resp.Code)
		t.Logf("Response Message: %s", resp.Msg)
		t.Logf("Converted Address: %s", resp.Address)
	})

	t.Run("Empty Public Key", func(t *testing.T) {
		req := &account.ConvertAddressRequest{
			Chain:     ChainName,
			Network:   "mainnet",
			PublicKey: "",
		}
		resp, err := adaptor.ConvertAddress(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Invalid Public Key Format", func(t *testing.T) {
		req := &account.ConvertAddressRequest{
			Chain:     ChainName,
			Network:   "mainnet",
			PublicKey: invalidPublicKey,
		}
		resp, err := adaptor.ConvertAddress(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestChainAdaptor_ValidAddress(t *testing.T) {
	const (
		validAddress   = "9VhPRjzizPY95TyBrve7heeJTZnofgkQYJpLxRSZGZ3H"
		invalidAddress = "invalid_address"
	)

	adaptor := &ChainAdaptor{}

	t.Run("Valid Address", func(t *testing.T) {
		req := &account.ValidAddressRequest{
			Chain:   ChainName,
			Network: "mainnet",
			Address: validAddress,
		}
		resp, err := adaptor.ValidAddress(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.Equal(t, true, resp.Valid)

		t.Logf("Response Code: %v", resp.Code)
		t.Logf("Response Message: %s", resp.Msg)
		t.Logf("Is Valid: %v", resp.Valid)
	})

	t.Run("Invalid Address", func(t *testing.T) {
		req := &account.ValidAddressRequest{
			Chain:   ChainName,
			Network: "mainnet",
			Address: invalidAddress,
		}
		resp, err := adaptor.ValidAddress(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_ERROR, resp.Code)
		assert.Equal(t, false, resp.Valid)

		t.Logf("Response Code: %v", resp.Code)
		t.Logf("Response Message: %s", resp.Msg)
		t.Logf("Is Valid: %v", resp.Valid)
	})

	t.Run("Empty Address", func(t *testing.T) {
		req := &account.ValidAddressRequest{
			Chain:   ChainName,
			Network: "mainnet",
			Address: "",
		}
		resp, err := adaptor.ValidAddress(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_ERROR, resp.Code)
		assert.Equal(t, false, resp.Valid)

		t.Logf("Response Code: %v", resp.Code)
		t.Logf("Response Message: %s", resp.Msg)
		t.Logf("Is Valid: %v", resp.Valid)
	})

}

func TestChainAdaptor_GetBlockByNumber(t *testing.T) {
	solClient, err := NewSolHttpClientAll(baseURL, withDebug)
	assert.NoError(t, err, "failed to initialize sol solclient")

	adaptor := &ChainAdaptor{
		solCli: solClient,
	}

	t.Run("Valid Block Number", func(t *testing.T) {
		req := &account.BlockNumberRequest{
			Chain:  ChainName,
			Height: 300944802,
			ViewTx: true,
		}

		resp, err := adaptor.GetBlockByNumber(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.NotEmpty(t, resp.Hash)
		assert.Equal(t, int64(300944802), resp.Height)

		assert.NotNil(t, resp.Transactions)
		if len(resp.Transactions) > 0 {
			tx := resp.Transactions[0]
			assert.NotEmpty(t, tx.Hash)

			t.Logf("Transaction Hash: %s", tx.Hash)
			t.Logf("From Address: %s", tx.From)
			t.Logf("To Address: %s", tx.To)
			t.Logf("Amount: %s", tx.Amount)
		}

		t.Logf("Block Height: %d", resp.Height)
		t.Logf("Block Hash: %s", resp.Hash)
		t.Logf("Transaction Count: %d", len(resp.Transactions))

		respJson, err := json.Marshal(resp)
		assert.NoError(t, err)
		t.Logf("Block json: %s", string(respJson))
	})

	t.Run("Zero Block Number", func(t *testing.T) {
		req := &account.BlockNumberRequest{
			Chain:  ChainName,
			Height: 0,
			ViewTx: true,
		}

		resp, err := adaptor.GetBlockByNumber(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.NotEmpty(t, resp.Hash)
		t.Logf("Genesis Block Hash: %s", resp.Hash)

		respJson, err := json.Marshal(resp)
		assert.NoError(t, err)
		t.Logf("Block json: %s", string(respJson))
	})

	t.Run("Invalid Block Number", func(t *testing.T) {
		req := &account.BlockNumberRequest{
			Chain:  ChainName,
			Height: 999999999999,
			ViewTx: true,
		}

		resp, err := adaptor.GetBlockByNumber(req)

		assert.Error(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_ERROR, resp.Code)
		assert.NotEmpty(t, resp.Msg)

		t.Logf("Error Message: %s", resp.Msg)
	})

}
