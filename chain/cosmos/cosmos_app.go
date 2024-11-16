package cosmos

//
//import (
//	"encoding/base64"
//	"encoding/hex"
//	"fmt"
//	"math/big"
//
//	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
//	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
//	"github.com/cosmos/cosmos-sdk/types"
//	"github.com/cosmos/cosmos-sdk/types/tx/signing"
//	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
//)
//
//func main1() {
//	// 示例参数
//	//params := map[string]interface{}{
//	//	"chainId":       "cosmoshub-4",
//	//	"from":          "cosmos1nvcgd368m4pm5mm3ppzawhsq6grra4ejnppplx",
//	//	"to":            "cosmos1f9cae0zhtjm4ul2uhq60vs8ulm9am4w4v5rw40",
//	//	"memo":          "test memo",
//	//	"amount_in":     "1000000",
//	//	"fee":           "1000",
//	//	"gas":           "200000",
//	//	"accountNumber": 12345,
//	//	"sequence":      1,
//	//	"decimal":       6,
//	//	"privateKey":    "your-private-key-in-hex",
//	//}
//	//
//	//txRaw, err := SignV2Transaction(params)
//	//if err != nil {
//	//	fmt.Println("Error signing transaction:", err)
//	//	return
//	//}
//	//
//	//fmt.Println("Signed transaction:", txRaw)
//	fmt.Println("Signed transaction:", "chenx")
//}
//
//func SignV2Transaction(params map[string]interface{}) (string, error) {
//	// 解析参数
//	chainId := params["chainId"].(string)
//	from := params["from"].(string)
//	to := params["to"].(string)
//	memo := params["memo"].(string)
//	amountIn := params["amount_in"].(string)
//	fee := params["fee"].(string)
//	gas := params["gas"].(string)
//	accountNumber := params["accountNumber"].(int)
//	sequence := params["sequence"].(int)
//	decimal := params["decimal"].(int)
//	privateKeyHex := params["privateKey"].(string)
//
//	// 转换金额和费用
//	amount := new(big.Int)
//	amount.SetString(amountIn, 10)
//	amount = amount.Mul(amount, big.NewInt(int64(10^decimal)))
//
//	feeAmount := new(big.Int)
//	feeAmount.SetString(fee, 10)
//	feeAmount = feeAmount.Mul(feeAmount, big.NewInt(int64(10^decimal)))
//
//	unit := "uatom"
//
//	// 验证输入
//	if amount.String()[0] == '.' || feeAmount.String()[0] == '.' {
//		return "", fmt.Errorf("input amount value invalid")
//	}
//
//	if !verifyAddress(from) || !verifyAddress(to) {
//		return "", fmt.Errorf("input address value invalid")
//	}
//
//	// 创建发送消息
//	sendMessage := createSendMessage(from, to, amount.String(), unit)
//	messages := []types.Msg{sendMessage}
//
//	// 创建交易体
//	txBody := createTxBody(messages, memo)
//
//	// 获取公钥
//	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
//	if err != nil {
//		return "", fmt.Errorf("error decoding private key: %w", err)
//	}
//
//	privKey := &secp256k1.PrivKey{Key: privateKeyBytes}
//	pubKey := privKey.PubKey()
//
//	// 创建 authInfo
//	authInfo := createAuthInfo(pubKey, sequence, feeAmount.String(), unit, gas)
//
//	// 获取签名文档
//	signDoc := getSignDoc(chainId, txBody, authInfo, accountNumber)
//
//	// 生成签名
//	signature, err := getDirectSignature(signDoc, privKey)
//	if err != nil {
//		return "", fmt.Errorf("error signing transaction: %w", err)
//	}
//
//	// 创建交易的原始字节
//	txRawBytes, err := createTxRawBytes(txBody, authInfo, signature)
//	if err != nil {
//		return "", fmt.Errorf("error creating transaction raw bytes: %w", err)
//	}
//
//	// 转换为 Base64 编码
//	txBytesBase64 := base64.StdEncoding.EncodeToString(txRawBytes)
//	txRaw := map[string]string{
//		"tx_bytes": txBytesBase64,
//		"mode":     "BROADCAST_MODE_SYNC",
//	}
//
//	// 返回 JSON 字符串
//	return fmt.Sprintf(`{"tx_bytes":"%s","mode":"%s"}`, txRaw["tx_bytes"], txRaw["mode"]), nil
//}
//
//func createSendMessage(from, to, amount, unit string) *banktypes.MsgSend {
//	fromAddr, _ := types.AccAddressFromBech32(from)
//	toAddr, _ := types.AccAddressFromBech32(to)
//	coin := types.NewCoin(unit, types.NewIntFromBigInt(new(big.Int).SetString(amount, 10)))
//	return banktypes.NewMsgSend(fromAddr, toAddr, types.NewCoins(coin))
//}
//
//func createTxBody(messages []types.Msg, memo string) *types.TxBody {
//	anyMessages := make([]*codectypes.Any, len(messages))
//	for i, msg := range messages {
//		anyMsg, _ := codectypes.NewAnyWithValue(msg)
//		anyMessages[i] = anyMsg
//	}
//	return &types.TxBody{
//		Messages: anyMessages,
//		Memo:     memo,
//	}
//}
//
//func createAuthInfo(pubKey types.PubKey, sequence int, feeAmount, unit, gas string) *types.AuthInfo {
//	signerInfo := &types.SignerInfo{
//		PublicKey: codectypes.UnsafePackAny(pubKey),
//		ModeInfo: &types.ModeInfo{
//			Sum: &types.ModeInfo_Single_{
//				Single: &types.ModeInfo_Single{Mode: signing.SignMode_SIGN_MODE_DIRECT},
//			},
//		},
//		Sequence: uint64(sequence),
//	}
//	fee := &types.Fee{
//		Amount:   types.NewCoins(types.NewCoin(unit, types.NewIntFromBigInt(new(big.Int).SetString(feeAmount, 10)))),
//		GasLimit: types.NewIntFromBigInt(new(big.Int).SetString(gas, 10)),
//	}
//	return &types.AuthInfo{
//		SignerInfos: []*types.SignerInfo{signerInfo},
//		Fee:         fee,
//	}
//}
//
//func getSignDoc(chainId string, txBody *types.TxBody, authInfo *types.AuthInfo, accountNumber int) *types.SignDoc {
//	return &types.SignDoc{
//		BodyBytes:     txBody.GetSignBytes(),
//		AuthInfoBytes: authInfo.GetSignBytes(),
//		ChainId:       chainId,
//		AccountNumber: uint64(accountNumber),
//	}
//}
//
//func getDirectSignature(signDoc *types.SignDoc, privKey *secp256k1.PrivKey) ([]byte, error) {
//	signature, err := privKey.Sign(signDoc.GetSignBytes())
//	if err != nil {
//		return nil, fmt.Errorf("error signing document: %w", err)
//	}
//	return signature, nil
//}
//
//func createTxRawBytes(txBody *types.TxBody, authInfo *types.AuthInfo, signature []byte) ([]byte, error) {
//	txRaw := &types.TxRaw{
//		BodyBytes:     txBody.GetSignBytes(),
//		AuthInfoBytes: authInfo.GetSignBytes(),
//		Signatures:    [][]byte{signature},
//	}
//	return txRaw.Marshal()
//}
//
//func verifyAddress(address string) bool {
//	_, err := types.AccAddressFromBech32(address)
//	return err == nil
//}
