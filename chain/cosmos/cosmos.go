package cosmos

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	authv1beta1 "cosmossdk.io/api/cosmos/auth/v1beta1"
	"github.com/cometbft/cometbft/types"
	types2 "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/golang/protobuf/ptypes"

	account2 "github.com/dapplink-labs/chain-explorer-api/common/account"
	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"
)

const (
	NetWork   = "mainnet"
	ChainName = "Cosmos"
)

type ChainAdaptor struct {
	client  CosmosClient
	cosData *CosmosData
	conf    *config.Config
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	cosmosClient, err := DialCosmosClient(context.Background(), conf)
	if err != nil {
		log.Error("new chain adaptor error (%w)", err)
		return nil, err
	}
	cosmosData, err := NewCosmosData(conf)
	if err != nil {
		log.Error("new chain cosmos data error (%w)", err)
		return nil, err
	}
	return &ChainAdaptor{
		client:  *cosmosClient,
		conf:    conf,
		cosData: cosmosData,
	}, nil
}

func (c *ChainAdaptor) GetSupportChains(req *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	supportList := []string{"stake", "cosmos", "atom"}

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
		Msg:     "Support this chain",
		Support: checkIf(req.Chain),
	}, nil
}

func (c *ChainAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	pubKeyBytes, err := hex.DecodeString(req.PublicKey)
	if err != nil {
		log.Error("ConvertAddress failed to decode hex : %v", err)
		return nil, err
	}
	addr := c.client.GetAddressFromPubKey(pubKeyBytes)
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
		Code:  common2.ReturnCode_SUCCESS,
		Msg:   "valid address success",
		Valid: true,
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
	balance, err := c.cosData.GetThirdNativeBalance(req.Address)
	if err != nil {
		log.Error("get account error (%w)", err)
		return nil, err
	}

	return &account.AccountResponse{
		Code:          common2.ReturnCode_SUCCESS,
		Msg:           "get account success",
		Network:       NetWork,
		AccountNumber: strconv.FormatUint(authAccount.AccountNumber, 10),
		Sequence:      strconv.FormatUint(authAccount.Sequence, 10),
		Balance:       balance.Response.AvailableBalance,
	}, nil
}

func (c *ChainAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	block, err := c.client.GetBlock(req.Height)
	if err != nil {
		log.Error("get block by number error (%w)", err)
		return nil, err
	}

	totalGas, blockTransactions := c.parseTx(block.Block.Txs)
	return &account.BlockResponse{
		Code:         common2.ReturnCode_SUCCESS,
		Msg:          "get block by number success",
		Height:       block.Block.Height,
		BaseFee:      strconv.FormatUint(totalGas, 10),
		Hash:         block.BlockID.Hash.String(),
		Transactions: blockTransactions,
	}, nil

}

func (c *ChainAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	block, err := c.client.GetBlockByHash(req.GetHash())
	if err != nil {
		log.Error("get block by hash error (%w)", err)
		return nil, err
	}
	totalGas, blockTransactions := c.parseTx(block.Block.Txs)
	return &account.BlockResponse{
		Code:         common2.ReturnCode_SUCCESS,
		Msg:          "get block by hash success",
		Transactions: blockTransactions,
		BaseFee:      strconv.FormatUint(totalGas, 10),
		Height:       block.Block.Height,
		Hash:         block.Block.Hash().String(),
	}, nil
}

