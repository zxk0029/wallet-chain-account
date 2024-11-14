package solana

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"

	account2 "github.com/dapplink-labs/chain-explorer-api/common/account"
	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"
	"github.com/gagliardetto/solana-go"
)

const ChainName = "Solana"

type ChainAdaptor struct {
	solCli    SolClient
	sdkClient *rpc.Client
	solData   *SolData
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	rpcUrl := conf.WalletNode.Sol.RpcUrl

	solHttpCli, err := NewSolHttpClient(rpcUrl)
	if err != nil {
		return nil, err
	}
	dataApiUrl := conf.WalletNode.Sol.DataApiUrl
	dataApiKey := conf.WalletNode.Sol.DataApiKey
	dataApiTimeOut := conf.WalletNode.Sol.TimeOut
	solData, err := NewSolScanClient(dataApiUrl, dataApiKey, time.Duration(dataApiTimeOut))
	if err != nil {
		return nil, err
	}

	sdkClient := rpc.New(rpcUrl)

	return &ChainAdaptor{
		solCli:    solHttpCli,
		sdkClient: sdkClient,
		solData:   solData,
	}, nil
}

func (c *ChainAdaptor) GetSupportChains(req *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	response := &account.SupportChainsResponse{
		Code:    common2.ReturnCode_ERROR,
		Msg:     "",
		Support: false,
	}
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		err := fmt.Errorf("GetSupportChains validateChainAndNetwork fail, err msg = %s", msg)
		log.Error("err", err)
		response.Msg = err.Error()
		return response, err
	}

	response.Msg = "Support this chain"
	response.Code = common2.ReturnCode_SUCCESS
	response.Support = true
	return response, nil
}

