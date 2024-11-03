package aptos

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetNodeInfo(t *testing.T) {
	// Test configuration
	const (
		baseURL   = "https://api.mainnet.aptoslabs.com/"
		apiKey    = "aptoslabs_7Gd8hUMMp85_JxF2SXZCDcmeP4tjuuBXjwFwqyY6nTFup"
		withDebug = false
	)

	// Initialize client
	client, err := NewAptosClientAll(baseURL, apiKey, withDebug)
	assert.NoError(t, err, "client initialization failed")
	assert.NotNil(t, client, "client should not be nil")

	// Test getting node info
	t.Run("Get Node Info", func(t *testing.T) {
		resp, err := client.GetNodeInfo()
		marshal, _ := json.Marshal(resp)
		fmt.Println("GetNodeInfo resp", string(marshal))
		assert.NoError(t, err, "failed to get node info")
		assert.NotNil(t, resp, "response should not be nil")

		// Validate node info data
		t.Run("Node Info Validation", func(t *testing.T) {
			assert.Greater(t, resp.ChainID, uint8(0), "chain id should be greater than 0")
			assert.NotEmpty(t, resp.Epoch, "epoch should not be empty")
			assert.NotEmpty(t, resp.LedgerVersion, "ledger version should not be empty")
			// oldest_ledger_version is "0"
			assert.NotNil(t, resp.OldestLedgerVersion, "oldest ledger version should not be nil")
			assert.NotEmpty(t, resp.LedgerTimestamp, "ledger timestamp should not be empty")
			assert.NotEmpty(t, resp.NodeRole, "node role should not be empty")
			// OldestBlockHeight is "0"
			assert.NotNil(t, resp.OldestBlockHeight, "oldest block height should not be nil")
			assert.NotEmpty(t, resp.BlockHeight, "block height should not be empty")
			assert.NotEmpty(t, resp.GitHash, "git hash should not be empty")
		})

		// Log node info details
		t.Run("Log Details", func(t *testing.T) {
			t.Log("Node Information:")
			t.Logf("Chain ID: %d", resp.ChainID)
			t.Logf("Epoch: %d", resp.Epoch)
			t.Logf("Ledger Version: %d", resp.LedgerVersion)
			t.Logf("Oldest Ledger Version: %d", resp.OldestLedgerVersion)
			t.Logf("Ledger Timestamp: %d", resp.LedgerTimestamp)
			t.Logf("Node Role: %s", resp.NodeRole)
			t.Logf("Oldest Block Height: %d", resp.OldestBlockHeight)
			t.Logf("Block Height: %d", resp.BlockHeight)
			t.Logf("Git Hash: %s", resp.GitHash)
		})
	})
}

func TestClient_GetAccount(t *testing.T) {
	const (
		baseURL   = "https://api.mainnet.aptoslabs.com/"
		apiKey    = "aptoslabs_7Gd8hUMMp85_JxF2SXZCDcmeP4tjuuBXjwFwqyY6nTFup"
		testAddr  = "0x070ebd0a6fffebd2913cbaa6c350db20b739cc93c1d83f6856d1d34900e3162e"
		withDebug = false
	)

	client, err := NewAptosClientAll(baseURL, apiKey, withDebug)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	t.Run("valid address", func(t *testing.T) {
		account, err := client.GetAccount(testAddr)
		assert.NoError(t, err)
		assert.NotNil(t, account)
		assert.NotEmpty(t, account.AuthenticationKey)
		assert.Greater(t, account.SequenceNumber, uint64(0))
		t.Logf("success: sequence=%d, auth_key=%s",
			account.SequenceNumber, account.AuthenticationKey)
	})

	t.Run("empty address", func(t *testing.T) {
		account, err := client.GetAccount("")
		assert.Error(t, err)
		assert.Nil(t, account)
		t.Logf("fail: %v", err)
	})

	t.Run("invalid address format", func(t *testing.T) {
		account, err := client.GetAccount("invalid-address")
		assert.Error(t, err)
		assert.Nil(t, account)
		t.Logf("fail: %v", err)
	})
}

func TestClient_GetGasPrice(t *testing.T) {
	client, err := NewAptosClientAll(
		"https://api.mainnet.aptoslabs.com/",
		"aptoslabs_7Gd8hUMMp85_JxF2SXZCDcmeP4tjuuBXjwFwqyY6nTFup",
		false,
	)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	gasPrice, err := client.GetGasPrice()

	assert.NoError(t, err)
	assert.NotNil(t, gasPrice)

	assert.GreaterOrEqual(t, gasPrice.DeprioritizedGasEstimate, 0)
	assert.GreaterOrEqual(t, gasPrice.GasEstimate, gasPrice.DeprioritizedGasEstimate)
	assert.GreaterOrEqual(t, gasPrice.PrioritizedGasEstimate, gasPrice.GasEstimate)

	t.Logf("Gas Price Details:")
	t.Logf("Deprioritized: %d", gasPrice.DeprioritizedGasEstimate)
	t.Logf("Standard: %d", gasPrice.GasEstimate)
	t.Logf("Prioritized: %d", gasPrice.PrioritizedGasEstimate)
}