func (c *ChainAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	header, err := c.client.GetHeaderByHash(req.GetHash())
	if err != nil {
		log.Error("get block header by hash error (%w)", err)
		return nil, err
	}
	// todo gas field
	blockHeader := &account.BlockHeader{
		Hash:       header.Header.Hash().String(),
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
	header, err := c.client.GetHeaderByHeight(req.Height)
	if err != nil {
		log.Error("get block header by number error (%w)", err)
		return nil, err
	}
	// todo gas field
	blockHeader := &account.BlockHeader{
		Hash:       header.Header.Hash().String(),
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

// todo-底层url不对需要重新封装
func (c *ChainAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {
	var resp *account2.TransactionResponse[account2.AccountTxResponse]
	var err error
	if req.ContractAddress != "0x00" {
		resp, err = c.client.GetTxByAddress(uint64(req.Page), uint64(req.Pagesize), req.Address, "token")
		if err != nil {
			return nil, err
		}
	} else {
		resp, err = c.client.GetTxByAddress(uint64(req.Page), uint64(req.Pagesize), req.Address, "normal")
		if err != nil {
			return nil, err
		}
	}
	txs := resp.TransactionList
	list := make([]*account.TxMessage, 0, len(txs))
	for i := 0; i < len(txs); i++ {
		list = append(list, &account.TxMessage{
			Hash:   txs[i].TxId,
			Tos:    []*account.Address{{Address: txs[i].To}},
			Froms:  []*account.Address{{Address: txs[i].From}},
			Fee:    txs[i].TxFee,
			Status: account.TxStatus_Success,
			Values: []*account.Value{{Value: txs[i].Amount}},
			Type:   1,
			Height: txs[i].Height,
		})
	}
	return &account.TxAddressResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "get transactions by address success",
		Tx:   list,
	}, err
}

func (c *ChainAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {
	txResult, err := c.client.GetTx(req.GetHash())
	if err != nil {
		return &account.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get tx by hash fail",
		}, err
	}

	msgIndex, fromAddr, toAddr, amount := "", "", "", ""
	for _, event := range txResult.GetTxResponse().Events {
		if event.Type == "transfer" && len(event.GetAttributes()) == 4 {
			for _, attr := range event.Attributes {
				if attr.Key == "recipient" {
					toAddr = attr.Value
				}
				if attr.Key == "sender" {
					fromAddr = attr.Value
				}
				if attr.Key == "amount" {
					amount = attr.Value
				}
				if attr.Key == "msg_index" {
					msgIndex = attr.Value
				}
			}
		}
	}

	values := []*account.Value{{Value: amount}}
	index, _ := strconv.ParseUint(msgIndex, 10, 32)
	return &account.TxHashResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "get tx by hash success",
		Tx: &account.TxMessage{
			Hash:            txResult.GetTxResponse().TxHash,
			Index:           uint32(index),
			Froms:           []*account.Address{{Address: fromAddr}},
			Tos:             []*account.Address{{Address: toAddr}},
			Values:          values,
			Fee:             strconv.FormatInt(txResult.GetTxResponse().GasUsed, 10),
			Status:          account.TxStatus_Success,
			Type:            0,
			Height:          strconv.FormatInt(txResult.GetTxResponse().Height, 10),
			ContractAddress: "",
			Datetime:        txResult.GetTxResponse().Timestamp,
			Data:            txResult.GetTxResponse().Data,
		},
	}, nil
}

// todo
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

func (c *ChainAdaptor) parseTx(txs types.Txs) (uint64, []*account.BlockInfoTransactionList) {
	totalGas := uint64(0)
	var recipient, sender, amount = "", "", ""
	var blockTransactions []*account.BlockInfoTransactionList

	for _, txData := range txs {
		// hash ok
		txHash := fmt.Sprintf("%x", string(txData.Hash()))
		txResult, err := c.client.Tx(txHash, true)
		if err != nil {
			log.Error("get block by number error (%w)", err)
			continue
		}
		totalGas += uint64(txResult.TxResult.GasUsed)
		for _, event := range txResult.TxResult.Events {
			eventLen := len(event.GetAttributes())
			if event.Type == "transfer" && eventLen == 4 {
				for _, attr := range event.GetAttributes() {
					if attr.GetKey() == "recipient" {
						recipient = attr.GetValue()
					} else if attr.GetKey() == "sender" {
						sender = attr.GetValue()
					} else if attr.GetKey() == "amount" {
						amount = attr.GetValue()
					}
				}
				blockTransaction := &account.BlockInfoTransactionList{
					From:   sender,
					To:     recipient,
					Hash:   txHash,
					Amount: amount,
				}
				blockTransactions = append(blockTransactions, blockTransaction)
			}
		}
	}

	return totalGas, blockTransactions
}
