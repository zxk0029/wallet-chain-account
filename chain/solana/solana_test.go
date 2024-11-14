package solana

import (
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

func TestChainAdaptor_BuildSignedTransaction(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.BuildSignedTransaction(&account.SignedTransactionRequest{
		Chain:    ChainName,
		Network:  "mainnet",
		Base64Tx: createTestBase64Tx(),
	})
	if err != nil {
		log.Error("TestChainAdaptor_BuildSignedTransaction failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.SignedTx)
}
func TestChainAdaptor_VerifySignedTransaction(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.VerifySignedTransaction(&account.VerifyTransactionRequest{
		Chain:     ChainName,
		Network:   "mainnet",
		Signature: "4MSGKyqDVHeGMWmSPHm5yDEmUcMdycc3LQsEq1Mu8HHGUjrUNRLV4TJPT1sbTyhyRhpMSdRW1ANty84asVQnEsmb2RupewnYX2jNjDobQ2deRA5q6sMcCrBVTeKjZ25PKuGKcxYgXDSEe2SZ6DPvg9BTLZgEKWTxNKKKGP4VPgYefhQ7grm3X9DHBnkLEpfxLUDzeGeMbESCPVw62wk1SVN1rzEGzpTfauvq3SzQb8n1PjAVaeSLkHqyy734yMJVvwWonPBDWQMAVwWomf4cFMKfQboR8ZBsp9cU7",
	})
	if err != nil {
		log.Error("TestChainAdaptor_VerifySignedTransaction failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
}
func createTestBase64Tx() string {

	testTx := TxStructure{
		Nonce:           "7cNmDJkzZLyXqP9q6ccznkuy4UkxiJCEu9QnWAXcrwDe",
		FromAddress:     "7YcpSkLK7gnSJ4JpysHR9BQgwe2gfffRQmMxHDbNf5ve",
		ToAddress:       "EUVrmoaKaSsHNkMFw7mVARR522wwH41BFRMha3WC8gha",
		Value:           "0.66",
		FromPrivateKey:  "",
		ContractAddress: "So11111111111111111111111111111111111111112", //5VzPuctbhMdqZBpxgxHCyH41sSckqPEKZ7qxbdgMN29Fbvmnpy3x6GcmUFxFw98oy3LcEEVCxwdr4gyQwcboSW6C
		//ContractAddress: "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB",//3L64aQvAmdhbaZJFdWXTSjLgmH1GwBhNE8eezqCFAHRvj9a76bwXoarivTSjzAJLiJ48CxtZ5Zke3djnfhuckKs
	}

	jsonBytes, err := json.Marshal(testTx)
	if err != nil {
		panic(err)
	}

	base64Str := base64.StdEncoding.EncodeToString(jsonBytes)
	return base64Str
}
