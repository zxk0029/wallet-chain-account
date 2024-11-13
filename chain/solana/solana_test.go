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
	// 创建测试数据结构
	testTx := TxStructure{
		Nonce:           "hsuXVct3kUjH1uxj1Wi8z93TkLTJo7YrHV5hzV61gYe",
		FromAddress:     "4wHd9tf4x4FkQ3JtgsMKyiEofEHSaZH5rYzfFKLvtESD",
		ToAddress:       "AaWEWZJZq2M4AUytd9XQGUTUXSpD85qERzbVEfXRjF7B",
		Value:           "0.01",
		FromPrivateKey:  "55a70321542da0b6123f37180e61993d5769f0a5d727f9c817151c1270c290963a7b3874ba467be6b81ea361e3d7453af8b81c88aedd24b5031fdda0bc71ad32",
		ContractAddress: "So11111111111111111111111111111111111111112",
	}

	// 转换为 JSON
	jsonBytes, err := json.Marshal(testTx)
	if err != nil {
		panic(err)
	}

	// 转换为 base64
	base64Str := base64.StdEncoding.EncodeToString(jsonBytes)
	return base64Str
}
