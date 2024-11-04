package cosmos

import (
	"context"
	authv1beta1 "cosmossdk.io/api/cosmos/auth/v1beta1"
	"encoding/json"
	types2 "github.com/cosmos/cosmos-sdk/types"
	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/golang/protobuf/ptypes"
	"strconv"
	"strings"
)

const ChainName = "Cosmos"
const NetWork = "mainnet"

type ChainAdaptor struct {
	client CosmosClient
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	cosmosClient, err := DialCosmosClient(context.Background(), conf.WalletNode.Cosmos.RPCs[0].RPCURL)
	if err != nil {
		log.Error("new chain adaptor error (%w)", err)
		return nil, err
	}
	return &ChainAdaptor{
		client: *cosmosClient,
	}, nil
}

func (c *ChainAdaptor) GetSupportChains(req *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	supportList := []string{"stake", "atom"}

	checkIf := func(s string) bool {
		for _, v := range supportList {
			if strings.EqualFold(v, s) {
				return true
			}
		}
		return false
	}

	return &account.SupportChainsResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Support: checkIf(req.Chain),
	}, nil
}

func (c *ChainAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	addr := c.client.GetAddressFromPubKey(req.PublicKey)

	return &account.ConvertAddressResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "convert address success",
		Address: addr,
	}, nil
}

func (c *ChainAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	_, err := types2.AccAddressFromBech32(req.Address)
	if err != nil {
		return &account.ValidAddressResponse{
			Code:  common2.ReturnCode_ERROR,
			Msg:   err.Error(),
			Valid: false,
		}, err
	}
	return &account.ValidAddressResponse{
		Valid: true,
	}, nil
}

func (c *ChainAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	block, err := c.client.GetBlock(&req.Height)
	if err != nil {
		log.Error("get block by number error (%w)", err)
		return nil, err
	}

	transactions, err := c.client.DecodeBlockTx(block.Block)
	if err != nil {
		log.Error("decode block tx error (%w)", err)
		return nil, err
	}
	// BaseFee
	return &account.BlockResponse{
		Code:         common2.ReturnCode_SUCCESS,
		Msg:          "get block by number success",
		Height:       block.Block.Height,
		Hash:         string(block.Block.Hash()),
		Transactions: transactions,
	}, nil
}

func (c *ChainAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	block, err := c.client.GetBlockByHash([]byte(req.Hash))
	if err != nil {
		log.Error("get block by hash error (%w)", err)
		return nil, err
	}
	transactions, err := c.client.DecodeBlockTx(block.Block)
	if err != nil {
		log.Error("decode block tx error (%w)", err)
		return nil, err
	}
	// BaseFee
	return &account.BlockResponse{
		Code:         common2.ReturnCode_SUCCESS,
		Msg:          "get block by hash success",
		Height:       block.Block.Height,
		Hash:         string(block.Block.Hash()),
		Transactions: transactions,
	}, nil
}

func (c *ChainAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	header, err := c.client.GetHeaderByHash([]byte(req.GetHash()))
	if err != nil {
		log.Error("get block header by hash error (%w)", err)
		return nil, err
	}
	// todo gas field
	blockHeader := &account.BlockHeader{
		TxHash:     header.Header.DataHash.String(),
		ParentHash: header.Header.AppHash.String(),
		Number:     strconv.FormatInt(header.Header.Height, 10),
		Time:       uint64(header.Header.Time.Unix()),
	}
	return &account.BlockHeaderResponse{
		Code:        common2.ReturnCode_SUCCESS,
		Msg:         "get block header by hash success",
		BlockHeader: blockHeader,
	}, nil
}

func (c *ChainAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	header, err := c.client.GetHeaderByHeight(&req.Height)
	if err != nil {
		log.Error("get block header by number error (%w)", err)
		return nil, err
	}
	// todo gas field
	blockHeader := &account.BlockHeader{
		TxHash:     header.Header.DataHash.String(),
		ParentHash: header.Header.AppHash.String(),
		Number:     strconv.FormatInt(header.Header.Height, 10),
		Time:       uint64(header.Header.Time.Unix()),
	}
	return &account.BlockHeaderResponse{
		Code:        common2.ReturnCode_SUCCESS,
		Msg:         "get block header by number success",
		BlockHeader: blockHeader,
	}, nil
}

