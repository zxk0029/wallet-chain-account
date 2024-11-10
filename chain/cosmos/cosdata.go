package cosmos

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/dapplink-labs/wallet-chain-account/config"
)

type CosmosData struct {
	DataApiUrl   string
	DataApiKey   string
	DataApiToken string
	TimeOut      uint64
	Client       *http.Client
}

type NativeBalanceResponse struct {
	Code     string `json:"code"`
	Msg      string `json:"msg"`
	Response struct {
		Address                string `json:"address"`
		AvailableBalance       string `json:"availableBalance"`
		Delegated              string `json:"delegated"`
		DelegatedReward        string `json:"delegatedReward"`
		RewardRecipientAddress string `json:"rewardRecipientAddress"`
		Unbonding              string `json:"unbonding"`
		Symbol                 string `json:"symbol"`
		Commission             string `json:"commission"`
		Incentive              string `json:"incentive"`
		EthereumCoChainAddress string `json:"ethereumCoChainAddress"`
	} `json:"data"`
}

func NewCosmosData(conf *config.Config) (*CosmosData, error) {
	return &CosmosData{
		DataApiUrl:   conf.WalletNode.Cosmos.DataApiUrl,
		DataApiKey:   conf.WalletNode.Cosmos.DataApiKey,
		DataApiToken: conf.WalletNode.Cosmos.DataApiToken,
		TimeOut:      conf.WalletNode.Cosmos.TimeOut,
		Client:       &http.Client{},
	}, nil
}

func (d *CosmosData) GetThirdNativeBalance(address string) (*NativeBalanceResponse, error) {
	apiURL := fmt.Sprintf("%s/api/v5/explorer/address/balance-cosmos", d.DataApiUrl)
	u, err := url.ParseRequestURI(apiURL)
	if err != nil {
		log.Error("GetThirdNativeBalance NewRequest fail", "err", err)
		return nil, err
	}
	params := url.Values{}
	params.Add("address", address)
	params.Add("chainShortName", ChainName)
	u.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Error("GetThirdNativeBalance NewRequest fail", "err", err)
		return nil, err
	}

	// set header
	req.Header.Set("Ok-Access-Key", d.DataApiKey)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := d.Client.Do(req)
	if err != nil {
		log.Error("GetThirdNativeBalance Do fail", "err", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("GetThirdNativeBalance failed to read response body: %v", "err", err)
		return nil, err
	}

	var nativeBalanceResponse NativeBalanceResponse
	err = json.Unmarshal(body, &nativeBalanceResponse)
	if err != nil {
		log.Error("GetThirdNativeBalance failed to unmarshal nativeBalanceResponse: %v", "err", err)
		return nil, err
	}

	return &nativeBalanceResponse, nil
}

func GetThirdBlock() {

}

func GetThirdBlockTxs() {

}