func (c *ChainAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	response := &account.ConvertAddressResponse{
		Code:    common2.ReturnCode_ERROR,
		Msg:     "",
		Address: "",
	}
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		err := fmt.Errorf("ConvertAddress validateChainAndNetwork fail, err msg = %s", msg)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	pubKeyHex := req.PublicKey
	if ok, msg := validatePublicKey(pubKeyHex); !ok {
		err := fmt.Errorf("ConvertAddress validatePublicKey fail, err msg = %s", msg)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	accountAddress, err := PubKeyHexToAddress(pubKeyHex)
	if err != nil {
		err := fmt.Errorf("ConvertAddress PubKeyHexToAddress failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	response.Code = common2.ReturnCode_SUCCESS
	response.Msg = "convert address success"
	response.Address = accountAddress
	return response, nil
}

func (c *ChainAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	response := &account.ValidAddressResponse{
		Code:  common2.ReturnCode_ERROR,
		Msg:   "",
		Valid: false,
	}

	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		err := fmt.Errorf("ValidAddress validateChainAndNetwork failed: %s", msg)
		log.Error("err", err)
		response.Msg = err.Error()
		return response, err
	}
	address := req.Address
	if len(address) == 0 {
		err := fmt.Errorf("ValidAddress address is empty")
		log.Error("err", err)
		response.Msg = err.Error()
		return response, err
	}
	if len(address) != 43 && len(address) != 44 {
		err := fmt.Errorf("invalid Solana address length: expected 43 or 44 characters, got %d", len(address))
		response.Msg = err.Error()
		return response, err
	}
	response.Code = common2.ReturnCode_SUCCESS
	response.Valid = true
	return response, nil
}

func (c *ChainAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	response := &account.BlockResponse{
		Code:         common2.ReturnCode_ERROR,
		Msg:          "",
		Height:       0,
		Hash:         "",
		BaseFee:      "",
		Transactions: nil,
	}
	if ok, msg := validateChainAndNetwork(req.Chain, ""); !ok {
		err := fmt.Errorf("GetBlockByNumber validateChainAndNetwork failed: %s", msg)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	resultSlot := uint64(req.Height)
	if req.Height == 0 {
		latestSlot, err := c.solCli.GetSlot(Finalized)
		if err != nil {
			err := fmt.Errorf("GetBlockByNumber GetSlot failed: %w", err)
			log.Error("err", err)
			response.Msg = err.Error()
			return nil, err
		}
		resultSlot = latestSlot
	}

	blockResult := &BlockResult{}
	if req.ViewTx {
		tempBlockBySlot, err := c.solCli.GetBlockBySlot(resultSlot, Signatures)
		if err != nil {
			err := fmt.Errorf("GetBlockByNumber GetBlockBySlot failed: %w", err)
			log.Error("err", err)
			response.Msg = err.Error()
			return nil, err
		}
		blockResult = tempBlockBySlot
	} else {
		tempBlockBySlot, err := c.solCli.GetBlockBySlot(resultSlot, None)
		if err != nil {
			err := fmt.Errorf("GetBlockByNumber GetBlockBySlot failed: %w", err)
			log.Error("err", err)
			response.Msg = err.Error()
			return nil, err
		}
		blockResult = tempBlockBySlot
	}

	response.Hash = blockResult.BlockHash
	response.Height = int64(resultSlot)
	response.Code = common2.ReturnCode_SUCCESS
	response.Msg = "GetBlockByNumber success"
	if req.ViewTx {
		response.Transactions = make([]*account.BlockInfoTransactionList, 0, len(blockResult.Signatures))
		for _, signature := range blockResult.Signatures {
			txInfo := &account.BlockInfoTransactionList{
				Hash: signature,
			}
			response.Transactions = append(response.Transactions, txInfo)
		}
	}
	return response, nil
}

func (c *ChainAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	response := &account.BlockResponse{
		Code:         common2.ReturnCode_ERROR,
		Msg:          "",
		Height:       0,
		Hash:         "",
		BaseFee:      "",
		Transactions: nil,
	}
	if ok, msg := validateChainAndNetwork(req.Chain, ""); !ok {
		response.Msg = msg
		err := fmt.Errorf("GetBlockByHash validateChainAndNetwork fail, err msg = %s", msg)
		return response, err
	}
	panic("implement me")
}

func (c *ChainAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	response := &account.BlockHeaderResponse{
		Code:        common2.ReturnCode_ERROR,
		Msg:         "",
		BlockHeader: nil,
	}

	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		response.Msg = msg
		err := fmt.Errorf("GetBlockHeaderByHash validateChainAndNetwork fail, err msg = %s", msg)
		return response, err
	}
	panic("implement me")
}

func (c *ChainAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	response := &account.BlockHeaderResponse{
		Code:        common2.ReturnCode_ERROR,
		Msg:         "",
		BlockHeader: nil,
	}
	if ok, msg := validateChainAndNetwork(req.Chain, ""); !ok {
		err := fmt.Errorf("GetBlockHeaderByNumber validateChainAndNetwork failed: %s", msg)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}

	resultSlot := uint64(req.Height)
	if req.Height == 0 {
		latestSlot, err := c.solCli.GetSlot(Finalized)
		if err != nil {
			err := fmt.Errorf("GetBlockHeaderByNumber GetSlot failed: %w", err)
			log.Error("err", err)
			response.Msg = err.Error()
			return nil, err
		}
		resultSlot = latestSlot
	}

	blockResult, err := c.solCli.GetBlockBySlot(resultSlot, None)
	if err != nil {
		err := fmt.Errorf("GetBlockHeaderByNumber GetBlockBySlot failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	blockHead := &account.BlockHeader{
		Hash:       blockResult.BlockHash,
		Number:     strconv.FormatUint(resultSlot, 10),
		ParentHash: blockResult.PreviousBlockhash,
		Time:       uint64(blockResult.BlockTime),
	}

	response.BlockHeader = blockHead
	response.Code = common2.ReturnCode_SUCCESS
	response.Msg = "GetBlockHeaderByNumber success"
	return response, nil
}

func (c *ChainAdaptor) GetAccount(req *account.AccountRequest) (*account.AccountResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	//TxResponse, err := c.solCli.SendTx(req.RawTx)
	//if err != nil {
	//	return &account.SendTxResponse{
	//		Code:   common2.ReturnCode_ERROR,
	//		Msg:    "get tx response error",
	//		TxHash: "0",
	//	}, nil
	//}
	return &account.SendTxResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "get tx response success",
		//TxHash: TxResponse,
	}, nil
}

func (c *ChainAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {
	var resp *account2.TransactionResponse[account2.AccountTxResponse]
	var err error
	fmt.Println("req.ContractAddress", req.ContractAddress)
	if req.ContractAddress != "0x00" && req.ContractAddress != "" {
		log.Info("Spl token transfer record")
		resp, err = c.solData.GetTxByAddress(uint64(req.Page), uint64(req.Pagesize), req.Address, "spl")
	} else {
		log.Info("Sol transfer record")
		resp, err = c.solData.GetTxByAddress(uint64(req.Page), uint64(req.Pagesize), req.Address, "sol")
	}
	if err != nil {
		log.Error("get GetTxByAddress error", "err", err)
		return &account.TxAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get tx list fail",
			Tx:   nil,
		}, err
	} else {
		txs := resp.TransactionList
		list := make([]*account.TxMessage, 0, len(txs))
		for i := 0; i < len(txs); i++ {
			list = append(list, &account.TxMessage{
				Hash:   txs[i].TxId,
				Tos:    []*account.Address{{Address: txs[i].To}},
				Froms:  []*account.Address{{Address: txs[i].From}},
				Fee:    txs[i].TxId,
				Status: account.TxStatus_Success,
				Values: []*account.Value{{Value: txs[i].Amount}},
				Type:   1,
				Height: txs[i].Height,
			})
		}
		fmt.Println("resp", resp)
		return &account.TxAddressResponse{
			Code: common2.ReturnCode_SUCCESS,
			Msg:  "get tx list success",
			Tx:   list,
		}, nil
	}
}

func (c *ChainAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {
	//tx, err := c.solCli.GetTxByHash(req.Hash)
	//if err != nil {
	//	return &account.TxHashResponse{
	//		Code: common2.ReturnCode_ERROR,
	//		Msg:  err.Error(),
	//		Tx:   nil,
	//	}, err
	//}
	//var value_list []*account.Value
	//value_list = append(value_list, &account.Value{Value: tx.Value})
	return &account.TxHashResponse{
		Tx: &account.TxMessage{
			//Hash:  tx.Hash,
			//Tos:   []*account.Address{{Address: tx.To}},
			//Froms: []*account.Address{{Address: tx.From}},

			//Fee:    tx.Fee,
			//Status: account.TxStatus_Success,
			//Values: value_list,
			//Type:   tx.Type,
			//Height: tx.Height,
		},
	}, nil
}

func (c *ChainAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	//TODO implement me
	panic("implement me")
}
func (c *ChainAdaptor) CreateUnSignTransaction(req *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {

	jsonBytes, err := base64.StdEncoding.DecodeString(req.Base64Tx)
	if err != nil {
		log.Error("decode string fail", "err", err)
		return nil, err
	}
	var data TxStructure
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		log.Error("parse json fail", "err", err)
		return nil, err
	}
	valueFloat, err := strconv.ParseFloat(data.Value, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse value: %w", err)
	}
	value := uint64(valueFloat * 1000000000)
	if err != nil {
		return nil, err
	}
	fromPubkey, err := solana.PublicKeyFromBase58(data.FromAddress)
	if err != nil {
		return nil, err
	}
	toPubkey, err := solana.PublicKeyFromBase58(data.ToAddress)
	if err != nil {
		return nil, err
	}
	var tx *solana.Transaction
	if isSOLTransfer(data.ContractAddress) {
		tx, err = solana.NewTransaction(
			[]solana.Instruction{
				system.NewTransferInstruction(
					value,
					fromPubkey,
					toPubkey,
				).Build(),
			},
			solana.MustHashFromBase58(data.Nonce),
			solana.TransactionPayer(fromPubkey),
		)

	} else {
		// SPL Token 转账
		mintPubkey := solana.MustPublicKeyFromBase58(data.ContractAddress)

		// 获取或创建发送方的代币账户
		fromTokenAccount, _, err := solana.FindAssociatedTokenAddress(
			fromPubkey,
			mintPubkey,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to find from token account: %w", err)
		}

		// 获取或创建接收方的代币账户
		toTokenAccount, _, err := solana.FindAssociatedTokenAddress(
			toPubkey,
			mintPubkey,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to find to token account: %w", err)
		}

		tx, err = solana.NewTransaction(
			[]solana.Instruction{
				token.NewTransferInstruction(
					value,
					fromTokenAccount, // 使用找到的代币账户
					toTokenAccount,   // 使用找到的代币账户
					fromPubkey,
					[]solana.PublicKey{},
				).Build(),
			},
			solana.MustHashFromBase58(data.Nonce),
			solana.TransactionPayer(fromPubkey),
		)
	}

	//https://github.com/gagliardetto/solana-go/tree/main?tab=readme-ov-file#transfer-sol-from-one-wallet-to-another-wallet
	return &account.UnSignTransactionResponse{
		Code:     common2.ReturnCode_SUCCESS,
		Msg:      "create un sign tx success",
		UnSignTx: tx.String(),
	}, nil
}
func (c ChainAdaptor) BuildSignedTransaction(req *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	jsonBytes, err := base64.StdEncoding.DecodeString(req.Base64Tx)
	if err != nil {
		log.Error("decode string fail", "err", err)
		return nil, err
	}
	var data TxStructure
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		log.Error("parse json fail", "err", err)
		return nil, err
	}
	valueFloat, err := strconv.ParseFloat(data.Value, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse value: %w", err)
	}
	value := uint64(valueFloat * 1000000000)
	if err != nil {
		return nil, err
	}
	fromPubkey, err := solana.PublicKeyFromBase58(data.FromAddress)
	if err != nil {
		return nil, err
	}
	privateKeyBytes, err := hex.DecodeString(data.FromPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}
	fromPrikey := solana.PrivateKey(privateKeyBytes)

	toPubkey, err := solana.PublicKeyFromBase58(data.ToAddress)
	if err != nil {
		return nil, err
	}
	var tx *solana.Transaction
	if isSOLTransfer(data.ContractAddress) {
		tx, err = solana.NewTransaction(
			[]solana.Instruction{
				system.NewTransferInstruction(
					value,
					fromPubkey,
					toPubkey,
				).Build(),
			},
			solana.MustHashFromBase58(data.Nonce),
			solana.TransactionPayer(fromPubkey),
		)

	} else {
		// SPL Token 转账
		mintPubkey := solana.MustPublicKeyFromBase58(data.ContractAddress)

		// 获取或创建发送方的代币账户
		fromTokenAccount, _, err := solana.FindAssociatedTokenAddress(
			fromPubkey,
			mintPubkey,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to find from token account: %w", err)
		}

		// 获取或创建接收方的代币账户
		toTokenAccount, _, err := solana.FindAssociatedTokenAddress(
			toPubkey,
			mintPubkey,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to find to token account: %w", err)
		}

		tx, err = solana.NewTransaction(
			[]solana.Instruction{
				token.NewTransferInstruction(
					value,
					fromTokenAccount, // 使用找到的代币账户
					toTokenAccount,   // 使用找到的代币账户
					fromPubkey,
					[]solana.PublicKey{},
				).Build(),
			},
			solana.MustHashFromBase58(data.Nonce),
			solana.TransactionPayer(fromPubkey),
		)
	}

	//https://github.com/gagliardetto/solana-go/tree/main?tab=readme-ov-file#transfer-sol-from-one-wallet-to-another-wallet
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			return &fromPrikey
		},
	)

	return &account.SignedTransactionResponse{
		Code:     common2.ReturnCode_SUCCESS,
		Msg:      "create un sign tx success",
		SignedTx: tx.String(),
	}, nil
}

func (c *ChainAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	//TODO implement me
	panic("implement me")
}
func isSOLTransfer(coinAddress string) bool {
	// SOL 的 wrapped token address 或空字符串
	return coinAddress == "" ||
		coinAddress == "So11111111111111111111111111111111111111112"
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

func validatePublicKey(pubKey string) (bool, string) {
	if pubKey == "" {
		return false, "public key cannot be empty"
	}
	pubKeyWithoutPrefix := strings.TrimPrefix(pubKey, "0x")

	if len(pubKeyWithoutPrefix) != 64 {
		return false, "invalid public key length"
	}
	if _, err := hex.DecodeString(pubKeyWithoutPrefix); err != nil {
		return false, "invalid public key format: must be hex string"
	}

	return true, ""
}
