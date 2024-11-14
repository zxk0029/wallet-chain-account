package cosmos

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/log"
)

// gasLimit=200000, feeAmount=1000, from=cosmos19thxsunl9lzywglsndth5a278wtavawzzpv44q
// to=cosmos1l6vul20q74gw56fped8srkjq2x8d9m305gnxr2
func BuildUnSignTransaction(txStruct *TxStructure) ([]byte, error) {
	//ctx := client.Context{}.WithChainID("cosmoshub-4")
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)
	txBuilder := txConfig.NewTxBuilder()

	txBuilder.SetGasLimit(txStruct.GasLimit)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(txStruct.FeeAmount))))
	txBuilder.SetMemo(txStruct.Memo)

	txBuilder.SetSignatures(signingtypes.SignatureV2{Sequence: txStruct.Sequence})

	fromAddr, err := sdk.AccAddressFromBech32(txStruct.FromAddress)
	if err != nil {
		log.Error("from address AccAddressFromBech32 fail", "err", err)
		return nil, err
	}
	toAddr, err := sdk.AccAddressFromBech32(txStruct.ToAddress)
	if err != nil {
		log.Error("to address AccAddressFromBech32 fail", "err", err)
		return nil, err
	}
	amount := sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(txStruct.Amount)))
	msg := banktypes.NewMsgSend(fromAddr, toAddr, amount)

	err = txBuilder.SetMsgs(msg)
	if err != nil {
		log.Error("set tx message fail", "err", err)
		return nil, err
	}

	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		log.Error("tx encoder fail", "err", err)
		return nil, err
	}

	log.Info("Unsigned Tx: ", "data", txBytes)

	return txBytes, nil
}

//
//func CreateUnSignTransaction() {
//	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
//	txConfig := tx.NewTxConfig(cdc, tx.DefaultSignModes)
//	txBuilder := txConfig.NewTxBuilder()
//
//	txBuilder.SetGasLimit(200000)
//	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(1000))))
//	txBuilder.SetMemo("Test Transaction")
//
//	fromAddr, err := sdk.AccAddressFromBech32("cosmos19thxsunl9lzywglsndth5a278wtavawzzpv44q")
//	if err != nil {
//		log.Fatal(err)
//	}
//	toAddr, err := sdk.AccAddressFromBech32("cosmos1l6vul20q74gw56fped8srkjq2x8d9m305gnxr2")
//	if err != nil {
//		log.Fatal(err)
//	}
//	amount := sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(1000000)))
//	msg := types.NewMsgSend(fromAddr, toAddr, amount)
//
//	// 将消息添加到交易中
//	err = txBuilder.SetMsgs(msg)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	err = txBuilder.SetSignatures(sig)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 序列化交易
//	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 打印未签名的交易
//	fmt.Printf("Unsigned Tx: %X\n", txBytes)
//}
//
//func CreateSignTransaction() {
//	// 创建一个 client.Context
//	ctx := client.Context{}.WithChainID("cosmoshub-4")
//
//	// 创建一个 codec.Marshaler
//	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
//
//	// 创建一个 TxConfig
//	txConfig := tx.NewTxConfig(cdc, tx.DefaultSignModes)
//
//	// 创建一个 TxBuilder
//	txBuilder := txConfig.NewTxBuilder()
//
//	// 设置交易参数
//	txBuilder.SetGasLimit(200000)
//	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(1000))))
//	txBuilder.SetMemo("Test Transaction")
//
//	// 创建一个转账消息
//	fromAddr, err := sdk.AccAddressFromBech32("cosmos1abc...")
//	if err != nil {
//		log.Fatal(err)
//	}
//	toAddr, err := sdk.AccAddressFromBech32("cosmos1def...")
//	if err != nil {
//		log.Fatal(err)
//	}
//	amount := sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(1000000)))
//	msg := types.NewMsgSend(fromAddr, toAddr, amount)
//
//	// 将消息添加到交易中
//	err = txBuilder.SetMsgs(msg)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 签名交易
//	privKey := "23905eb45d6d199d97791a3a1e99f9b69959ddccfd420038727aa5c57ea626d8" // 获取私钥
//	signerData := signing.SignerData{
//		ChainID:       ctx.ChainID,
//		AccountNumber: 2424228, // 账户编号
//		Sequence:      9,       // 序列号
//	}
//
//	txsign.Sign
//	sig, err := clienttx.Sign(ctx, txConfig.SignModeHandler().DefaultMode(), signerData, txBuilder, privKey, txConfig, 0)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 将签名添加到交易中
//	err = txBuilder.SetSignatures(sig)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 获取交易对象
//	tx := txBuilder.GetTx()
//
//	// 序列化交易
//	txBytes, err := txConfig.TxEncoder()(tx)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 打印已签名的交易
//	fmt.Printf("Signed Tx: %X\n", txBytes)
//}
