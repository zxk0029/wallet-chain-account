package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/log"
)

func main() {

	encodedHash := "NSkPkTFwZDB7WxqaRO77HPP2b2jqrUU5vNaluhOGbpA="
	decodedHash, err := base64.StdEncoding.DecodeString(encodedHash)
	if err != nil {
		log.Error("解码失败: %v", err)
	}
	fmt.Printf("解码后的二进制数据: %x\n", string(decodedHash))
	fmt.Printf("数据长度: %d 字节\n", len(decodedHash))

	// 示例数据
	prefix := "cosmos"
	data := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}

	// 编码 Bech32 字符串
	bech32Str, err := bech32.ConvertAndEncode(prefix, data)
	if err != nil {
		log.Error("failed to encode bech32 string: %v", err)
	}
	fmt.Printf("Encoded Bech32 String: %s\n", bech32Str)

	// 解码 Bech32 字符串
	hrp, decodedData, err := bech32.DecodeAndConvert(bech32Str)
	if err != nil {
		log.Error("failed to decode bech32 string: %v", err)
	}
	fmt.Printf("Decoded Bech32 Prefix: %s\n", hrp)
	fmt.Printf("Decoded Bech32 Data: %X\n", decodedData)

	hrp1, decodedData1, err1 := bech32.DecodeAndConvert("NSkPkTFwZDB7WxqaRO77HPP2b2jqrUU5vNaluhOGbpA=")
	if err1 != nil {
		log.Error("failed to decode bech32 string: %v", err)
	}
	fmt.Printf("Decoded Bech32 Prefix: %s\n", hrp1)
	fmt.Printf("Decoded Bech32 Data: %X\n", decodedData1)
	//// 替换为你的 Cosmos 节点的 REST API 地址
	//restURL := "https://cosmos-rest.publicnode.com"
	//
	//// 获取最新的区块高度
	//latestBlockHeight := getLatestBlockHeight(restURL)
	//
	//// 获取指定高度的区块
	//block := getBlock(restURL, latestBlockHeight)
	//
	//// 打印区块信息
	//fmt.Printf("Block Height: %s\n", block.Block.Header.Height)
	//fmt.Printf("Block Time: %s\n", block.Block.Header.Time)
	//
	//// 解析区块中的交易
	//for _, txBytes := range block.Block.Data.Txs {
	//	decodeTx(restURL, txBytes)
	//}

	//var f = flag.String("c", "config.yml", "config path")
	//flag.Parse()
	//conf, err := config.New(*f)
	//if err != nil {
	//	panic(err)
	//}
	//dispatcher, err := chaindispatcher.New(conf)
	//if err != nil {
	//	log.Error("Setup dispatcher failed", "err", err)
	//	panic(err)
	//}
	//
	//grpcServer := grpc.NewServer(grpc.UnaryInterceptor(dispatcher.Interceptor))
	//defer grpcServer.GracefulStop()
	//
	//wallet2.RegisterWalletAccountServiceServer(grpcServer, dispatcher)
	//
	//listen, err := net.Listen("tcp", ":"+conf.Server.Port)
	//if err != nil {
	//	log.Error("net listen failed", "err", err)
	//	panic(err)
	//}
	//reflection.Register(grpcServer)
	//
	//log.Info("dapplink wallet rpc services start success", "port", conf.Server.Port)
	//
	//if err := grpcServer.Serve(listen); err != nil {
	//	log.Error("grpc server serve failed", "err", err)
	//	panic(err)
	//}
}

func getAccountBalance(queryClient types.QueryClient, address string) {
	// 创建查询请求
	req := &types.QueryAllBalancesRequest{
		Address: address,
		Pagination: &query.PageRequest{
			Limit: 100, // 设置分页限制
		},
	}

	// 发送查询请求
	res, err := queryClient.AllBalances(context.Background(), req)
	if err != nil {
		panic(err)
	}

	// 打印账户余额
	fmt.Printf("Account: %s\n", address)
	for _, coin := range res.Balances {
		fmt.Printf("Denom: %s, Amount: %s\n", coin.Denom, coin.Amount)
	}
}

