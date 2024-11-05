package aptos

import (
	"fmt"
	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateDifferentAccounts(t *testing.T) {
	//privKey.ToHex()
	//privKey.Bytes()
	//privKey.PubKey()
	//
	//pubKey.ToHex()
	//pubKey.Bytes()
	//pubKey.AuthKey()
	//
	//account.Address()

	t.Run("Test NewEd25519Account", func(t *testing.T) {
		// Private (hex): 0x482607f684af7b554fb0a5cff46f3e13f55c386a8c2078c85fdd23482cd46d0d
		// pubkey (hex): 0xe053ceb35bdeeb607c404a76ee6a123e2e9527c158a4d283133b9a8053e9599b
		// address: 0xd9a32936401d4c547938ac6af3def73e74d2eea4912e5fc00276d41c04c6f814
		// authKey: 0xd9a32936401d4c547938ac6af3def73e74d2eea4912e5fc00276d41c04c6f814
		// signMessage: 0xc7049e4670af30c404dbd90eaaf36482a06d3dd00f9c8913130862714860961cdfd8d4793d6b74ed2720be72511750989e2c03a67c93d1eaac6e658647b76504
		account, err := aptos.NewEd25519Account()
		if err != nil {
			panic(err)
		}

		privateKey := account.Signer.(*crypto.Ed25519PrivateKey)
		fmt.Printf(" (hex): %s\n", privateKey.ToHex())

		publicKey := account.PubKey()
		fmt.Printf("pubkey (hex): %s\n", publicKey.ToHex())

		address := account.Address
		fmt.Printf("address: %s\n", address.String())

		authKey := account.AuthKey()
		fmt.Printf("authKey: %s\n", authKey.ToHex())

		message := []byte("Hello Aptos!")
		signature, err := account.SignMessage(message)
		if err != nil {
			panic(err)
		}
		fmt.Printf("signMessage: %s\n", signature.ToHex())
	})

	t.Run("Test NewEd25519Account", func(t *testing.T) {
		// Private (hex): 0x437f72e66d2955e6c65d3e97aa6aa7345d63b4afda0c58e78971247bd073843a
		// pubkey (hex): 0xdda59a045b6dc46c378fa84218c0d116c97795c63d8e8b56868fe2492f3bd45b
		// address: 0x34d8e3074323789467ce1e5d2c538312985dcd3b8889f29ce23e08b0d798404d
		// authKey: 0x34d8e3074323789467ce1e5d2c538312985dcd3b8889f29ce23e08b0d798404d
		// signMessage: 0xd18492c1a4f906d88c6842d6b9d3c6ac94a1cce79b2128e53e37ce7dd3e60810d0b88f18c18915175a9b0e9836b9023705080053b40f02a36c9c7bc199cab00c
		account, err := aptos.NewEd25519Account()
		if err != nil {
			panic(err)
		}

		privateKey := account.Signer.(*crypto.Ed25519PrivateKey)
		fmt.Printf("Private (hex): %s\n", privateKey.ToHex())

		publicKey := account.PubKey()
		fmt.Printf("pubkey (hex): %s\n", publicKey.ToHex())

		address := account.Address
		fmt.Printf("address: %s\n", address.String())

		authKey := account.AuthKey()
		fmt.Printf("authKey: %s\n", authKey.ToHex())

		message := []byte("Hello Aptos!")
		signature, err := account.SignMessage(message)
		if err != nil {
			panic(err)
		}
		fmt.Printf("signMessage: %s\n", signature.ToHex())
	})
}

func TestValidateAptosAddress(t *testing.T) {
	t.Run("验证 Aptos 地址", func(t *testing.T) {
		// 1. 创建 AccountAddress 实例
		address := &aptos.AccountAddress{}

		// 2. 测试地址
		testAddress := "0x34d8e3074323789467ce1e5d2c538312985dcd3b8889f29ce23e08b0d798404d"

		// 3. 解析地址
		err := address.ParseStringRelaxed(testAddress)
		if err != nil {
			t.Fatalf("地址解析失败: %v", err)
		}

		// 4. 验证地址格式
		// 验证标准格式输出
		assert.Equal(t, testAddress, address.String(), "地址格式化输出应该匹配")

		// 5. 验证完整格式（补零）
		fullAddress := address.StringLong()
		assert.Equal(t, 66, len(fullAddress), "完整地址长度应为66（包含0x前缀）")

		// 6. 测试不同格式的地址
		// 测试无0x前缀
		addressWithoutPrefix := testAddress[2:]
		err = address.ParseStringRelaxed(addressWithoutPrefix)
		assert.NoError(t, err, "应该能解析无0x前缀的地址")

		// 9. 打印验证结果
		fmt.Printf("地址验证成功:\n")
		fmt.Printf("address 标准格式: %s\n", address.String())
		fmt.Printf("address 完整格式: %s\n", address.StringLong())
		fmt.Printf("address 完整格式 len: %d\n", len(address))
	})
}

func TestAptosTransfer(t *testing.T) {
	// 创建客户端
	client, err := aptos.NewClient(aptos.DevnetConfig)
	assert.NoError(t, err, "创建客户端失败")

	// 创建测试账户
	alice, err := aptos.NewEd25519Account()
	assert.NoError(t, err, "创建 alice 账户失败")

	bob, err := aptos.NewEd25519Account()
	assert.NoError(t, err, "创建 bob 账户失败")

	// 为 alice 账户充值
	const fundAmount = 100_000_000
	err = client.Fund(alice.Address, fundAmount)
	assert.NoError(t, err, "为 alice 充值失败")
	t.Logf("已为账户 %s 充值 %d APT", alice.Address, fundAmount)

	// 获取初始余额
	aliceBalance, err := client.AccountAPTBalance(alice.Address)
	assert.NoError(t, err, "获取 alice 余额失败")
	assert.Equal(t, uint64(fundAmount), aliceBalance, "alice 初始余额不正确")
	t.Logf("Alice 初始余额: %d APT", aliceBalance)

	// 构建转账交易
	const transferAmount = 1_000
	transferPayload, err := aptos.CoinTransferPayload(nil, bob.Address, transferAmount)
	assert.NoError(t, err, "构建转账 payload 失败")

	// 构建原始交易
	rawTxn, err := client.BuildTransaction(alice.AccountAddress(),
		aptos.TransactionPayload{Payload: transferPayload})
	assert.NoError(t, err, "构建交易失败")

	// 模拟交易
	simulationResult, err := client.SimulateTransaction(rawTxn, alice)
	assert.NoError(t, err, "模拟交易失败")
	assert.Equal(t, "Executed successfully", simulationResult[0].VmStatus, "交易模拟状态不正确")
	t.Logf("交易模拟成功，预计 gas 费用: %d", simulationResult[0].GasUsed)

	// 签名交易
	signedTxn, err := rawTxn.SignedTransaction(alice)
	assert.NoError(t, err, "签名交易失败")

	// 提交交易
	submitResult, err := client.SubmitTransaction(signedTxn)
	assert.NoError(t, err, "提交交易失败")

	// 等待交易完成
	txn, err := client.WaitForTransaction(submitResult.Hash)
	assert.NoError(t, err, "等待交易完成失败")
	assert.True(t, txn.Success, "交易执行失败")
	t.Logf("交易已确认，交易哈希: %s", submitResult.Hash)

	// 验证最终余额
	newAliceBalance, err := client.AccountAPTBalance(alice.Address)
	assert.NoError(t, err, "获取 alice 新余额失败")

	bobBalance, err := client.AccountAPTBalance(bob.Address)
	assert.NoError(t, err, "获取 bob 余额失败")

	// 验证转账结果（需要考虑 gas 费用）
	assert.Less(t, newAliceBalance, aliceBalance-transferAmount, "alice 余额扣除不正确")
	assert.Equal(t, uint64(transferAmount), bobBalance, "bob 余额增加不正确")
	t.Logf("转账后 Alice 余额: %d APT", newAliceBalance)
	t.Logf("转账后 Bob 余额: %d APT", bobBalance)
	t.Logf("实际 gas 费用: %d", aliceBalance-transferAmount-newAliceBalance)
}
