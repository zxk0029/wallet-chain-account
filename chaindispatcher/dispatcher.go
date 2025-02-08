package chaindispatcher

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/chain/aptos"
	"github.com/dapplink-labs/wallet-chain-account/chain/cosmos"
	"github.com/dapplink-labs/wallet-chain-account/chain/ethereum"
	"github.com/dapplink-labs/wallet-chain-account/chain/polygon"
	"github.com/dapplink-labs/wallet-chain-account/chain/solana"
	"github.com/dapplink-labs/wallet-chain-account/chain/sui"
	"github.com/dapplink-labs/wallet-chain-account/chain/ton"
	"github.com/dapplink-labs/wallet-chain-account/chain/tron"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	"github.com/dapplink-labs/wallet-chain-account/rpc/common"
)

// 通用请求接口（所有gRPC请求必须实现）
type CommonRequest interface {
	GetChain() string // 获取请求的目标区块链
}

type ChainType = string

// 链调度器核心结构
type ChainDispatcher struct {
	registry map[string]chain.IChainAdaptor // 区块链适配器注册表（key: 链名称）
}

// 初始化函数（项目启动时调用）
func New(conf *config.Config) (*ChainDispatcher, error) {
	dispatcher := ChainDispatcher{
		registry: make(map[ChainType]chain.IChainAdaptor),
	}

	// 链适配器工厂函数映射表（key: 链名称）
	chainAdaptorFactoryMap := map[string]func(conf *config.Config) (chain.IChainAdaptor, error){
		ethereum.ChainName: ethereum.NewChainAdaptor, // 以太坊适配器
		cosmos.ChainName:   cosmos.NewChainAdaptor,   // Cosmos适配器
		solana.ChainName:   solana.NewChainAdaptor,   // Solana适配器
		tron.ChainName:     tron.NewChainAdaptor,     // Tron适配器
		aptos.ChainName:    aptos.NewChainAdaptor,    // Aptos适配器
		sui.ChainName:      sui.NewSuiAdaptor,        // Sui适配器
		ton.ChainName:      ton.NewChainAdaptor,      // Ton适配器
		polygon.ChainName:  polygon.NewChainAdaptor,  // Polygon适配器
	}

	supportedChains := []string{
		ethereum.ChainName,
		cosmos.ChainName,
		solana.ChainName,
		tron.ChainName,
		sui.ChainName,
		ton.ChainName,
		aptos.ChainName,
		polygon.ChainName,
	}

	// 遍历配置文件中的链名称，初始化链适配器
	for _, c := range conf.Chains {
		if factory, ok := chainAdaptorFactoryMap[c]; ok {
			adaptor, err := factory(conf) // 使用工厂函数创建适配器
			if err != nil {
				log.Crit("failed to setup chain", "chain", c, "error", err)
			}
			dispatcher.registry[c] = adaptor // 将适配器注册到调度器中
		} else {
			log.Error("unsupported chain", "chain", c, "supportedChains", supportedChains)
		}
	}
	return &dispatcher, nil
}

// gRPC拦截器（所有请求都会先经过这里）
func (d *ChainDispatcher) Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error("panic error", "msg", e)
			log.Debug(string(debug.Stack()))
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()

	// 解析请求方法名（如 /service/method）
	pos := strings.LastIndex(info.FullMethod, "/")
	method := info.FullMethod[pos+1:]

	chainName := req.(CommonRequest).GetChain() // 获取请求的链名称
	log.Info(method, "chain", chainName, "req", req)

	resp, err = handler(ctx, req) // 调用实际的请求处理函数，继续处理请求
	log.Debug("Finish handling", "resp", resp, "err", err)
	return
}

// 预处理请求（检查链是否存在）
func (d *ChainDispatcher) preHandler(req interface{}) error {
	chainName := req.(CommonRequest).GetChain() // 获取请求的链名称
	if _, ok := d.registry[chainName]; !ok {    // 检查链是否注册
		return fmt.Errorf("%s (chain: %s)", config.UnsupportedOperation, chainName)
	}
	return nil
}

func (d *ChainDispatcher) GetSupportChains(ctx context.Context, request *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	if err := d.preHandler(request); err != nil {
		return &account.SupportChainsResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  fmt.Sprintf("GetSupportChains Fail: %v", err),
		}, nil
	}
	return d.registry[request.Chain].GetSupportChains(request)
}

func (d *ChainDispatcher) ConvertAddress(ctx context.Context, request *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	if err := d.preHandler(request); err != nil {
		return &account.ConvertAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  fmt.Sprintf("ConvertAddress Fail: %v", err),
		}, nil
	}
	return d.registry[request.Chain].ConvertAddress(request)
}

func (d *ChainDispatcher) ValidAddress(ctx context.Context, request *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	if err := d.preHandler(request); err != nil {
		return &account.ValidAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  fmt.Sprintf("ValidAddress Fail: %v", err),
		}, nil
	}
	return d.registry[request.Chain].ValidAddress(request)
}

func (d *ChainDispatcher) GetBlockByNumber(ctx context.Context, request *account.BlockNumberRequest) (*account.BlockResponse, error) {
	if err := d.preHandler(request); err != nil {
		return &account.BlockResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  fmt.Sprintf("GetBlockByNumber Fail: %v", err),
		}, nil
	}
	return d.registry[request.Chain].GetBlockByNumber(request)
}