type BlockResponse struct {
	Block struct {
		Header struct {
			Height string `json:"height"`
			Time   string `json:"time"`
		} `json:"header"`
		Data struct {
			Txs []string `json:"txs"`
		} `json:"data"`
	} `json:"block"`
}

type DecodeTxRequest struct {
	Tx string `json:"tx_bytes"`
}

type DecodeTxResponse struct {
	Tx struct {
		Body struct {
			Messages []banktypes.MsgSend `json:"messages"`
		} `json:"body"`
	} `json:"tx"`
}

func getLatestBlockHeight(restURL string) string {
	resp, err := http.Get(restURL + "/cosmos/base/tendermint/v1beta1/blocks/latest")
	if err != nil {
		log.Error("failed to get latest block height: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("failed to read response body: %v", err)
	}

	var blockResp BlockResponse
	err = json.Unmarshal(body, &blockResp)
	if err != nil {
		log.Error("failed to unmarshal response: %v", err)
	}

	return blockResp.Block.Header.Height
}

func getBlock(restURL, height string) BlockResponse {
	resp, err := http.Get(fmt.Sprintf("%s/cosmos/base/tendermint/v1beta1/blocks/%s", restURL, height))
	if err != nil {
		log.Error("failed to get block: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("failed to read response body: %v", err)
	}

	var blockResp BlockResponse
	err = json.Unmarshal(body, &blockResp)
	if err != nil {
		log.Error("failed to unmarshal response: %v", err)
	}

	return blockResp
}

func parseTx(txBytes string) {
	// 创建一个 codec 用于解码交易
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	// 解码交易
	txRaw, err := tx.DefaultTxDecoder(cdc)([]byte(txBytes))
	if err != nil {
		log.Error("failed to decode transaction: %v", err)
	}

	// 打印交易信息
	//fmt.Printf("Transaction Hash: %X\n", bytes.HexBytes(txBytes).Hash())
	fmt.Printf("Transaction Type: %T\n", txRaw)

	// 解析交易中的消息
	for _, msg := range txRaw.GetMsgs() {
		fmt.Printf("Message Type: %T\n", msg)
		fmt.Printf("Message Data: %v\n", msg)
	}
}

func decodeTx(restURL, txData string) {
	// 创建解码交易请求
	decodeTxReq := DecodeTxRequest{Tx: txData}
	jsonReq, err := json.Marshal(decodeTxReq)
	if err != nil {
		log.Error("failed to marshal decode tx request: %v", err)
	}

	// 发送 POST 请求到 tx/v1beta1/decode 端点
	resp, err := http.Post(restURL+"/cosmos/tx/v1beta1/decode", "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Error("failed to send decode tx request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("failed to read decode tx response: %v", err)
	}

	// 解析响应
	var decodeTxResp DecodeTxResponse
	err = json.Unmarshal(body, &decodeTxResp)
	if err != nil {
		log.Error("failed to unmarshal decode tx response: %v", err)
	}

	// 打印解码后的交易信息
	fmt.Printf("Decoded Transaction:\n")
	for i, msg := range decodeTxResp.Tx.Body.Messages {
		fmt.Printf("Message %d: %s\n", i+1, msg.String())
	}
}

//func getBlock(tmRPC *http.HTTP, height int64) {
//	// 获取区块信息
//	block, err := tmRPC.Block(context.Background(), &height)
//	if err != nil {
//		panic(err)
//	}
//
//	// 打印区块信息
//	fmt.Printf("Block Height: %d\n", block.Block.Height)
//	fmt.Printf("Block Hash: %s\n", block.Block.Hash())
//	fmt.Printf("Block Time: %s\n", block.Block.Time)
//}
//
//func getTransaction(tmRPC *http.HTTP, txHash string) {
//	// 获取交易信息
//	txResult, err := tmRPC.Tx(context.Background(), []byte(txHash), true)
//	if err != nil {
//		panic(err)
//	}
//
//	// 打印交易信息
//	fmt.Printf("Transaction Hash: %s\n", txResult.Hash)
//	fmt.Printf("Height: %d\n", txResult.Height)
//	fmt.Printf("Index: %d\n", txResult.Index)
//	fmt.Printf("TxResult: %+v\n", txResult.TxResult)
//	fmt.Printf("Tx: %s\n", txResult.Tx)
//}
