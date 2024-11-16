package cosmos

import (
	"encoding/base64"
	"encoding/json"
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
		PublicKey: "032d535553c70dfbb9c13f32cb6d1002a4b421beff39009670e29a7e51fb88ec3f",
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

func TestCosmos_CreateUnSignTransaction(t *testing.T) {
	txStruct := &TxStructure{
		ChainId:         "cosmoshub-4",
		FromAddress:     "cosmos1qgas8xpptnp09lyl32kfp60hldges6guu28qmk",
		ToAddress:       "cosmos1l6vul20q74gw56fped8srkjq2x8d9m305gnxr2",
		ContractAddress: "",
		Amount:          10000,
		GasLimit:        137674,
		FeeAmount:       10000,
		Sequence:        0,
		AccountNumber:   3014650,
		Decimal:         6,
		Memo:            "10086",
		PubKey:          "03f16c9160c81b806a04da3c27d9200fe684aad79b9fcdcccfac8aa60ad7f0a56a",
	}

	txBytes, err := json.Marshal(txStruct)
	assert.NoError(t, err)

	base64Tx := base64.StdEncoding.EncodeToString(txBytes)

	chainAdaptor, err := getChainAdaptor()
	assert.NoError(t, err)

	request := &account.UnSignTransactionRequest{
		Chain:    ChainName,
		Network:  NetWork,
		Base64Tx: base64Tx,
	}

	response, err := chainAdaptor.CreateUnSignTransaction(request)
	assert.NoError(t, err)
	fmt.Println("response", response)
}

func TestCosmos_BuildSignedTransaction(t *testing.T) {
	txStruct := &TxStructure{
		ChainId:         "cosmoshub-4",
		FromAddress:     "cosmos1qgas8xpptnp09lyl32kfp60hldges6guu28qmk",
		ToAddress:       "cosmos1l6vul20q74gw56fped8srkjq2x8d9m305gnxr2",
		ContractAddress: "",
		Amount:          10000,
		GasLimit:        137674,
		FeeAmount:       10000,
		Sequence:        0,
		AccountNumber:   3014650,
		Decimal:         6,
		Memo:            "10086",
		PubKey:          "03f16c9160c81b806a04da3c27d9200fe684aad79b9fcdcccfac8aa60ad7f0a56a",
	}

	txBytes, err := json.Marshal(txStruct)
	assert.NoError(t, err)

	base64Tx := base64.StdEncoding.EncodeToString(txBytes)

	chainAdaptor, err := getChainAdaptor()
	assert.NoError(t, err)

	signStr := "7e115128104f37ba74f990807c4926555a483ab88e55e157428253e686fb01a92f4ad6892bbb1a080b77795c815eb54fa94002a938ca98a3edf2f220a8b31ba001"
	request := &account.SignedTransactionRequest{
		Chain:     ChainName,
		Network:   NetWork,
		Signature: signStr,
		Base64Tx:  base64Tx,
		PublicKey: "03f16c9160c81b806a04da3c27d9200fe684aad79b9fcdcccfac8aa60ad7f0a56a",
	}
	response, err := chainAdaptor.BuildSignedTransaction(request)
	assert.NoError(t, err)
	fmt.Println("response", response)
}

func TestCosmos_SendTx(t *testing.T) {
	chainAdaptor, err := getChainAdaptor()
	assert.NoError(t, err)

	req := &account.SendTxRequest{
		RawTx: "0a98010a8e010a1c2f636f736d6f732e62616e6b2e763162657461312e4d736753656e64126e0a2d636f736d6f73317167617338787070746e7030396c796c33326b66703630686c64676573366775753238716d6b122d636f736d6f73316c3676756c323071373467773536667065643873726b6a7132783864396d333035676e7872321a0e0a057561746f6d120531303030301205313030383612660a4e0a460a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912230a2103f16c9160c81b806a04da3c27d9200fe684aad79b9fcdcccfac8aa60ad7f0a56a12040a02080112140a0e0a057561746f6d1205313030303010cab3081a417e115128104f37ba74f990807c4926555a483ab88e55e157428253e686fb01a92f4ad6892bbb1a080b77795c815eb54fa94002a938ca98a3edf2f220a8b31ba001",
	}
	response, err := chainAdaptor.SendTx(req)
	assert.NoError(t, err)
	fmt.Printf("response TxHash=%s \n", response.TxHash)
}
