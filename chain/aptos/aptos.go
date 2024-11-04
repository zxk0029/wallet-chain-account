package aptos

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/log"
	"golang.org/x/crypto/sha3"

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

func (c *ChainAdaptor) GetSupportChains(req *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.SupportChainsResponse{
			Code:    common2.ReturnCode_ERROR,
			Msg:     msg,
			Support: false,
		}, nil
	}
	return &account.SupportChainsResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "Support this chain",
		Support: true,
	}, nil
}

func (c *ChainAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.ConvertAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
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

func (c *ChainAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.ValidAddressResponse{
			Code:  common2.ReturnCode_ERROR,
			Msg:   msg,
			Valid: false,
		}, nil
	}
	if len(req.Address) != 66 || !strings.HasPrefix(req.Address, "0x") {
		return &account.ValidAddressResponse{
			Code:  common2.ReturnCode_SUCCESS,
			Msg:   "invalid address: wrong length or missing 0x prefix",
			Valid: false,
		}, nil
	}
	ok := regexp.MustCompile("^[0-9a-fA-F]{64}$").MatchString(req.Address[2:])
	if !ok {
		return &account.ValidAddressResponse{
			Code:  common2.ReturnCode_SUCCESS,
			Msg:   "invalid address: contains invalid characters",
			Valid: false,
		}, nil
	}
	if strings.TrimPrefix(req.Address, "0x") == strings.Repeat("0", 64) {
		return &account.ValidAddressResponse{
			Code:  common2.ReturnCode_SUCCESS,
			Msg:   "invalid address: cannot be all zeros",
			Valid: false,
		}, nil
	}
	return &account.ValidAddressResponse{
		Code:  common2.ReturnCode_SUCCESS,
		Msg:   "valid address",
		Valid: true,
	}, nil
}

func (c *ChainAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, ""); !ok {
		return &account.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
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

func (c *ChainAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, ""); !ok {
		return &account.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
	// no implement
	panic("implement me")
}

func (c *ChainAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.BlockHeaderResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
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
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.BlockHeaderResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
	// no implement
	panic("implement me")
}

func (c *ChainAdaptor) GetAccount(req *account.AccountRequest) (*account.AccountResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.AccountResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
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
		Network:  req.Network,
	}, nil
}

func (c *ChainAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.FeeResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
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

func (c *ChainAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.SendTxResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.TxAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
	transactionsPtr, err := c.aptosClient.GetTransactionByAddress(req.Address)
	if err != nil {
		log.Error("GetTxByAddress GetTransactionByAddress fail", "err", err)
		return &account.TxAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "GetTxByAddress GetTransactionByAddress fail",
		}, nil
	}
	if transactionsPtr == nil {
		return &account.TxAddressResponse{
			Code: common2.ReturnCode_SUCCESS,
			Msg:  "GetTxByAddress success but no transactions found",
			Tx:   []*account.TxMessage{},
		}, nil
	}

	transactions := *transactionsPtr
	var txMessages []*account.TxMessage

	for _, tx := range transactions {
		var txStatus account.TxStatus
		if tx.Success {
			txStatus = account.TxStatus_Success
		} else {
			txStatus = account.TxStatus_Failed
		}

		feeStatement := GetFeeStatementFromEvents(tx.Events)
		var totalFee uint64
		if feeStatement != nil {
			totalFee = CalculateGasFee(tx.GasUnitPrice, feeStatement.TotalChargeGasUnits,
				feeStatement.StorageFeeOctas, feeStatement.StorageFeeRefundOctas)
		} else {
			totalFee = tx.GasUsed * tx.GasUnitPrice
		}
		fromAddr := &account.Address{
			Address: tx.Sender,
		}
		txMessage := &account.TxMessage{
			Hash:  tx.Hash,
			Froms: []*account.Address{fromAddr},
			//TODO to
			Tos: []*account.Address{},
			//TODO Value
			Values: []*account.Value{},
			Fee:    strconv.FormatUint(totalFee, 10),
			Status: txStatus,
			Type:   0,
			Height: strconv.FormatUint(tx.Version, 10),
			// ContractAddress:
			Datetime: strconv.FormatUint(tx.Timestamp, 10),
			Data:     convertExtraInfo(tx),
		}
		txMessages = append(txMessages, txMessage)
	}

	return &account.TxAddressResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "GetTxByAddress success",
		Tx:   txMessages,
	}, nil
}

func (c *ChainAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
	tx, err := c.aptosClient.GetTransactionByHash(req.Hash)
	if err != nil {
		log.Error("GetTransactionByHash error", "err", err)
		return &account.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "GetTransactionByHash error",
		}, nil
	}
	//TODO fromAddrList toAddrsList
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

func (c *ChainAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.BlockByRangeResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
	startVersion, err := strconv.ParseUint(req.Start, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid start version: %w", err)
	}
	endVersion, err := strconv.ParseUint(req.End, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid end version: %w", err)
	}
	if startVersion > endVersion {
		return nil, fmt.Errorf("start version (%d) cannot be greater than end version (%d)", startVersion, endVersion)
	}
	txs, err := c.aptosClient.GetTransactionByVersionRange(startVersion, endVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	response := &account.BlockByRangeResponse{
		Code:        common2.ReturnCode_SUCCESS,
		Msg:         "GetBlockByRange success",
		BlockHeader: make([]*account.BlockHeader, 0, len(txs)),
	}
	for _, tx := range txs {
		blockHeader := &account.BlockHeader{
			Hash: tx.Hash,
			//ParentHash:  tx.StateRootHash,
			//Root:        tx.StateRootHash,
			TxHash:      tx.Hash,
			ReceiptHash: tx.EventRootHash,
			//Number:      tx.Version,
			GasLimit: tx.MaxGasAmount,
			GasUsed:  tx.GasUsed,
			Time:     tx.Timestamp,
			Extra:    convertExtraInfo(tx),
			Nonce:    strconv.FormatUint(tx.SequenceNumber, 10),
		}

		response.BlockHeader = append(response.BlockHeader, blockHeader)
	}
	return response, nil
}

func (c *ChainAdaptor) CreateUnSignTransaction(req *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.UnSignTransactionResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) BuildSignedTransaction(req *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.SignedTransactionResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.DecodeTransactionResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
	return &account.DecodeTransactionResponse{
		Code:     common2.ReturnCode_SUCCESS,
		Msg:      "verify tx success",
		Base64Tx: "0x000000",
	}, nil
}

func (c *ChainAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.VerifyTransactionResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
	return &account.VerifyTransactionResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "verify tx success",
		Verify: true,
	}, nil
}

func (c *ChainAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.ExtraDataResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
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

func convertExtraInfo(tx TransactionResponse) string {
	extraInfo := map[string]interface{}{
		"vm_status":             tx.VMStatus,
		"accumulator_root_hash": tx.AccumulatorRootHash,
		"changes":               tx.Changes,
		"signature":             tx.Signature,
		"events":                tx.Events,
		"payload":               tx.Payload,
		"success":               tx.Success,
	}

	extraJSON, err := json.Marshal(extraInfo)
	if err != nil {
		return ""
	}
	return string(extraJSON)
}

func validateChainAndNetwork(chain, network string) (bool, string) {
	if chain != ChainName {
		return false, "invalid chain"
	}
	//if network != NetworkMainnet && network != NetworkTestnet {
	//	return false, "invalid network"
	//}
	return true, ""
}
