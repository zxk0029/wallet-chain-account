package aptos

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/log"
	"golang.org/x/crypto/sha3"
	"strconv"

	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"
)

const ChainName = "Aptos"

type ChainAdaptor struct {
	aptosClient *RestyClient
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	rpcUrl := conf.WalletNode.Aptos.RPCs[0].RPCURL
	apiKey := conf.WalletNode.Aptos.DataApiKey
	aptosClient, err := NewAptosClient(rpcUrl, apiKey)
	if err != nil {
		return nil, err
	}
	return &ChainAdaptor{
		aptosClient: aptosClient,
	}, nil
}

func (c ChainAdaptor) GetSupportChains(req *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	return &account.SupportChainsResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "Support this chain",
		Support: true,
	}, nil
}

func (c ChainAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	publicKeyBytes, err := hex.DecodeString(req.PublicKey)
	if err != nil {
		log.Error("ConvertAddress DecodeString fail", "err", err)
		return &account.ConvertAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "ConvertAddress DecodeString fail",
		}, nil
	}

	hasher := sha3.New256()
	hasher.Write(publicKeyBytes)
	hash := hasher.Sum(nil)

	aptosAddress := "0x" + hex.EncodeToString(hash)

	return &account.ConvertAddressResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "convert address success",
		Address: aptosAddress,
	}, nil
}

func (c ChainAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	// no implement
	panic("implement me")
}

func (c *ChainAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	if req.Height == 0 {
		nodeInfo, err := c.aptosClient.GetNodeInfo()
		if err != nil {
			log.Error("GetBlockHeaderByNumber GetNodeInfo fail", "err", err)
			return &account.BlockHeaderResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  "GetBlockHeaderByNumber GetNodeInfo fail",
			}, nil
		}
		seconds := nodeInfo.LedgerTimestamp / 1_000_000
		blockHead := &account.BlockHeader{
			Number: strconv.FormatUint(nodeInfo.BlockHeight, 10),
			Time:   seconds,
		}
		return &account.BlockHeaderResponse{
			Code:        common2.ReturnCode_SUCCESS,
			Msg:         "GetBlockHeaderByNumber GetNodeInfo success",
			BlockHeader: blockHead,
		}, nil
	}

	blockResponse, err := c.aptosClient.GetBlockByHeight(uint64(req.Height))
	if err != nil {
		log.Error("GetBlockHeaderByNumber GetBlockByHeight fail", "err", err)
		return &account.BlockHeaderResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "GetBlockHeaderByNumber GetBlockByHeight fail",
		}, nil
	}
	blockHead := &account.BlockHeader{
		Hash:   blockResponse.BlockHash,
		Number: strconv.FormatUint(blockResponse.BlockHeight, 10),
		Time:   blockResponse.BlockTimestamp,
	}
	return &account.BlockHeaderResponse{
		Code:        common2.ReturnCode_SUCCESS,
		Msg:         "GetBlockHeaderByNumber GetBlockByHeight success",
		BlockHeader: blockHead,
	}, nil
}

func (c *ChainAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	// no implement
	panic("implement me")
}

func (c ChainAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	if req.Height == 0 {
		nodeInfo, err := c.aptosClient.GetNodeInfo()
		if err != nil {
			log.Error("GetBlockByNumber GetNodeInfo fail", "err", err)
			return &account.BlockResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  "GetBlockByNumber fail",
			}, nil
		}
		if req.ViewTx {

		}
		return &account.BlockResponse{
			Code:   common2.ReturnCode_SUCCESS,
			Msg:    "GetBlockByNumber GetNodeInfo success",
			Height: int64(nodeInfo.BlockHeight),
			// TODO: Transactionsdasda
			Transactions: nil,
		}, nil
	}

	blockResponse, err := c.aptosClient.GetBlockByHeight(uint64(req.Height))
	if err != nil {
		log.Error("GetBlockByNumber GetBlockByHeight fail", "err", err)
		return &account.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "GetBlockByNumber GetBlockByHeight fail",
		}, nil
	}
	if req.ViewTx {

	}
	return &account.BlockResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "GetBlockByNumber GetBlockByHeight success",
		Height: int64(blockResponse.BlockHeight),
		Hash:   blockResponse.BlockHash,
		// TODO: Transactionsdasda
		Transactions: nil,
	}, nil
}

