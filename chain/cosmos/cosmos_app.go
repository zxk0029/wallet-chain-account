package cosmos

//
//import (
//	"context"
//	"cosmossdk.io/math"
//	"encoding/base64"
//	"encoding/hex"
//	"encoding/json"
//	"fmt"
//	"github.com/cosmos/cosmos-sdk/client"
//	"github.com/cosmos/cosmos-sdk/codec"
//	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
//	"math/big"
//	"strconv"
//
//	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
//	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
//	sdktypes "github.com/cosmos/cosmos-sdk/types"
//	"github.com/cosmos/cosmos-sdk/types/tx/signing"
//	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
//	"github.com/cosmos/cosmos-sdk/x/bank/types"
//)
//
//type SignV1TransactionParams struct {
//	ChainID       string
//	From          string
//	To            string
//	Memo          string
//	AmountIn      string
//	Fee           string
//	Gas           string
//	AccountNumber uint64
//	Sequence      uint64
//	Decimal       int
//	PrivateKey    string
//}
//
//func SignV1Transaction(params SignV1TransactionParams) (string, error) {
//	// Parse the amount and fee
//	amount, ok := new(big.Int).SetString(params.AmountIn, 10)
//	if !ok {
//		return "", fmt.Errorf("invalid amount value")
//	}
//	fee, ok := new(big.Int).SetString(params.Fee, 10)
//	if !ok {
//		return "", fmt.Errorf("invalid fee value")
//	}
//
//	// Convert amount and fee to the correct decimal
//	amount = amount.Mul(amount, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(params.Decimal)), nil))
//	fee = fee.Mul(fee, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(params.Decimal)), nil))
//
//	// Verify addresses
//	if err := verifyAddress(params.From); err != nil {
//		return "", fmt.Errorf("invalid from address: %w", err)
//	}
//	if err := verifyAddress(params.To); err != nil {
//		return "", fmt.Errorf("invalid to address: %w", err)
//	}
//
//	// Create the send message
//	sendMsg := types.NewMsgSend(
//		sdktypes.MustAccAddressFromBech32(params.From),
//		sdktypes.MustAccAddressFromBech32(params.To),
//		sdktypes.NewCoins(sdktypes.NewCoin("uatom", math.NewIntFromBigInt(amount))),
//	)
//
//	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
//	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)
//	txBuilder := txConfig.NewTxBuilder()
//
//	// Create the transaction builder
//	//clientCtx := client.Context{}.
//	//	WithChainID(params.ChainID).
//	//	WithTxConfig(txConfig)
//	//
//	//txBuilder := clientCtx.TxConfig.NewTxBuilder()
//
//	txBuilder.SetMemo(params.Memo)
//	gas, err := strconv.ParseUint(params.Gas, 10, 64)
//	if err != nil {
//		return "", fmt.Errorf("parse uint gas: %w", err)
//	}
//	txBuilder.SetGasLimit(gas)
//	txBuilder.SetFeeAmount(sdktypes.NewCoins(sdktypes.NewCoin("uatom", math.NewIntFromBigInt(fee))))
//
//	// Add the message to the transaction
//	if err := txBuilder.SetMsgs(sendMsg); err != nil {
//		return "", fmt.Errorf("failed to set messages: %w", err)
//	}
//
//	// Create the auth info
//	privKeyBytes, err := hex.DecodeString(params.PrivateKey)
//	if err != nil {
//		return "", fmt.Errorf("invalid private key: %w", err)
//	}
//	privKey := secp256k1.PrivKey{Key: privKeyBytes}
//	pubKey := privKey.PubKey()
//
//	//signerData := authsigning.SignerData{
//	//	ChainID:       params.ChainID,
//	//	AccountNumber: params.AccountNumber,
//	//	Sequence:      params.Sequence,
//	//}
//
//	// Sign the transaction
//	signMode := signing.SignMode_SIGN_MODE_DIRECT
//	sigData := signing.SingleSignatureData{
//		SignMode:  signMode,
//		Signature: nil,
//	}
//	sig := signing.SignatureV2{
//		PubKey:   pubKey,
//		Data:     &sigData,
//		Sequence: params.Sequence,
//	}
//	if err := txBuilder.SetSignatures(sig); err != nil {
//		return "", fmt.Errorf("failed to set signatures: %w", err)
//	}
//
//	// Create a context
//	ctx := context.Background()
//
//	// Create a tx factory
//	txFactory := clienttx.Factory{}.
//		WithChainID(params.ChainID).
//		WithKeybase(client.NewInMemoryKeyring()).
//		WithTxConfig(txConfig).
//		WithAccountNumber(params.AccountNumber).
//		WithSequence(params.Sequence)
//
//	// Sign the transaction
//	if err := clienttx.Sign(ctx, txFactory, "", txBuilder, true); err != nil {
//		return "", fmt.Errorf("failed to sign transaction: %w", err)
//	}
//
//	clienttx.BroadcastTx()
//	// Create the final transaction
//	txRaw := txBuilder.GetTx()
//
//	// Convert the transaction to base64
//	txBytes, err := txConfig.TxEncoder()(txRaw)
//	if err != nil {
//		return "", fmt.Errorf("failed to encode transaction: %w", err)
//	}
//	txBytesBase64 := base64.StdEncoding.EncodeToString(txBytes)
//
//	// Create the final JSON object
//	txRawJSON := map[string]string{
//		"tx_bytes": txBytesBase64,
//		"mode":     "BROADCAST_MODE_SYNC",
//	}
//	txJSON, err := json.Marshal(txRawJSON)
//	if err != nil {
//		return "", fmt.Errorf("failed to marshal JSON: %w", err)
//	}
//
//	return string(txJSON), nil
//}
//
//func verifyAddress(address string) error {
//	_, err := sdktypes.AccAddressFromBech32(address)
//	return err
//}
//
//func main() {
//	params := SignV1TransactionParams{
//		ChainID:       "cosmoshub-4",
//		From:          "cosmos1abc...",
//		To:            "cosmos1def...",
//		Memo:          "test transaction",
//		AmountIn:      "100",
//		Fee:           "1",
//		Gas:           "200000",
//		AccountNumber: 123,
//		Sequence:      456,
//		Decimal:       6,
//		PrivateKey:    "your_private_key_in_hex",
//	}
//
//	txJSON, err := SignV1Transaction(params)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//
//	fmt.Println("Signed Transaction:", txJSON)
//}
