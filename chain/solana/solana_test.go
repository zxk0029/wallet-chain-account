package solana

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dapplink-labs/wallet-chain-account/rpc/common"
	"testing"

	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"

	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
)

func setup() (chain.IChainAdaptor, error) {
	conf, err := config.New("../../config.yml")
	if err != nil {
		log.Error("load config failed, error:", err)
		return nil, err
	}
	adaptor, err := NewChainAdaptor(conf)
	if err != nil {
		log.Error("create chain adaptor failed, error:", err)
		return nil, err
	}
	return adaptor, nil
}

// tx, err := solana.NewTransaction(
// []solana.Instruction{
// system.NewTransferInstruction(
// value,
// fromPubkey,
// toPubkey,
// ).Build(),
// },
// solana.HashFromBytes(binary.BigEndian.AppendUint64(make([]byte, 24), data.Nonce)),
// solana.TransactionPayer(fromPubkey),
func TestChainAdaptor_CreateUnSignTransaction(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.CreateUnSignTransaction(&account.UnSignTransactionRequest{
		Chain:    ChainName,
		Network:  "mainnet",
		Base64Tx: createTestBase64Tx(),
	})
	if err != nil {
		log.Error("CreateUnSignTransaction failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.UnSignTx)
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