func (d *ChainDispatcher) GetBlockByHash(ctx context.Context, request *account.BlockHashRequest) (*account.BlockResponse, error) {
	if err := d.preHandler(request); err != nil {
		return &account.BlockResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  fmt.Sprintf("GetBlockByHash Fail: %v", err),
		}, nil
	}
	return d.registry[request.Chain].GetBlockByHash(request)
}

func (d *ChainDispatcher) GetBlockHeaderByHash(ctx context.Context, request *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	if err := d.preHandler(request); err != nil {
		return &account.BlockHeaderResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  fmt.Sprintf("GetBlockHeaderByHash Fail: %v", err),
		}, nil
	}
	return d.registry[request.Chain].GetBlockHeaderByHash(request)
}

func (d *ChainDispatcher) GetBlockHeaderByNumber(ctx context.Context, request *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	if err := d.preHandler(request); err != nil {
		return &account.BlockHeaderResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  fmt.Sprintf("GetBlockHeaderByNumber Fail: %v", err),
		}, nil
	}
	return d.registry[request.Chain].GetBlockHeaderByNumber(request)
}

func (d *ChainDispatcher) GetBlockHeaderByRange(ctx context.Context, request *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	if err := d.preHandler(request); err != nil {
		return &account.BlockByRangeResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  fmt.Sprintf("GetBlockByRange Fail: %v", err),
		}, nil
	}
	return d.registry[request.Chain].GetBlockByRange(request)
}

func (d *ChainDispatcher) GetAccount(ctx context.Context, request *account.AccountRequest) (*account.AccountResponse, error) {
	if err := d.preHandler(request); err != nil {
		return &account.AccountResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  fmt.Sprintf("GetAccount Fail: %v", err),
		}, nil
	}
	return d.registry[request.Chain].GetAccount(request)
}

func (d *ChainDispatcher) GetFee(ctx context.Context, request *account.FeeRequest) (*account.FeeResponse, error) {
	if err := d.preHandler(request); err != nil {
		return &account.FeeResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  fmt.Sprintf("GetFee Fail: %v", err),
		}, nil
	}
	return d.registry[request.Chain].GetFee(request)
}

func (d *ChainDispatcher) SendTx(ctx context.Context, request *account.SendTxRequest) (*account.SendTxResponse, error) {
	if err := d.preHandler(request); err != nil {
		return &account.SendTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  fmt.Sprintf("SendTx Fail: %v", err),
		}, nil
	}
	return d.registry[request.Chain].SendTx(request)
}

func (d *ChainDispatcher) GetTxByAddress(ctx context.Context, request *account.TxAddressRequest) (*account.TxAddressResponse, error) {
	if err := d.preHandler(request); err != nil {
		return &account.TxAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  fmt.Sprintf("GetTxByAddress Fail: %v", err),
		}, nil
	}
	return d.registry[request.Chain].GetTxByAddress(request)
}

func (d *ChainDispatcher) GetTxByHash(ctx context.Context, request *account.TxHashRequest) (*account.TxHashResponse, error) {
	if err := d.preHandler(request); err != nil {
		return &account.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  fmt.Sprintf("GetTxByHash Fail: %v", err),
		}, nil
	}
	return d.registry[request.Chain].GetTxByHash(request)
}

func (d *ChainDispatcher) GetBlockByRange(ctx context.Context, request *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	if err := d.preHandler(request); err != nil {
		return &account.BlockByRangeResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  fmt.Sprintf("GetBlockByRange Fail: %v", err),
		}, nil
	}
	return d.registry[request.Chain].GetBlockByRange(request)
}

func (d *ChainDispatcher) CreateUnSignTransaction(ctx context.Context, request *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	if err := d.preHandler(request); err != nil {
		return &account.UnSignTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  fmt.Sprintf("CreateUnSignTransaction Fail: %v", err),
		}, nil
	}
	return d.registry[request.Chain].CreateUnSignTransaction(request)
}

func (d *ChainDispatcher) BuildSignedTransaction(ctx context.Context, request *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	if err := d.preHandler(request); err != nil {
		return &account.SignedTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  fmt.Sprintf("BuildSignedTransaction Fail: %v", err),
		}, nil
	}
	return d.registry[request.Chain].BuildSignedTransaction(request)
}

func (d *ChainDispatcher) DecodeTransaction(ctx context.Context, request *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	if err := d.preHandler(request); err != nil {
		return &account.DecodeTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  fmt.Sprintf("DecodeTransaction Fail: %v", err),
		}, nil
	}
	return d.registry[request.Chain].DecodeTransaction(request)
}

func (d *ChainDispatcher) VerifySignedTransaction(ctx context.Context, request *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	if err := d.preHandler(request); err != nil {
		return &account.VerifyTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  fmt.Sprintf("VerifySignedTransaction Fail: %v", err),
		}, nil
	}
	return d.registry[request.Chain].VerifySignedTransaction(request)
}

func (d *ChainDispatcher) GetExtraData(ctx context.Context, request *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	if err := d.preHandler(request); err != nil {
		return &account.ExtraDataResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  fmt.Sprintf("GetExtraData Fail: %v", err),
		}, nil
	}
	return d.registry[request.Chain].GetExtraData(request)
}