func (c ChainAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	// no implement
	panic("implement me")
}

func (c ChainAdaptor) GetAccount(req *account.AccountRequest) (*account.AccountResponse, error) {
	accountResponse, err := c.aptosClient.GetAccount(req.Address)
	if err != nil {
		log.Error("GetAccount fail", "err", err)
		return &account.AccountResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "GetAccount fail",
		}, nil
	}
	return &account.AccountResponse{
		Code:     common2.ReturnCode_SUCCESS,
		Msg:      "get account response success",
		Sequence: strconv.FormatUint(accountResponse.SequenceNumber, 10),
	}, nil
}

func (c ChainAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {
	response, err := c.aptosClient.GetGasPrice()
	if err != nil {
		log.Error("GetFee fail", "err", err)
		return &account.FeeResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "GetFee fail",
		}, nil
	}
	return &account.FeeResponse{
		Code:      common2.ReturnCode_SUCCESS,
		Msg:       "GetFee success",
		SlowFee:   strconv.FormatUint(response.DeprioritizedGasEstimate, 10),
		NormalFee: strconv.FormatUint(response.GasEstimate, 10),
		FastFee:   strconv.FormatUint(response.PrioritizedGasEstimate, 10),
	}, nil
}

func (c ChainAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {
	tx, err := c.aptosClient.GetTransactionByHash(req.Hash)
	if err != nil {
		log.Error("GetTransactionByHash error", "err", err)
		return &account.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "GetTransactionByHash error",
		}, nil
	}
	var fromAddrList []*account.Address
	var toAddrsList []*account.Address
	var valueList []*account.Value
	var txStatus account.TxStatus
	if tx.Success {
		txStatus = account.TxStatus_Success
	} else {
		txStatus = account.TxStatus_Failed
	}

	feeStatement := GetFeeStatementFromEvents(tx.Events)
	totalFee := CalculateGasFee(tx.GasUnitPrice, feeStatement.TotalChargeGasUnits, feeStatement.StorageFeeOctas, feeStatement.StorageFeeRefundOctas)

	return &account.TxHashResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "GetTransactionByHash success",
		Tx: &account.TxMessage{
			Hash:   tx.Hash,
			Froms:  fromAddrList,
			Tos:    toAddrsList,
			Values: valueList,
			Fee:    strconv.FormatUint(totalFee, 10),
			Status: txStatus,
			Type:   0,
			//Height:          tx.,
			//ContractAddress: tx.To().String(),
			//Data: hexutils.BytesToHex(tx.Data()),
		},
	}, nil
}

func (c ChainAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) CreateUnSignTransaction(req *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) BuildSignedTransaction(req *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	return &account.DecodeTransactionResponse{
		Code:     common2.ReturnCode_SUCCESS,
		Msg:      "verify tx success",
		Base64Tx: "0x000000",
	}, nil
}

func (c ChainAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	return &account.VerifyTransactionResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "verify tx success",
		Verify: true,
	}, nil
}

func (c ChainAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	return &account.ExtraDataResponse{
		Code:  common2.ReturnCode_SUCCESS,
		Msg:   "get extra data success",
		Value: "not data",
	}, nil
}

func CalculateGasFee(gasUnitPrice, totalChargeGasUnits, storageFeeOctas, storageRefundOctas uint64) uint64 {
	// calc base gas fee
	gasFee := gasUnitPrice * totalChargeGasUnits

	// Storage Fee
	netStorageFee := storageFeeOctas - storageRefundOctas

	// totalFee
	totalFee := gasFee + netStorageFee

	return totalFee
}

func GetFeeStatementFromEvents(events []Event) *FeeStatement {
	for _, event := range events {
		if event.Type == "0x1::transaction_fee::FeeStatement" {
			return &event.Data
		}
	}
	return nil
}
