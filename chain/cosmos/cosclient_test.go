package cosmos

import (
	"context"
	authv1beta1 "cosmossdk.io/api/cosmos/auth/v1beta1"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

const (
	defaultRpcAddress = "https://cosmos-rpc.publicnode.com:443"
)

func TestClient_GetAccount(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()
	c, err := DialCosmosClient(ctx, defaultRpcAddress)
	assert.NoError(t, err)

	response, err := c.GetAccount("cosmos19thxsunl9lzywglsndth5a278wtavawzzpv44q")
	assert.NoError(t, err)

	authAccount := new(authv1beta1.BaseAccount)
	err = ptypes.UnmarshalAny(response.Account, authAccount)
	assert.NoError(t, err)
	fmt.Printf("sequence: %s, account number: %s, address: %s \n",
		strconv.FormatUint(authAccount.GetSequence(), 10),
		strconv.FormatUint(authAccount.GetAccountNumber(), 10),
		authAccount.GetAddress())
}

func TestClient_GetBalance(t *testing.T) {
	c, err := DialCosmosClient(context.Background(), defaultRpcAddress)
	assert.NoError(t, err)

	ret, err := c.GetBalance("uatom", "cosmos19thxsunl9lzywglsndth5a278wtavawzzpv44q")
	assert.NoError(t, err)
	fmt.Println("here", ret)
}

func TestClient_GetTxByHash(t *testing.T) {
	c, err := DialCosmosClient(context.Background(), defaultRpcAddress)
	assert.NoError(t, err)

	ret, err := c.GetTxByHash("85C84677F466D71C0BB6E744439C3040ABB35B8F2B838CC7B73CD1BFF33D0B88")
	assert.NoError(t, err)
	fmt.Println("here", ret)
}

func TestClient_GetBlock(t *testing.T) {
	c, err := DialCosmosClient(context.Background(), defaultRpcAddress)
	assert.NoError(t, err)

	height := int64(22879895)
	ret, err := c.GetBlock(&height)
	assert.NoError(t, err)
	fmt.Println("here", ret)
}

func TestClient_GetTxByEvent(t *testing.T) {
	c, err := DialCosmosClient(context.Background(), defaultRpcAddress)
	assert.NoError(t, err)

	event := []string{"send"}
	ret, err := c.GetTxByEvent(event, 0, 10)
	assert.NoError(t, err)
	fmt.Println("here", ret)
}
