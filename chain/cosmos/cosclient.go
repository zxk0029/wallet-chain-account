package cosmos

import (
	"context"
	authv1beta1 "cosmossdk.io/api/cosmos/auth/v1beta1"
	"fmt"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	cometypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	ed255192 "github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dapplink-labs/wallet-chain-account/common/retry"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	"github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"time"
)

const (
	defaultDialTimeout    = 5 * time.Second
	defaultDialAttempts   = 5
	defaultRequestTimeout = 10 * time.Second
)

//type CosmosClient interface {
//	GetAccount(ctx context.Context, addr string) (*authv1beta1.QueryAccountResponse, error)
//
//	GetBalance(ctx context.Context, coin, addr string) (*types.Coin, error)
//
//	GetTxByHash(ctx context.Context, hash string) (*tx.GetTxResponse, error)
//
//	GetTxByEvent(ctx context.Context, event []string, page, limit uint64) (*tx.GetTxsEventResponse, error)
//
//	SendTx(ctx context.Context, fromAddr, toAddr, coin string, amount int64) (*banktypes.MsgSendResponse, error)
//
//	GetAddressFromPubKey(key []byte) string
//
//	BroadcastTx(ctx context.Context, txByte []byte) (*tx.BroadcastTxResponse, error)
//	// block
//	GetBlock(ctx context.Context, height *int64) (*ctypes.ResultBlock, error)
//	////
//	//Header(ctx context.Context, height *int64) (*ctypes.ResultHeader, error)
//	////
//	//BlockByHash(hash []byte) (*ctypes.ResultBlock, error)
//	////
//	//BlockchainInfo(minHeight, maxHeight int64) (*ctypes.ResultBlockchainInfo, error)
//	////
//	//Tx(hash []byte, prove bool) (*ctypes.ResultTx, error)
//
//	Close() error
//}

type CosmosClient struct {
	context client.Context
	rpchttp *rpchttp.HTTP
	codec   *codec.ProtoCodec

	grpcConn        *grpc.ClientConn
	bankClient      banktypes.QueryClient
	txServiceClient sdktx.ServiceClient
	authClient      authv1beta1.QueryClient
	bankKeeper      keeper.BaseKeeper
}

func DialCosmosClient(ctx context.Context, nodeUrl string) (*CosmosClient, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultDialTimeout)
	defer cancel()

	bOff := retry.Exponential()
	connClient, err := retry.Do(ctx, defaultDialAttempts, bOff, func() (*CosmosClient, error) {
		//grpcConn, err := grpc.Dial(nodeUrl, grpc.WithInsecure())
		//if err != nil {
		//	log.Error("failed to retry dial nodeUrl (%s): %w", nodeUrl, err)
		//	return nil, err
		//}
		//return grpcConn, nil

		//if !IsURLAvailable(rpcUrl) {
		//	return nil, fmt.Errorf("address unavailable (%s)", rpcUrl)
		//}
		//
		ctx := client.Context{}
		conn, err := client.NewClientFromNode(nodeUrl)
		if err != nil {
			log.Error("failed to retry dial nodeUrl (%s): %w", nodeUrl, err)
			return nil, err
		}
		codec := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
		ctx = ctx.WithClient(conn).WithKeyring(keyring.NewInMemory(codec))
		return &CosmosClient{
			context:         ctx,
			rpchttp:         conn,
			codec:           codec,
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

func (c *CosmosClient) GetTxByHash(hash string) (*sdktx.GetTxResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	return c.txServiceClient.GetTx(ctx, &sdktx.GetTxRequest{Hash: hash})
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
	// todo check
	pub := ed255192.PubKey{Key: key}
	return pub.Address().String()
}

func (c *CosmosClient) BroadcastTx(txByte []byte) (*sdktx.BroadcastTxResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	return c.txServiceClient.BroadcastTx(ctx, &sdktx.BroadcastTxRequest{
		TxBytes: txByte,
		Mode:    sdktx.BroadcastMode_BROADCAST_MODE_SYNC,
	})
}

func (c *CosmosClient) GetBlock(height *int64) (*ctypes.ResultBlock, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	result, err := c.rpchttp.Block(ctx, height)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *CosmosClient) DecodeBlockTx(block *cometypes.Block) ([]*account.BlockInfoTransactionList, error) {
	var blockTransactions []*account.BlockInfoTransactionList

	// 获取区块的总 Gas 使用量
	//totalGas := uint64(0)
	//for _, txHash := range block.Txs {
	//	// 获取交易结果
	//	txResult, err := c.rpchttp.Tx(context.Background(), txHash.Hash(), false)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	// 累加 Gas 使用量
	//	totalGas += uint64(txResult.TxResult.GasUsed)
	//}

	txDecoder := authtx.DefaultTxDecoder(c.codec)
	for _, txBytes := range block.Txs {
		// 解码交易
		tx, err := txDecoder(txBytes)
		if err != nil {
			return nil, err
		}
		// 获取交易中的消息
		msgs := tx.GetMsgs()
		for i, msg := range msgs {
			log.Info("Message %d: %+v", i, msg)
			// 获取消息的具体信息
			switch msg := msg.(type) {
			case *banktypes.MsgSend:
				log.Info("Sender: %s", msg.FromAddress)
				log.Info("Receiver: %s", msg.ToAddress)
				log.Info("Amount: %s", msg.Amount)

				blockTransaction := &account.BlockInfoTransactionList{
					From:   msg.FromAddress,
					To:     msg.ToAddress,
					Hash:   string(txBytes.Hash()),
					Amount: msg.Amount.String(),
				}
				blockTransactions = append(blockTransactions, blockTransaction)
			// 其他消息类型的处理
			default:
				log.Info("Unknown message type: %T", msg)
			}
		}
	}
	return blockTransactions, nil
}

func (c *CosmosClient) GetHeaderByHeight(height *int64) (*ctypes.ResultHeader, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	result, err := c.rpchttp.Header(ctx, height)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *CosmosClient) GetHeaderByHash(hash []byte) (*ctypes.ResultHeader, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	result, err := c.rpchttp.HeaderByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *CosmosClient) GetBlockByHash(hash []byte) (*ctypes.ResultBlock, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	result, err := c.rpchttp.BlockByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return result, nil
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

func (c *CosmosClient) Tx(hash []byte, prove bool) (*ctypes.ResultTx, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	result, err := c.rpchttp.Tx(ctx, hash, prove)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *CosmosClient) Close() error {
	//c.rpcHttp.OnStop()
	//return c.grpcConn.Close()
	return nil
}
