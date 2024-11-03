package cosmos

import (
	"context"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

const ChainName = "Cosmos"

type CosmosChainAdaptor struct {
	client CosmosClient
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	cosmosClient, err := DialCosmosClient(context.Background(), conf.WalletNode.Cosmos.RPCs[0].RPCURL)
	if err != nil {
		return nil, err
	}
	return &CosmosChainAdaptor{
		client: cosmosClient
	}, nil
}


func (c CosmosChainAdaptor) GetSupportChains(req *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	return &account.SupportChainsResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "Support this chain",
		Support: true,
	}, nil
}

func (c CosmosChainAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	addressCommon := common.BytesToAddress(crypto.Keccak256(req.PublicKey[1:])[12:])
	return &account.ConvertAddressResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "convert address successs",
		Address: addressCommon.String(),
	}, nil
}

func (c CosmosChainAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	return nil, nil
}

func (c CosmosChainAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	block, err := c.cosmosClient.BlockByNumber(big.NewInt(req.Height))
	if err != nil {
		log.Error("block by number error", err)
		return &account.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "block by number error",
		}, nil
	}
	var txListRet []*account.BlockInfoTransactionList
	for _, v := range block.Transactions {
		bitlItem := &account.BlockInfoTransactionList{
			From:   "0x000",
			To:     v.To,
			Hash:   v.Hash,
			Time:   "0",
			Amount: "10",
			Fee:    "0",
			Status: "1",
		}
		txListRet = append(txListRet, bitlItem)
	}
	return &account.BlockResponse{
		Code:         common2.ReturnCode_SUCCESS,
		Msg:          "block by number success",
		Hash:         block.Hash.String(),
		BaseFee:      block.BaseFee,
		Transactions: txListRet,
	}, nil
}

func (c CosmosChainAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c CosmosChainAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c CosmosChainAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c CosmosChainAdaptor) GetAccount(req *account.AccountRequest) (*account.AccountResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c CosmosChainAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c CosmosChainAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c CosmosChainAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c CosmosChainAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c CosmosChainAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c CosmosChainAdaptor) CreateUnSignTransaction(req *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c CosmosChainAdaptor) BuildSignedTransaction(req *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c CosmosChainAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c CosmosChainAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c CosmosChainAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	//TODO implement me
	panic("implement me")
}
