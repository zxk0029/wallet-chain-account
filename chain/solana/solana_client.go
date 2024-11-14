package solana

import (
	"context"
	"errors"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	gresty "github.com/go-resty/resty/v2"
)

type SolClient interface {
	GetHealth() (string, error)

	GetAccountInfo(inputAddr string) (*GetAccountInfoResponse, error)
	GetBalance(inputAddr string) (*GetBalanceResponse, error)

	//GetBlockHeight() (uint64, error)

	GetFeeForMessage(message string) (uint64, error)
	GetRecentPrioritizationFees() ([]PrioritizationFee, error)

	GetSlot(commitment CommitmentType) (uint64, error)
	GetBlocksWithLimit(startSlot uint64, limit uint64) ([]uint64, error)
	GetBlockBySlot(slot uint64, detailType TransactionDetailsType) (*BlockResult, error)

	GetTransaction(signature string) (*GetTransactionResponse, error)
	GetTransactionRange(signatures []string) ([]*GetTransactionResponse, error)
	GetTxForAddress(
		address string,
		commitment CommitmentType,
		limit uint64,
		beforeSignature string,
		untilSignature string,
	) ([]SignatureInfo, error)
}

type CommitmentType string

const (
	// Finalized Confirmed Processed
	// Finalized wait 32 slot
	Finalized CommitmentType = "finalized"
	// Confirmed wait 2-3 slot
	Confirmed CommitmentType = "confirmed"
	// Processed wait 0 slot
	Processed CommitmentType = "processed"
)

type TransactionDetailsType string

const (
	Full       TransactionDetailsType = "full"
	Accounts   TransactionDetailsType = "accounts"
	Signatures TransactionDetailsType = "signatures"
	None       TransactionDetailsType = "none"
)

const (
	defaultRequestTimeout   = 30 * time.Second
	defaultRetryCount       = 3
	defaultRetryWaitTime    = 10 * time.Second
	defaultRetryMaxWaitTime = 30 * time.Second
	defaultWithDebug        = false

	blockLimit = 50_0000
)

type solclient struct {
	grestyClient *gresty.Client
}

var (
	errHTTPError       = errors.New("aptos http error")
	errInvalidAddress  = errors.New("invalid address")
	errInvalidResponse = errors.New("invalid response")
)

func NewSolHttpClient(baseUrl string) (SolClient, error) {
	return NewSolHttpClientAll(baseUrl, defaultWithDebug)
}

func NewSolHttpClientAll(baseUrl string, withDebug bool) (SolClient, error) {
	grestyClient := gresty.New()
	grestyClient.SetBaseURL(baseUrl)
	grestyClient.SetTimeout(defaultRequestTimeout)
	grestyClient.SetRetryCount(defaultRetryCount)
	grestyClient.SetRetryWaitTime(defaultRetryWaitTime)
	grestyClient.SetRetryMaxWaitTime(defaultRetryMaxWaitTime)
	grestyClient.SetDebug(withDebug)

	// Retry Condition
	//grestyClient.AddRetryCondition(func(r *gresty.Response, err error) bool {
	//	return err != nil || r.StatusCode() >= 500
	//})

	grestyClient.OnBeforeRequest(func(c *gresty.Client, r *gresty.Request) error {
		log.Printf("Making request to %s (Attempt %d)", r.URL, r.Attempt)
		return nil
	})

	grestyClient.OnAfterResponse(func(c *gresty.Client, r *gresty.Response) error {
		statusCode := r.StatusCode()
		attempt := r.Request.Attempt
		method := r.Request.Method
		url := r.Request.URL
		log.Printf("Response received: Method=%s, URL=%s, Status=%d, Attempt=%d",
			method, url, statusCode, attempt)

		if statusCode >= 400 {
			if statusCode == 404 {
				return fmt.Errorf("%d resource not found %s %s: %w",
					statusCode, method, url, errHTTPError)
			}
			if statusCode >= 500 {
				return fmt.Errorf("%d server error %s %s: %w",
					statusCode, method, url, errHTTPError)
			}
			return fmt.Errorf("%d cannot %s %s: %w",
				statusCode, method, url, errHTTPError)
		}
		return nil
	})
	return &solclient{grestyClient: grestyClient}, nil
}

func (s *solclient) GetHealth() (string, error) {
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getHealth",
		"params":  []interface{}{},
	}

	resp := struct {
		Jsonrpc string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Result  string `json:"result"`
		Error   struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error,omitempty"`
	}{}

	httpResp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(&resp).
		Post("/")
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}

	if httpResp.IsError() {
		return "", fmt.Errorf("failed to get health: %w", errHTTPError)
	}

	if resp.Error.Message != "" {
		return "", fmt.Errorf("rpc error: %s (code: %d)",
			resp.Error.Message, resp.Error.Code)
	}

	return resp.Result, nil
}

