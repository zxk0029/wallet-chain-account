package main // 声明这是可执行程序的入口包

import (
	"flag" // 命令行参数解析
	"net"  // 网络相关操作

	"google.golang.org/grpc"            // gRPC框架
	"google.golang.org/grpc/reflection" // gRPC反射服务（用于测试调试）

	"github.com/ethereum/go-ethereum/log" // 以太坊的日志库

	"github.com/dapplink-labs/wallet-chain-account/chaindispatcher"     // 项目自定义的链调度器
	"github.com/dapplink-labs/wallet-chain-account/config"              // 配置管理
	wallet2 "github.com/dapplink-labs/wallet-chain-account/rpc/account" // RPC服务实现
)

// 程序主入口
func main() {
	// 解析命令行参数（-c 指定配置文件路径）
	var f = flag.String("c", "config.yml", "config path")
	flag.Parse()

	// 加载配置文件
	conf, err := config.New(*f)
	if err != nil {
		panic(err) // 如果加载失败直接终止程序
	}

	// 创建链调度器实例（核心业务逻辑）
	dispatcher, err := chaindispatcher.New(conf)
	if err != nil {
		log.Error("Setup dispatcher failed", "err", err)
		panic(err)
	}

	// 创建gRPC服务器，添加调度器的拦截器
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(dispatcher.Interceptor))
	defer grpcServer.GracefulStop() // 程序退出时优雅关闭服务

	// 注册钱包账户服务到gRPC服务器
	wallet2.RegisterWalletAccountServiceServer(grpcServer, dispatcher)

	// 创建TCP监听端口（端口号来自配置文件）
	listen, err := net.Listen("tcp", ":"+conf.Server.Port)
	if err != nil {
		log.Error("net listen failed", "err", err)
		panic(err)
	}
	reflection.Register(grpcServer) // 启用gRPC反射服务（方便客户端调试）

	// 输出启动成功日志
	log.Info("dapplink wallet rpc services start success", "port", conf.Server.Port)

	// 启动gRPC服务器（阻塞式运行）
	if err := grpcServer.Serve(listen); err != nil {
		log.Error("grpc server serve failed", "err", err)
		panic(err)
	}
}
