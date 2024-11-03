package aptos

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	gresty "github.com/go-resty/resty/v2"
)

type Client interface {
	GetNodeInfo() (*NodeInfo, error)

	GetAccount(inputAddr string) (*AccountResponse, error)
	GetGasPrice() (*EstimateGasPriceResponse, error)
	SubmitTransaction(req *SubmitTransactionRequest) (*SubmitTransactionResponse, error)
	GetTransactionByHash(txHash string) (*TransactionResponse, error)
	GetTransactionByAddress(address string) (*[]TransactionResponse, error)

	GetBlockByHeight(height uint64) (*BlockResponse, error)
}

var (
	errHTTPError       = errors.New("aptos http error")
	errInvalidAddress  = errors.New("invalid address")
	errInvalidResponse = errors.New("invalid response")
)

const (
	defaultRequestTimeout = 10 * time.Second
	defaultRetryCount     = 3
	defaultWithDebug      = false
	apikeyHeader          = "api-key"

	baseAPIPath = "/v1"

	// API Url
	pathNodeInfo = baseAPIPath + "/"

	pathGetSequence = baseAPIPath + "/accounts/%s"
	pathGasPrice    = baseAPIPath + "/estimate_gas_price"

	pathTransactions = baseAPIPath + "/transactions"
	pathTxByAddr     = baseAPIPath + "/accounts/%s/transactions"
	pathTxByHash     = baseAPIPath + "/transactions/by_hash/%s"

	pathBlockByHeight = baseAPIPath + "/blocks/by_height/%s"
)

type RestyClient struct {
	client *gresty.Client
}

func NewAptosClient(baseUrl, apiKey string) (*RestyClient, error) {
	return NewAptosClientAll(baseUrl, apiKey, defaultWithDebug)
}

func NewAptosClientAll(baseUrl, apiKey string, withDebug bool) (*RestyClient, error) {
	client := gresty.New()
	client.SetBaseURL(baseUrl)
	client.SetTimeout(defaultRequestTimeout)
	client.SetRetryCount(defaultRetryCount)
	client.SetDebug(withDebug)
	if apiKey != "" {
		client.SetHeader(apikeyHeader, apiKey)
	}

	client.OnAfterResponse(func(c *gresty.Client, r *gresty.Response) error {
		statusCode := r.StatusCode()
		if statusCode >= 400 {
			method := r.Request.Method
			url := r.Request.URL
			return fmt.Errorf("%d cannot %s %s: %w", statusCode, method, url, errHTTPError)
		}
		return nil
	})
	return &RestyClient{
		client: client,
	}, nil
}

func (c *RestyClient) GetNodeInfo() (*NodeInfo, error) {
	response := &NodeInfo{}
	resp, err := c.client.R().
		SetResult(response).
		Get(pathNodeInfo)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get node info: %w", errHTTPError)
	}

	return response, nil
}

// GetAccount info
func (c *RestyClient) GetAccount(inputAddr string) (*AccountResponse, error) {
	if !IsValidAddress(inputAddr) {
		return nil, fmt.Errorf("invalid address %s: %w", inputAddr, errInvalidAddress)
	}
	dealAddr := strings.TrimSpace(inputAddr)

	account := &AccountResponse{}
	resp, err := c.client.R().
		SetResult(account).
		Get(fmt.Sprintf(pathGetSequence, dealAddr))
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get sequence number: %w", errHTTPError)
	}
	return account, nil
}

// GetGasPrice Get estimate gas price
func (c *RestyClient) GetGasPrice() (*EstimateGasPriceResponse, error) {
	gasPrice := &EstimateGasPriceResponse{}
	resp, err := c.client.R().
		SetResult(gasPrice).
		Get(fmt.Sprintf(pathGasPrice))

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("failed to get gas price: %w", errHTTPError)
	}

	return gasPrice, nil
}

func (c *RestyClient) SubmitTransaction(req *SubmitTransactionRequest) (*SubmitTransactionResponse, error) {
	// check req
	if err := ValidateSubmitTransaction(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	response := &SubmitTransactionResponse{}
	resp, err := c.client.R().
		SetBody(req).
		SetResult(response).
		Post(pathTransactions)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to broadcast transaction: %w", errHTTPError)
	}

	return response, nil
}

func (c *RestyClient) GetTransactionByAddress(inputAddr string) (*[]TransactionResponse, error) {
	if !IsValidAddress(inputAddr) {
		return nil, fmt.Errorf("invalid address %s: %w", inputAddr, errInvalidAddress)
	}
	path := fmt.Sprintf(pathTxByAddr, inputAddr)

	var response []TransactionResponse
	resp, err := c.client.R().
		SetResult(&response).
		Get(path)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get transactions: %w", errHTTPError)
	}
	return &response, nil
}

func (c *RestyClient) GetTransactionByHash(txHash string) (*TransactionResponse, error) {
	if txHash == "" {
		return nil, fmt.Errorf("transaction hash cannot be empty")
	}
	path := fmt.Sprintf(pathTxByHash, txHash)

	response := &TransactionResponse{}
	resp, err := c.client.R().
		SetResult(response).
		Get(path)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get transaction: %w", errHTTPError)
	}

	return response, nil
}

func (c *RestyClient) GetBlockByHeight(height uint64) (*BlockResponse, error) {
	if height < 0 {
		return nil, fmt.Errorf("invalid block height")
	}

	path := fmt.Sprintf(pathBlockByHeight, fmt.Sprint(height))

	response := &BlockResponse{}
	resp, err := c.client.R().
		SetResult(response).
		Get(path)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get block: %w", errHTTPError)
	}

	return response, nil
}

func IsValidAddress(inputAddr string) bool {
	if len(inputAddr) == 0 {
		return false
	}
	// Prefix 0x
	if !strings.HasPrefix(inputAddr, "0x") {
		return false
	}
	// Trim Prefix 0x
	trimPrefix0xAddr := strings.TrimPrefix(inputAddr, "0x")

	// white space
	trimmedAddr := strings.TrimSpace(trimPrefix0xAddr)
	if len(trimmedAddr) == 0 {
		return false
	}

	isAllZeros := true
	for _, c := range trimPrefix0xAddr {
		if c != '0' {
			isAllZeros = false
			break
		}
	}
	if isAllZeros {
		return false
	}

	// check all hex str
	_, err := hex.DecodeString(trimPrefix0xAddr)
	return err == nil
}

func ValidateSubmitTransaction(req *SubmitTransactionRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	if !IsValidAddress(req.Sender) {
		return fmt.Errorf("invalid sender address: %s", req.Sender)
	}

	// require req
	if req.SequenceNumber == "" {
		return errors.New("sequence number is required")
	}
	if req.MaxGasAmount == "" {
		return errors.New("max gas amount is required")
	}
	if req.GasUnitPrice == "" {
		return errors.New("gas unit price is required")
	}
	if req.ExpirationTimestampSecs == "" {
		return errors.New("expiration timestamp is required")
	}

	return nil
}
