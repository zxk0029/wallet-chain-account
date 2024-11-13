package solana

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"strconv"
	"time"

	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/mr-tron/base58"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	//"github.com/gagliardetto/solana-go"

	account2 "github.com/dapplink-labs/chain-explorer-api/common/account"
	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"

	"encoding/json"
)

const ChainName = "Solana"

type ChainAdaptor struct {
	solCli  SolClient
	solData *SolData
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	cli, err := NewSolClient(conf)
	if err != nil {
		return nil, err
	}
	sol, err := NewSolScanClient(conf.WalletNode.Sol.DataApiUrl, conf.WalletNode.Sol.DataApiKey, time.Duration(conf.WalletNode.Sol.TimeOut))
	if err != nil {
		return nil, err
	}
	return &ChainAdaptor{
		solCli:  *cli,
		solData: sol,
	}, nil

}

func (c ChainAdaptor) GetSupportChains(req *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	return &account.SupportChainsResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "Support solana chain",
		Support: true,
	}, nil
}

func (c ChainAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	publicKeyBytes, err := hex.DecodeString(req.PublicKey)
	if err != nil {
		return &account.ConvertAddressResponse{
			Code:    common2.ReturnCode_ERROR,
			Msg:     "convert address fail",
			Address: common.Address{}.String(),
		}, nil
	}
	address := solana.PublicKeyFromBytes(publicKeyBytes)
	if err != nil {
		return &account.ConvertAddressResponse{
			Code:    common2.ReturnCode_ERROR,
			Msg:     "invalid public key",
			Address: common.Address{}.String(),
		}, nil
	}

	if !address.IsOnCurve() {
		return &account.ConvertAddressResponse{
			Code:    common2.ReturnCode_ERROR,
			Msg:     "public key is not on the curve",
			Address: common.Address{}.String(),
		}, nil
	}

	return &account.ConvertAddressResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "convert address success",
		Address: address.String(),
	}, nil
}

func (c ChainAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	if len(req.Address) == 0 {
		return &account.ValidAddressResponse{
			Code:  common2.ReturnCode_SUCCESS,
			Msg:   "invalid address: empty address",
			Valid: false,
		}, nil
	}

	decoded, err := base58.Decode(req.Address)
	if err != nil {
		return &account.ValidAddressResponse{
			Code:  common2.ReturnCode_SUCCESS,
			Msg:   "invalid address: not base58 encoded",
			Valid: false,
		}, nil
	}

	if len(decoded) != 32 {
		return &account.ValidAddressResponse{
			Code:  common2.ReturnCode_SUCCESS,
			Msg:   "invalid address: wrong length",
			Valid: false,
		}, nil
	}

	return &account.ValidAddressResponse{
		Code:  common2.ReturnCode_SUCCESS,
		Msg:   "valid address",
		Valid: true,
	}, nil
}

func (c ChainAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) GetAccount(req *account.AccountRequest) (*account.AccountResponse, error) {
	//req.ContractAddress as nonceAddress
	nonceResult, err := c.solCli.GetNonce(req.ContractAddress)
	if err != nil {
		log.Error("get nonce by address fail", "err", err)
		return &account.AccountResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get nonce by address fail",
		}, nil
	}
	balanceResult, err := c.solCli.GetBalance(req.Address)
	if err != nil {
		return &account.AccountResponse{
			Code:    common2.ReturnCode_ERROR,
			Msg:     "get token balance fail",
			Balance: "0",
		}, err
	}
	log.Info("balance result", "balance=", balanceResult, "balanceStr=", balanceResult)
	return &account.AccountResponse{
		Code:          common2.ReturnCode_SUCCESS,
		Msg:           "get account response success",
		AccountNumber: "0",
		Sequence:      nonceResult,
		Balance:       balanceResult,
	}, nil
}

func (c ChainAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	TxResponse, err := c.solCli.SendTx(req.RawTx)
	if err != nil {
		return &account.SendTxResponse{
			Code:   common2.ReturnCode_ERROR,
			Msg:    "get tx response error",
			TxHash: "0",
		}, nil
	}
	return &account.SendTxResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "get tx response success",
		TxHash: TxResponse,
	}, nil
}

func (c ChainAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {
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

func (c ChainAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {
	tx, err := c.solCli.GetTxByHash(req.Hash)
	if err != nil {
		return &account.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
			Tx:   nil,
		}, err
	}
	var value_list []*account.Value
	value_list = append(value_list, &account.Value{Value: tx.Value})
	return &account.TxHashResponse{
		Tx: &account.TxMessage{
			Hash:  tx.Hash,
			Tos:   []*account.Address{{Address: tx.To}},
			Froms: []*account.Address{{Address: tx.From}},

			Fee:    tx.Fee,
			Status: account.TxStatus_Success,
			Values: value_list,
			Type:   tx.Type,
			Height: tx.Height,
		},
	}, nil
}

func (c ChainAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
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

func (c ChainAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	//TODO implement me
	panic("implement me")
}
func isSOLTransfer(coinAddress string) bool {
	// SOL 的 wrapped token address 或空字符串
	return coinAddress == "" ||
		coinAddress == "So11111111111111111111111111111111111111112"
}