func (s *solclient) GetAccountInfo(inputAddr string) (*GetAccountInfoResponse, error) {
	dealAddr := strings.TrimSpace(inputAddr)

	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getAccountInfo",
		"params": []interface{}{
			dealAddr,
			map[string]string{
				"encoding": "base64",
			},
		},
	}
	account := &GetAccountInfoResponse{}
	resp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(account).
		Post("/")
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get account info: %w", errHTTPError)
	}
	return account, nil
}

func (s *solclient) GetBalance(inputAddr string) (*GetBalanceResponse, error) {
	dealAddr := strings.TrimSpace(inputAddr)

	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getBalance",
		"params": []interface{}{
			dealAddr,
		},
	}

	balance := &GetBalanceResponse{}
	resp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(balance).
		Post("/")
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get balance: %w", errHTTPError)
	}
	return balance, nil
}

// GetBlockHeight get latest height
func (s *solclient) GetBlockHeight() (uint64, error) {
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getBlockHeight",
	}

	resp := &BlockHeightResponse{}
	httpResp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(resp).
		Post("/")
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}

	if httpResp.IsError() {
		return 0, fmt.Errorf("failed to get block height: %w", errHTTPError)
	}

	return resp.Result, nil
}

// GetSlot get latest slot
func (s *solclient) GetSlot(commitment CommitmentType) (uint64, error) {

	config := GetSlotRequest{
		Commitment: commitment,
	}

	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getSlot",
		"params":  []interface{}{config},
	}

	resp := &GetSlotResponse{}
	httpResp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(resp).
		Post("/")
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}

	if httpResp.IsError() {
		return 0, fmt.Errorf("failed to get slot: %w", errHTTPError)
	}

	return resp.Result, nil
}

// GetBlocksWithLimit returns a list of confirmed blocks starting at the given slot with limit
func (s *solclient) GetBlocksWithLimit(startSlot uint64, limit uint64) ([]uint64, error) {
	if limit > blockLimit {
		return nil, fmt.Errorf("limit must not exceed %d blocks", blockLimit)
	}

	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getBlocksWithLimit",
		"params":  []uint64{startSlot, limit},
	}

	resp := &GetBlocksWithLimitResponse{}
	httpResp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(resp).
		Post("/")

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if httpResp.IsError() {
		return nil, fmt.Errorf("failed to get blocks with limit: %w", errHTTPError)
	}

	return resp.Result, nil
}

func (s *solclient) GetBlockBySlot(slot uint64, detailType TransactionDetailsType) (*BlockResult, error) {
	//s.grestyClient.SetTimeout(120 * time.Second)
	//defer s.grestyClient.SetTimeout(defaultRequestTimeout)

	config := GetBlockRequest{
		Commitment:                     Finalized,
		Encoding:                       "json",
		MaxSupportedTransactionVersion: 0,
		TransactionDetails:             string(detailType),
		Rewards:                        false,
	}

	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getBlock",
		"params":  []interface{}{slot, config},
	}

	resp := &GetBlockResponse{}
	httpResp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(resp).
		Post("/")
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if httpResp.IsError() {
		return nil, fmt.Errorf("failed to get block: %w", errHTTPError)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("RPC error: (code: %d) %s", resp.Error.Code, resp.Error.Message)
	}

	return &resp.Result, nil
}

func (s *solclient) GetTransaction(signature string) (*GetTransactionResponse, error) {
	config := map[string]interface{}{
		"encoding":                       "json",
		"commitment":                     Finalized,
		"maxSupportedTransactionVersion": 0,
	}

	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getTransaction",
		"params":  []interface{}{signature, config},
	}

	resp := &GetTransactionResponse{}
	httpResp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(resp).
		Post("/")

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if httpResp.IsError() {
		return nil, fmt.Errorf("failed to get transaction: %w", errHTTPError)
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("RPC error: (code: %d) %s", resp.Error.Code, resp.Error.Message)
	}

	return resp, nil
}