func (c *ChainAdaptor) GetAccount(req *account.AccountRequest) (*account.AccountResponse, error) {
	response, err := c.client.GetAccount(req.Address)
	if err != nil {
		log.Error("get account error (%w)", err)
		return nil, err
	}

	authAccount := new(authv1beta1.BaseAccount)
	if err := ptypes.UnmarshalAny(response.Account, authAccount); err != nil {
		log.Error("get account error (%w)", err)
		return nil, err
	}
	return &account.AccountResponse{
		Code:          common2.ReturnCode_SUCCESS,
		Msg:           "get account success",
		Network:       NetWork,
		AccountNumber: authAccount.GetAddress(),
		Sequence:      strconv.FormatUint(authAccount.Sequence, 10),
		Balance:       strconv.FormatUint(authAccount.AccountNumber, 10),
	}, nil
}

func (c *ChainAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	ret, err := c.client.BroadcastTx([]byte(req.RawTx))
	if err != nil {
		return &account.SendTxResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "BroadcastTx fail",
		}, err
	}
	return &account.SendTxResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "send tx success",
		TxHash: ret.TxResponse.TxHash,
	}, nil
}

func (c *ChainAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {
	//TODO 需接第三方
	panic("implement me")
}

func (c *ChainAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {
	ret, err := c.client.GetTxByHash(req.Hash)
	if err != nil {
		return &account.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get tx by hash fail",
		}, err
	}

	var resp []*Tx
	if err := json.Unmarshal([]byte(ret.TxResponse.Logs.String()), &resp); err != nil {
		return nil, err
	}

	fromAddr, toAddr := "", ""
	for _, v := range resp[0].Events {
		if v.Type == "transfer" {
			for _, attr := range v.Attributes {
				if attr.Key == "recipient" {
					toAddr = attr.Value
				}
				if attr.Key == "sender" {
					fromAddr = attr.Value
				}
			}
		}
	}

	index := 0
	if resp[0] != nil {
		index = resp[0].MsgIndex
	}

	return &account.TxHashResponse{
		Tx: &account.TxMessage{
			Hash:            req.Hash,
			Index:           uint32(index),
			Froms:           []*account.Address{{Address: fromAddr}},
			Tos:             []*account.Address{{Address: toAddr}},
			Values:          nil,
			Fee:             strconv.FormatInt(ret.TxResponse.GasUsed, 10),
			Status:          account.TxStatus_Success,
			Type:            0,
			Height:          strconv.FormatInt(ret.TxResponse.Height, 10),
			ContractAddress: "",
			Datetime:        ret.TxResponse.Timestamp,
		},
	}, nil
}

func (c *ChainAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	minHeight, err := strconv.ParseInt(req.GetStart(), 10, 64)
	if err != nil {
		log.Error("min height invalid", err)
		return &account.BlockByRangeResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get block range fail ! start height err",
		}, err
	}
	maxHeight, err := strconv.ParseInt(req.GetEnd(), 10, 64)
	if err != nil {
		log.Error("max height invalid", err)
		return &account.BlockByRangeResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get block range fail ! start height err",
		}, err
	}
	blockInfo, err := c.client.BlockchainInfo(minHeight, maxHeight)
	log.Info("block metas len: %d", len(blockInfo.BlockMetas))
	return &account.BlockByRangeResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "get block by range success",
	}, nil
}

func (c *ChainAdaptor) CreateUnSignTransaction(req *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	return &account.UnSignTransactionResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "create unsigned transaction success",
	}, nil
}

func (c *ChainAdaptor) BuildSignedTransaction(req *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	return &account.SignedTransactionResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "build signed transaction success",
	}, nil
}

func (c *ChainAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	return &account.DecodeTransactionResponse{
		Code:     common2.ReturnCode_SUCCESS,
		Msg:      "decode transaction success",
		Base64Tx: "0x000000",
	}, nil
}

func (c *ChainAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	return &account.VerifyTransactionResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "verify signed transaction success",
		Verify: true,
	}, nil
}

func (c *ChainAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	return &account.ExtraDataResponse{
		Code:  common2.ReturnCode_SUCCESS,
		Msg:   "get extra data success",
		Value: "not data",
	}, nil
}

// /////////////////////////////////////////////////////////////////////////
type Tx struct {
	MsgIndex int        `json:"msg_index"`
	Events   []*TxEvent `json:"events"`
}

type TxEvent struct {
	Type       string              `json:"type"`
	Attributes []*TxEventAttribute `json:"attributes"`
}

type TxEventAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
