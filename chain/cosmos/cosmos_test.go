package cosmos

import (
	"flag"
	"fmt"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	defaultCosmosRpcAddress = "https://cosmos-rpc.publicnode.com:443"
)

func TestCosmos_GetBlockByNumber(t *testing.T) {
	var f = flag.String("c", "../../config.yml", "config path")
	flag.Parse()
	conf, err := config.New(*f)
	assert.NoError(t, err)

	chainAdaptor, err := NewChainAdaptor(conf)
	assert.NoError(t, err)
	request := &account.BlockNumberRequest{
		Height: int64(22879895),
	}
	response, err := chainAdaptor.GetBlockByNumber(request)
	assert.NoError(t, err)
	fmt.Println("response", response)
}