func (s *solclient) GetTransactionRange(signatures []string) ([]*GetTransactionResponse, error) {
	if len(signatures) == 0 {
		return nil, fmt.Errorf("empty signatures")
	}

	if len(signatures) == 1 {
		tx, err := s.GetTransaction(signatures[0])
		if err != nil {
			return nil, err
		}
		return []*GetTransactionResponse{tx}, nil
	}

	count := len(signatures)
	transactions := make([]*GetTransactionResponse, count)
	var wg sync.WaitGroup

	const groupSize = 20
	numGroups := (count-1)/groupSize + 1
	errChan := make(chan error, numGroups)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	rateLimiter := time.NewTicker(100 * time.Millisecond)
	defer rateLimiter.Stop()

	for i := 0; i < count; i += groupSize {
		start := i
		end := i + groupSize - 1
		if end >= count {
			end = count - 1
		}
		wg.Add(1)

		go func(start, end int) {
			defer wg.Done()

			for j := start; j <= end; j++ {
				select {
				case <-ctx.Done():
					errChan <- ctx.Err()
					return
				case <-rateLimiter.C:
					tx, err := s.GetTransaction(signatures[j])
					if err != nil {
						if strings.Contains(err.Error(), "Transaction not found") {
							transactions[j] = nil
							continue
						}
						errChan <- fmt.Errorf("failed to get transaction %s: %w", signatures[j], err)
						return
					}
					transactions[j] = tx
				}
			}
		}(start, end)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	validTransactions := make([]*GetTransactionResponse, 0)
	for _, tx := range transactions {
		if tx != nil {
			validTransactions = append(validTransactions, tx)
		}
	}

	return validTransactions, nil
}

func (s *solclient) GetFeeForMessage(message string) (uint64, error) {
	config := GetFeeForMessageRequest{
		Commitment: string(Finalized),
	}

	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getFeeForMessage",
		"params":  []interface{}{message, config},
	}

	resp := &GetFeeForMessageResponse{}
	httpResp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(resp).
		Post("/")
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}

	if httpResp.IsError() {
		return 0, fmt.Errorf("failed to get fee for message: %w", errHTTPError)
	}

	if resp.Result.Value == nil {
		return 0, fmt.Errorf("invalid message or unable to estimate fee")
	}

	return *resp.Result.Value, nil
}

func (s *solclient) GetRecentPrioritizationFees() ([]PrioritizationFee, error) {
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getRecentPrioritizationFees",
		"params":  []interface{}{},
	}

	resp := &getRecentPrioritizationFeesResponse{}
	httpResp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(resp).
		Post("/")
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if httpResp.IsError() {
		return nil, fmt.Errorf("failed to get prioritization fees: %w", errHTTPError)
	}

	if len(resp.Result) == 0 {
		return []PrioritizationFee{}, nil
	}

	return resp.Result, nil
}

func GetSuggestedPriorityFee(fees []PrioritizationFee) uint64 {
	if len(fees) == 0 {
		return 0
	}

	priorityFees := make([]uint64, len(fees))
	for i, fee := range fees {
		priorityFees[i] = fee.PrioritizationFee
	}

	sort.Slice(priorityFees, func(i, j int) bool {
		return priorityFees[i] < priorityFees[j]
	})

	index := int(float64(len(priorityFees)) * 0.75)
	return priorityFees[index]
}

func (s *solclient) GetTxForAddress(
	address string,
	commitment CommitmentType,
	limit uint64,
	beforeSignature string,
	untilSignature string,
) ([]SignatureInfo, error) {
	if address == "" {
		return nil, fmt.Errorf("empty address")
	}
	config := &GetSignaturesRequest{
		Commitment: string(commitment),
		Limit:      limit,
		Before:     beforeSignature,
		Until:      untilSignature,
	}

	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getSignaturesForAddress",
		"params":  []interface{}{address, config},
	}

	resp := &GetSignaturesResponse{}
	httpResp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(resp).
		Post("/")
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if httpResp.IsError() {
		return nil, fmt.Errorf("failed to get signatures: %w", errHTTPError)
	}

	return resp.Result, nil
}

func (s *solclient) SimulateTransaction(
	signedTx string,
	config *SimulateRequest,
) (*SimulateTransactionResponse, error) {
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "simulateTransaction",
		"params": []interface{}{
			signedTx,
			config,
		},
	}

	resp := &SimulateTransactionResponse{}
	httpResp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(resp).
		Post("/")
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if httpResp.IsError() {
		return nil, fmt.Errorf("failed to simulate transaction: %w", errHTTPError)
	}

	return resp, nil
}

func (s *solclient) SendTransaction(
	signedTx string,
	config *SendTransactionRequest,
) (string, error) {
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "sendTransaction",
		"params": []interface{}{
			signedTx,
			config,
		},
	}

	resp := struct {
		Jsonrpc string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Result  string `json:"result"`
		Error   struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error,omitempty"`
	}{}

	httpResp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(&resp).
		Post("/")
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}

	if httpResp.IsError() {
		return "", fmt.Errorf("failed to send transaction: %w", errHTTPError)
	}

	if resp.Error.Message != "" {
		return "", fmt.Errorf("rpc error: %s (code: %d)",
			resp.Error.Message, resp.Error.Code)
	}

	return resp.Result, nil
}

func GetNonceAccount(client *rpc.Client, nonceAccountPubkey solana.PublicKey) error {
	nonceAccountInfo, err := client.GetAccountInfo(ctx, nonceAccountPubkey)
	if err != nil {
		return fmt.Errorf("get nonce account error: %w", err)
	}

	// 2. 解析 nonce 数据
	nonceData, err := system.NonceAccount(nonceAccountInfo.Value.Data.GetBinary())
	if err != nil {
		return fmt.Errorf("parse nonce data error: %w", err)
	}
}
