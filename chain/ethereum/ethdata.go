package ethereum

import (
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/dapplink-labs/chain-explorer-api/common/account"
	"github.com/dapplink-labs/chain-explorer-api/common/chain"
	"github.com/dapplink-labs/chain-explorer-api/explorer/etherscan"
)

type EthData struct {
	EthDataCli *etherscan.ChainExplorerAdaptor
}

func NewEthDataClient(baseUrl, apiKey string, timeout time.Duration) (*EthData, error) {
	etherscanCli, err := etherscan.NewChainExplorerAdaptor(apiKey, baseUrl, false, time.Duration(timeout))
	if err != nil {
		log.Error("New etherscan client fail", "err", err)
		return nil, err
	}
	return &EthData{EthDataCli: etherscanCli}, err
}

func (ss *EthData) GetTxByAddress(page, pagesize uint64, address string, action account.ActionType) (*account.TransactionResponse[account.AccountTxResponse], error) {
	request := &account.AccountTxRequest{
		PageRequest: chain.PageRequest{
			Page:  page,
			Limit: pagesize,
		},
		Action:  action,
		Address: address,
	}
	txData, err := ss.EthDataCli.GetTxByAddress(request)
	if err != nil {
		return nil, err
	}
	return txData, nil
}
