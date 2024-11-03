package cosmos

import (
	"context"
	"fmt"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/dapplink-labs/wallet-chain-account/common/retry"
	"github.com/dapplink-labs/wallet-chain-account/common/util"
	"time"
)

const (
	defaultDialTimeout    = 5 * time.Second
	defaultDialAttempts   = 5
	defaultRequestTimeout = 10 * time.Second
)

type CosmosClient interface {
	// block
	Block(height *int64) (*ctypes.ResultBlock, error)
	//
	Header(ctx context.Context, height *int64) (*ctypes.ResultHeader, error)
	//
	BlockByHash(hash []byte) (*ctypes.ResultBlock, error)
	//
	BlockchainInfo(minHeight, maxHeight int64) (*ctypes.ResultBlockchainInfo, error)
	//
	Tx(hash []byte, prove bool) (*ctypes.ResultTx, error)

	Close()
}

type CosmosContext struct {
	context client.Context
	rpcHttp *rpchttp.HTTP
}

func DialCosmosClient(ctx context.Context, nodeUrl string) (CosmosClient, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultDialTimeout)
	defer cancel()

	bOff := retry.Exponential()
	httpClient, err := retry.Do(ctx, defaultDialAttempts, bOff, func() (*rpchttp.HTTP, error) {
		if !util.IsURLAvailable(nodeUrl) {
			return nil, fmt.Errorf("address unavailable (%s)", nodeUrl)
		}
		client, err := client.NewClientFromNode(nodeUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to dial address (%s): %w", nodeUrl, err)
		}
		return client, nil
	})

	if err != nil {
		return nil, err
	}

	return &CosmosContext{context: client.Context{}.WithClient(httpClient), rpcHttp: httpClient}, nil
}

func (c *CosmosContext) Block(height *int64) (*ctypes.ResultBlock, error) {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	result, err := c.rpcHttp.Block(ctxwt, height)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *CosmosContext) Header(ctx context.Context, height *int64) (*ctypes.ResultHeader, error) {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	result, err := c.rpcHttp.Header(ctxwt, height)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *CosmosContext) BlockByHash(hash []byte) (*ctypes.ResultBlock, error) {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	result, err := c.rpcHttp.BlockByHash(ctxwt, hash)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *CosmosContext) BlockchainInfo(minHeight, maxHeight int64) (*ctypes.ResultBlockchainInfo, error) {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	result, err := c.rpcHttp.BlockchainInfo(ctxwt, minHeight, maxHeight)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *CosmosContext) Tx(hash []byte, prove bool) (*ctypes.ResultTx, error) {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	result, err := c.rpcHttp.Tx(ctxwt, hash, prove)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *CosmosContext) Close() {
	//c.rpcHttp.OnStop()
}
