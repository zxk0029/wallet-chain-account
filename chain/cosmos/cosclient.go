package cosmos

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	sdk256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/dapplink-labs/chain-explorer-api/common/chain"
	"io/ioutil"
	"net/http"
	"time"

	authv1beta1 "cosmossdk.io/api/cosmos/auth/v1beta1"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	cttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	autx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	commaccount "github.com/dapplink-labs/chain-explorer-api/common/account"
	"github.com/dapplink-labs/chain-explorer-api/explorer/oklink"
	"github.com/dapplink-labs/wallet-chain-account/common/helpers"
	"github.com/dapplink-labs/wallet-chain-account/common/retry"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	"github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"
)

const (
	defaultDialTimeout    = 30 * time.Second
	defaultDialAttempts   = 5
	defaultRequestTimeout = 10 * time.Second
)

type CosmosClient struct {
	context client.Context
	rpchttp *rpchttp.HTTP
	codec   *codec.ProtoCodec

	dataClient      *oklink.ChainExplorerAdaptor
	bankClient      banktypes.QueryClient
	txServiceClient sdktx.ServiceClient
	authClient      authv1beta1.QueryClient
	bankKeeper      keeper.BaseKeeper
}

type BlockResponse struct {
	BlockId struct {
		Hash string `json:"hash"`
	} `json:"block_id"`
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
	Response struct {
		Height    string `json:"height"`
		Txhash    string `json:"txhash"`
		GasWanted string `json:"gas_wanted"`
		GasUsed   string `json:"gas_used"`
	} `json:"tx_response"`
}

type TxResponse struct {
	//Tx struct {
	//	Body struct {
	//		Messages []string `json:"messages"`
	//	} `json:"body"`
	//} `json:"tx"`
	Response struct {
		Height    string     `json:"height"`
		Txhash    string     `json:"txhash"`
		GasWanted string     `json:"gas_wanted"`
		GasUsed   string     `json:"gas_used"`
		Events    []*TxEvent `json:"events"`
		Timestamp string     `json:"timestamp"`
	} `json:"tx_response"`
}

type TxEvent struct {
	Type       string              `json:"type"`
	Attributes []*TxEventAttribute `json:"attributes"`
}

type TxEventAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func DialCosmosClient(ctx context.Context, conf *config.Config) (*CosmosClient, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultDialTimeout)
	defer cancel()

	nodeUrl := conf.WalletNode.Cosmos.RpcUrl
	dataApiKey := conf.WalletNode.Cosmos.DataApiKey
	dataApiUrl := conf.WalletNode.Cosmos.DataApiUrl
	timeOut := conf.WalletNode.Cosmos.TimeOut
	bOff := retry.Exponential()
	connClient, err := retry.Do(ctx, defaultDialAttempts, bOff, func() (*CosmosClient, error) {
		if !helpers.IsURLAvailable(nodeUrl) {
			return nil, fmt.Errorf("address unavailable (%s)", nodeUrl)
		}

		ctx := client.Context{}
		conn, err := client.NewClientFromNode(nodeUrl)
		if err != nil {
			log.Error("failed to retry dial nodeUrl (%s): %w", nodeUrl, err)
			return nil, err
		}

		dataClient, err := oklink.NewChainExplorerAdaptor(dataApiKey, dataApiUrl+"/", false, time.Duration(timeOut))
		if err != nil {
			log.Error("failed cosmos  new chain explorer adaptor", "err", err)
			return nil, err
		}
		codec := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
		ctx = ctx.WithClient(conn).WithKeyring(keyring.NewInMemory(codec))
		return &CosmosClient{
			context: ctx,
			rpchttp: conn,
			codec:   codec,

			dataClient:      dataClient,
			bankClient:      banktypes.NewQueryClient(ctx),
			txServiceClient: sdktx.NewServiceClient(ctx),
			authClient:      authv1beta1.NewQueryClient(ctx),
		}, nil

	})
	if err != nil {
		log.Error("failed to dial nodeUrl (%s): %w", nodeUrl, err)
		return nil, err
	}

	return connClient, nil
}

func (c *CosmosClient) GetBalance(coin, addr string) (*sdk.Coin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	address, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, err
	}

	resp, err := c.bankClient.Balance(ctx, &banktypes.QueryBalanceRequest{
		Address: address.String(),
		Denom:   coin,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetBalance(), nil
}

func (c *CosmosClient) GetTxByHash(hash string) (*ctypes.ResultTx, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		log.Error("failed to get block by hash: %v", err)
		return nil, err
	}

	return c.rpchttp.Tx(ctx, hashBytes, true)
}

