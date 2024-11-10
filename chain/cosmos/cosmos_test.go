package cosmos

import (
	"flag"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
)

const (
	defaultConfigPath = "../../config.yml"
)

func getChainAdaptor() (chain.IChainAdaptor, error) {
	var f = flag.String("c", defaultConfigPath, "config path")
	flag.Parse()
	conf, _ := config.New(*f)
	return NewChainAdaptor(conf)
}

// success
func TestCosmos_ConvertAddress(t *testing.T) {
	chainAdaptor, err := getChainAdaptor()
	assert.NoError(t, err)
	request := &account.ConvertAddressRequest{
		PublicKey: "036e8f01c6e68d9c5c66ab172ff0898234c7a889997802437bce6e7d1f89161fc1",
	}
	response, err := chainAdaptor.ConvertAddress(request)
	assert.NoError(t, err)
	fmt.Println("response", response)
}

// success
func TestCosmos_ValidAddress(t *testing.T) {
	chainAdaptor, err := getChainAdaptor()
	assert.NoError(t, err)
	request := &account.ValidAddressRequest{
		Address: "cosmos19thxsunl9lzywglsndth5a278wtavawzzpv44q",
	}
	response, err := chainAdaptor.ValidAddress(request)
	assert.NoError(t, err)
	fmt.Println("response", response)
}

// success
func TestCosmos_GetAccount(t *testing.T) {
	chainAdaptor, err := getChainAdaptor()
	assert.NoError(t, err)
	request := &account.AccountRequest{
		Address: "cosmos19thxsunl9lzywglsndth5a278wtavawzzpv44q",
	}
	response, err := chainAdaptor.GetAccount(request)
	assert.NoError(t, err)
	fmt.Println("response", response)
}

// success
func TestCosmos_GetBlockByNumber(t *testing.T) {
	chainAdaptor, err := getChainAdaptor()
	assert.NoError(t, err)
	request := &account.BlockNumberRequest{
		Height: int64(22879895),
	}
	response, err := chainAdaptor.GetBlockByNumber(request)
	assert.NoError(t, err)
	fmt.Println("response", response)
}

func TestCosmos_GetBlockByHash(t *testing.T) {
	chainAdaptor, err := getChainAdaptor()
	assert.NoError(t, err)

	req := &account.BlockHashRequest{
		Hash: "35290F91317064307B5B1A9A44EEFB1CF3F66F68EAAD4539BCD6A5BA13866E90",
	}
	response, err := chainAdaptor.GetBlockByHash(req)
	assert.NoError(t, err)
	fmt.Println("response", response)
}

func TestCosmos_GetBlockHeaderByHash(t *testing.T) {
	chainAdaptor, err := getChainAdaptor()
	assert.NoError(t, err)

	req := &account.BlockHeaderHashRequest{
		Hash: "35290F91317064307B5B1A9A44EEFB1CF3F66F68EAAD4539BCD6A5BA13866E90",
	}
	response, err := chainAdaptor.GetBlockHeaderByHash(req)
	assert.NoError(t, err)
	fmt.Println("response", response)
}

// success
func TestCosmos_GetTxByHash(t *testing.T) {
	chainAdaptor, err := getChainAdaptor()
	assert.NoError(t, err)
	request := &account.TxHashRequest{
		Hash: "85C84677F466D71C0BB6E744439C3040ABB35B8F2B838CC7B73CD1BFF33D0B88", //"35290F91317064307B5B1A9A44EEFB1CF3F66F68EAAD4539BCD6A5BA13866E90",
	}
	response, err := chainAdaptor.GetTxByHash(request)
	assert.NoError(t, err)
	fmt.Println("response", response)
}

func TestCosmos_GetTxByAddress(t *testing.T) {
	chainAdaptor, err := getChainAdaptor()
	assert.NoError(t, err)

	request := &account.TxAddressRequest{
		Chain:    "Cosmos",
		Address:  "cosmos1nvcgd368m4pm5mm3ppzawhsq6grra4ejnppplx",
		Pagesize: 1,
		Page:     1,
	}

	response, err := chainAdaptor.GetTxByAddress(request)
	assert.NoError(t, err)
	fmt.Println("response", response)
}

func TestCosmos_GetBlockByRange(t *testing.T) {
	chainAdaptor, err := getChainAdaptor()
	assert.NoError(t, err)

	request := &account.BlockByRangeRequest{
		Start: "22879895",
		End:   "22879896",
	}
	response, err := chainAdaptor.GetBlockByRange(request)
	assert.NoError(t, err)
	fmt.Println("response", response)
}