func TestClient_GetTransactionByHash(t *testing.T) {
	// Test configuration
	const (
		baseURL     = "https://api.mainnet.aptoslabs.com/"
		apiKey      = "aptoslabs_7Gd8hUMMp85_JxF2SXZCDcmeP4tjuuBXjwFwqyY6nTFup"
		withDebug   = true
		validTxHash = "0x43531969ff8e93de962ea65e5609c2b05de3aa5e78933d8925613e75d3d53772"
	)

	// Initialize client
	client, err := NewAptosClientAll(baseURL, apiKey, withDebug)
	assert.NoError(t, err, "client initialization failed")
	assert.NotNil(t, client, "client should not be nil")

	t.Run("Valid Transaction", func(t *testing.T) {
		resp, err := client.GetTransactionByHash(validTxHash)
		assert.NoError(t, err, "failed to get transaction")
		assert.NotNil(t, resp, "response should not be nil")

		t.Run("Transaction Metadata", func(t *testing.T) {
			assert.Equal(t, "1878359809", resp.Version)
			assert.Equal(t, validTxHash, resp.Hash)
			assert.Equal(t, "13", resp.GasUsed)
			assert.True(t, resp.Success)
			assert.Equal(t, "Executed successfully", resp.VMStatus)
		})

		// 验证变更
		t.Run("Transaction Changes", func(t *testing.T) {
			assert.NotEmpty(t, resp.Changes, "changes should not be empty")

			change := resp.Changes[0]
			assert.Equal(t, "write_resource", change.Type)
			assert.NotEmpty(t, change.Address)
			assert.NotEmpty(t, change.StateKeyHash)

			//data := change.Data.Data
			//assert.NotEmpty(t, data, "data should not be empty")
			//
			//rewards := data.Rewards
			//assert.NotEmpty(t, rewards, "rewards should not be empty")

			//lastReward := rewards[len(rewards)-1]
			//assert.Equal(t, "1730609200", lastReward.Eid)
			//assert.Equal(t, "12", lastReward.Share)
			//assert.Equal(t, "1730609324", lastReward.Unlock)
		})

		t.Run("Log Details", func(t *testing.T) {
			t.Log("Transaction Details:")
			t.Logf("Version: %s", resp.Version)
			t.Logf("Hash: %s", resp.Hash)
			t.Logf("Gas Used: %s", resp.GasUsed)
			t.Logf("Success: %v", resp.Success)
			t.Logf("VM Status: %s", resp.VMStatus)

			if len(resp.Changes) > 0 {
				t.Log("\nTransaction Changes:")
				marshal, _ := json.Marshal(resp.Changes)
				t.Logf("resp.Changes: %s", string(marshal))
			}
		})
	})

	t.Run("Error Cases", func(t *testing.T) {
		t.Run("Empty Hash", func(t *testing.T) {
			resp, err := client.GetTransactionByHash("")
			assert.Error(t, err)
			assert.Nil(t, resp)
			assert.Contains(t, err.Error(), "transaction hash cannot be empty")
		})

		t.Run("Invalid Hash", func(t *testing.T) {
			invalidHash := "0x1234567890abcdef"
			resp, err := client.GetTransactionByHash(invalidHash)
			assert.Error(t, err)
			assert.Nil(t, resp)
		})
	})
}

func TestClient_GetBlockByHeight(t *testing.T) {
	// Test configuration
	const (
		baseURL   = "https://api.mainnet.aptoslabs.com/"
		apiKey    = "aptoslabs_7Gd8hUMMp85_JxF2SXZCDcmeP4tjuuBXjwFwqyY6nTFup"
		withDebug = false
		// Use a known block height for testing
		validHeight uint64 = 247279394
	)

	// Initialize client
	client, err := NewAptosClientAll(baseURL, apiKey, withDebug)
	assert.NoError(t, err, "client initialization failed")
	assert.NotNil(t, client, "client should not be nil")

	// Test valid block height
	t.Run("Valid Block Height", func(t *testing.T) {
		resp, err := client.GetBlockByHeight(validHeight)
		assert.NoError(t, err, "failed to get block")
		assert.NotNil(t, resp, "response should not be nil")

		// Validate block data
		t.Run("Block Metadata", func(t *testing.T) {
			assert.NotEmpty(t, resp.BlockHeight, "block height should not be empty")
			assert.NotEmpty(t, resp.BlockHash, "block hash should not be empty")
			assert.NotEmpty(t, resp.BlockTimestamp, "block timestamp should not be empty")
			assert.NotEmpty(t, resp.FirstVersion, "first version should not be empty")
			assert.NotEmpty(t, resp.LastVersion, "last version should not be empty")
		})

		// Log block details
		t.Run("Log Details", func(t *testing.T) {
			t.Log("Block Details:")
			t.Logf("Block Height: %d", resp.BlockHeight)
			t.Logf("Block Hash: %s", resp.BlockHash)
			t.Logf("Timestamp: %d", resp.BlockTimestamp)
			t.Logf("First Version: %d", resp.FirstVersion)
			t.Logf("Last Version: %d", resp.LastVersion)
		})
	})

	// Test error cases
	t.Run("Error Cases", func(t *testing.T) {
		// Test invalid block height
		t.Run("Invalid Height", func(t *testing.T) {
			// Use a very large block height that shouldn't exist yet
			invalidHeight := uint64(999999999999)
			resp, err := client.GetBlockByHeight(invalidHeight)
			assert.Error(t, err, "should return an error")
			assert.Nil(t, resp, "response should be nil")
			t.Logf("Expected error: %v", err)
		})
	})
}