func (c *CosmosClient) GetTxByEvent(event []string, page, limit uint64) (*sdktx.GetTxsEventResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 1000
	}

	if len(event) == 0 {
		return nil, errors.New("event cannot be empty")
	}
	eventTmp := make([]string, 0)
	for _, v := range event {
		_ = v
		eventTmp = append(eventTmp, fmt.Sprintf("message.sender='%s'", "cosmos1qd99t24whd3hfg22r53x8uw9ps3rrctwxqvn4m"))
	}

	tmp := &sdktx.GetTxsEventRequest{
		Events: eventTmp,
		Pagination: &query.PageRequest{
			Offset: page,
			Limit:  limit,
		},
	}

	return c.txServiceClient.GetTxsEvent(ctx, tmp)
}

func (c *CosmosClient) GetAccount(addr string) (*authv1beta1.QueryAccountResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	return c.authClient.Account(ctx, &authv1beta1.QueryAccountRequest{Address: addr})
}

func (c *CosmosClient) SendTx(fromAddr, toAddr, coin string, amount int64) (*banktypes.MsgSendResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()
	newCtx := sdk.UnwrapSDKContext(ctx)

	fromAddress, err := sdk.AccAddressFromBech32(fromAddr)
	if err != nil {
		return nil, err
	}
	toAddress, err := sdk.AccAddressFromBech32(toAddr)
	if err != nil {
		return nil, err
	}

	if c.bankKeeper.BlockedAddr(toAddress) {
		return nil, errors.New(fmt.Sprintf("%s is not allowed to receive funds", toAddress))
	}

	msg := banktypes.NewMsgSend(fromAddress, toAddress, sdk.NewCoins(sdk.NewInt64Coin(coin, amount)))
	if err := c.bankKeeper.SendCoins(newCtx, fromAddress, toAddress, msg.Amount); err != nil {
		return nil, err
	}

	//defer func() {
	//	for _, a := range msg.Amount {
	//		if a.Amount.IsInt64() {
	//			telemetry.SetGaugeWithLabels(
	//				[]string{"tx", "msg", "send"},
	//				float32(a.Amount.Int64()),
	//				[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
	//			)
	//		}
	//	}
	//}()

	return &banktypes.MsgSendResponse{}, nil
}

func (c *CosmosClient) GetAddressFromPubKey(key []byte) string {
	pubKey := &sdk256k1.PubKey{Key: key}
	accAddress := sdk.AccAddress(pubKey.Address())
	prefix := "cosmos"
	address, _ := bech32.ConvertAndEncode(prefix, accAddress)
	return address
}

func (c *CosmosClient) BroadcastTx(txByte []byte) (*sdktx.BroadcastTxResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	return c.txServiceClient.BroadcastTx(ctx, &sdktx.BroadcastTxRequest{
		TxBytes: txByte,
		Mode:    sdktx.BroadcastMode_BROADCAST_MODE_SYNC,
	})
}

func (c *CosmosClient) GetBlock(height int64) (*ctypes.ResultBlock, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	return c.rpchttp.Block(ctx, &height)
}

func (c *CosmosClient) TxParse(txData cttypes.Tx) {

	// 创建一个 authTx.TxDecoder
	txDecoder := autx.DefaultTxDecoder(c.codec)

	txStr := "CpQBCpEBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEnEKLWNvc21vczFya2V3ZXhod3dubWFnNXgzczZscnpxMjRzdWc0andmMDc5dWEyeRItY29zbW9zMWE5NjU3NmE3dHh4cWozOTJwNDAzZGcyZXRycHMwam14aGpmNWp3GhEKBXVhdG9tEggxNjcyMzM4OBJlCk4KRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECC6cnNR8l9z9mQsK/7jKyYHb6sXY4Kv37HZD91Ujoq+MSBAoCCH8SEwoNCgV1YXRvbRIEMzUzNRDJ0AgaQIMMnx49o57CM6l5vd/b/32ciczK0CuTuTv2qFA6DOvNMMI13sT2IqnIA8YVAMUXVOEEb/UmeluqFm7njlzYSIQ="
	txStr = txData.String()
	tx0, err := txDecoder([]byte(txStr))
	fmt.Printf("Tx0 Type: %T\n", tx0)

	decodeBytes1, _ := base64.StdEncoding.DecodeString(string(txStr))
	tx1, err := txDecoder(decodeBytes1)
	fmt.Printf("Tx1 Type: %T\n", tx1)

	decodeBytes2, _ := hex.DecodeString(string(txStr))
	tx2, err := txDecoder(decodeBytes2)
	fmt.Printf("Tx2 Type: %T\n", tx2)

	// 解码交易数据
	tx, err := txDecoder([]byte(txStr))
	if err != nil {
		log.Error("decoder tx fail!", err)
		return
	}

	// 打印交易信息
	//fmt.Printf("Tx Hash: %X\n", txBytes.
	fmt.Printf("Tx Type: %T\n", tx)

	// 解析交易中的各个字段
	for _, msg := range tx.GetMsgs() {
		fmt.Printf("Msg Type: %T\n", msg)

		switch msg := msg.(type) {
		case *banktypes.MsgSend:
			fmt.Printf("From: %s\n", msg.FromAddress)
			fmt.Printf("To: %s\n", msg.ToAddress)
			fmt.Printf("Amount: %s\n", msg.Amount)
		default:
			fmt.Println("Unknown Msg Type")
		}
	}
}

