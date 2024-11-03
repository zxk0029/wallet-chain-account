package main

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/tendermint/tendermint/rpc/client/http"
)

func main() {
	ctx := client.Context{}
	nodeURI := "https://cosmos-rpc.publicnode.com:443" // 默认的 Tendermint RPC 端口
	client, _ := client.NewClientFromNode(nodeURI)
	ctx = ctx.WithClient(client)

	client.Status(context.Background())
	client.Header()

	// 初始化查询客户端
	queryClient := types.NewQueryClient(ctx)

	// 获取账户余额
	address := "cosmos19thxsunl9lzywglsndth5a278wtavawzzpv44q"
	getAccountBalance(queryClient, address)

	// 初始化 Cosmos SDK 客户端
	//ctx := client.NewClientContextFromViper(viper.GetViper()).
	//	WithChainID("your_chain_id").
	//	WithKeyring(keyring.NewInMemory()).
	//	WithClient(client.NewClient(nil, nil))

	//// 创建 Tendermint RPC 客户端
	//nodeURI := "https://cosmos-rpc.publicnode.com:443" // 默认的 Tendermint RPC 端口
	//tmClient, _ := http.New(nodeURI, "/websocket")
	//client.NewClientFromNode(nodeURI)
	//client.NewClientContextFromViper
	//
	//cc := tmClient.(client.CometRPC)
	//client.Context{}.WithClient()
	//
	//// 获取最新的区块
	//status, err := tmClient.Status(context.Background())
	//if err != nil {
	//	fmt.Printf("Failed to get status: %v\n", err)
	//	return
	//}
	//
	//// 获取区块高度
	//blockHeight := status.SyncInfo.LatestBlockHeight
	//
	//fmt.Printf("Current block height: %d\n", blockHeight)
	//
	//result, _ := tmClient.BlockByHash(context.Background(), []byte("85C84677F466D71C0BB6E744439C3040ABB35B8F2B838CC7B73CD1BFF33D0B88"))
	//
	//fmt.Printf("Current block height: %s\n", result)
	//
	//result, _ = tmClient.Block(context.Background(), &blockHeight)
	//fmt.Printf("Current block height: %s\n", result)

	//var f = flag.String("c", "config.yml", "config path")
	//flag.Parse()
	//conf, err := config.New(*f)
	//if err != nil {
	//	panic(err)
	//}
	//dispatcher, err := chaindispatcher.New(conf)
	//if err != nil {
	//	log.Error("Setup dispatcher failed", "err", err)
	//	panic(err)
	//}
	//
	//grpcServer := grpc.NewServer(grpc.UnaryInterceptor(dispatcher.Interceptor))
	//defer grpcServer.GracefulStop()
	//
	//wallet2.RegisterWalletAccountServiceServer(grpcServer, dispatcher)
	//
	//listen, err := net.Listen("tcp", ":"+conf.Server.Port)
	//if err != nil {
	//	log.Error("net listen failed", "err", err)
	//	panic(err)
	//}
	//reflection.Register(grpcServer)
	//
	//log.Info("dapplink wallet rpc services start success", "port", conf.Server.Port)
	//
	//if err := grpcServer.Serve(listen); err != nil {
	//	log.Error("grpc server serve failed", "err", err)
	//	panic(err)
	//}
}

func getAccountBalance(queryClient types.QueryClient, address string) {
	// 创建查询请求
	req := &types.QueryAllBalancesRequest{
		Address: address,
		Pagination: &query.PageRequest{
			Limit: 100, // 设置分页限制
		},
	}

	// 发送查询请求
	res, err := queryClient.AllBalances(context.Background(), req)
	if err != nil {
		panic(err)
	}

	// 打印账户余额
	fmt.Printf("Account: %s\n", address)
	for _, coin := range res.Balances {
		fmt.Printf("Denom: %s, Amount: %s\n", coin.Denom, coin.Amount)
	}
}

func getBlock(tmRPC *http.HTTP, height int64) {
	// 获取区块信息
	block, err := tmRPC.Block(context.Background(), &height)
	if err != nil {
		panic(err)
	}

	// 打印区块信息
	fmt.Printf("Block Height: %d\n", block.Block.Height)
	fmt.Printf("Block Hash: %s\n", block.Block.Hash())
	fmt.Printf("Block Time: %s\n", block.Block.Time)
}

func getTransaction(tmRPC *http.HTTP, txHash string) {
	// 获取交易信息
	txResult, err := tmRPC.Tx(context.Background(), []byte(txHash), true)
	if err != nil {
		panic(err)
	}

	// 打印交易信息
	fmt.Printf("Transaction Hash: %s\n", txResult.Hash)
	fmt.Printf("Height: %d\n", txResult.Height)
	fmt.Printf("Index: %d\n", txResult.Index)
	fmt.Printf("TxResult: %+v\n", txResult.TxResult)
	fmt.Printf("Tx: %s\n", txResult.Tx)
}
