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
		log.Error("CreateUnSignTransaction failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.SignedTx)
}
func createTestBase64Tx() string {

	testTx := TxStructure{
		Nonce:          "5k4cL62LqAaBU1hYh6nEhnQ5EonPPSxenLDxTn2VtMik",
		FromAddress:    "7YcpSkLK7gnSJ4JpysHR9BQgwe2gfffRQmMxHDbNf5ve",
		ToAddress:      "EUVrmoaKaSsHNkMFw7mVARR522wwH41BFRMha3WC8gha",
		Value:          "0.01",
		FromPrivateKey: "5XH7UVWU3q7qUwCga2sH55Gtfg1osdBmiRjGxCU3iECBEXqtvxSQcUjNGdpjWrLgS3dago1WGb15KuruTXMjJoR8",
		//ContractAddress: "So11111111111111111111111111111111111111112",//5VzPuctbhMdqZBpxgxHCyH41sSckqPEKZ7qxbdgMN29Fbvmnpy3x6GcmUFxFw98oy3LcEEVCxwdr4gyQwcboSW6C
		ContractAddress: "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB",
	}

	jsonBytes, err := json.Marshal(testTx)
	if err != nil {
		panic(err)
	}

	base64Str := base64.StdEncoding.EncodeToString(jsonBytes)
	return base64Str
}