func (c *CosmosClient) GetBlockByHash(hash string) (*ctypes.ResultBlock, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}
	return c.rpchttp.BlockByHash(ctx, hashBytes)
}

func (c *CosmosClient) TxDecode(txData []byte) (*sdktx.TxDecodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	request := &sdktx.TxDecodeRequest{
		TxBytes: txData,
	}
	return c.txServiceClient.TxDecode(ctx, request)
}

func (c *CosmosClient) DecodeBlockTx(restURL string, block *BlockResponse) ([]*account.BlockInfoTransactionList, error) {
	var blockTransactions []*account.BlockInfoTransactionList

	for _, txData := range block.Block.Data.Txs {
		// 创建解码交易请求
		decodeTxReq := DecodeTxRequest{Tx: txData}
		jsonReq, err := json.Marshal(decodeTxReq)
		if err != nil {
			log.Error("failed to marshal decode tx request: %v", err)
			return nil, nil
		}

		// 发送 POST 请求到 tx/v1beta1/decode 端点
		resp, err := http.Post(restURL+"/cosmos/tx/v1beta1/decode", "application/json", bytes.NewBuffer(jsonReq))
		if err != nil {
			log.Error("failed to send decode tx request: %v", err)
			return nil, nil
		}
		// 读取响应
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error("failed to read decode tx response: %v", err)
			return nil, nil
		}
		resp.Body.Close()
		var decodeTxResp DecodeTxResponse
		err = json.Unmarshal(body, &decodeTxResp)
		if err != nil {
			log.Error("failed to unmarshal decode tx response: %v", err)
			return nil, nil
		}
		for _, msg := range decodeTxResp.Tx.Body.Messages {
			log.Info("Sender: %s", msg.FromAddress)
			log.Info("Receiver: %s", msg.ToAddress)
			log.Info("Amount: %s", msg.Amount.String())

			blockTransaction := &account.BlockInfoTransactionList{
				From:   msg.FromAddress,
				To:     msg.ToAddress,
				Hash:   decodeTxResp.Response.Txhash,
				Amount: msg.Amount.String(),
			}
			blockTransactions = append(blockTransactions, blockTransaction)
		}
	}
	return blockTransactions, nil
}

func (c *CosmosClient) GetHeaderByHeight(height int64) (*ctypes.ResultHeader, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	return c.rpchttp.Header(ctx, &height)
}

func (c *CosmosClient) GetHeaderByHash(hash string) (*ctypes.ResultHeader, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}
	return c.rpchttp.HeaderByHash(ctx, hashBytes)
}

func (c *CosmosClient) BlockchainInfo(minHeight, maxHeight int64) (*ctypes.ResultBlockchainInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	result, err := c.rpchttp.BlockchainInfo(ctx, minHeight, maxHeight)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *CosmosClient) Tx(hash string, prove bool) (*ctypes.ResultTx, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}

	return c.rpchttp.Tx(ctx, hashBytes, prove)
}

func (c *CosmosClient) GetTxByAddress(page, pagesize uint64, address string, action commaccount.ActionType) (*commaccount.TransactionResponse[commaccount.AccountTxResponse], error) {
	request := &commaccount.AccountTxRequest{
		ChainShortName:   ChainName,
		ExplorerName:     oklink.ChainExplorerName,
		Action:           action,
		Address:          address,
		StartBlockHeight: 0,
		EndBlockHeight:   22990189,
		PageRequest: chain.PageRequest{
			Page:  page,
			Limit: pagesize,
		},
	}
	rep, err := c.dataClient.GetTxByAddress(request)
	return rep, err
}

func (c *CosmosClient) GetTx(hash string) (*sdktx.GetTxResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	req := &sdktx.GetTxRequest{
		Hash: hash,
	}
	return c.txServiceClient.GetTx(ctx, req)
}

func (c *CosmosClient) Close() error {
	c.rpchttp.OnStop()
	return nil
}
