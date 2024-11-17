package cosmos

import (
	"encoding/base64"
	"encoding/hex"
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
		ToAddress:       "cosmos19thxsunl9lzywglsndth5a278wtavawzzpv44q",
		ContractAddress: "",
		Amount:          10000,
		GasLimit:        137674,
		FeeAmount:       10000,
		Sequence:        1,
		AccountNumber:   3014650,
		Decimal:         6,
		Memo:            "10086",
		PubKey:          "032d535553c70dfbb9c13f32cb6d1002a4b421beff39009670e29a7e51fb88ec3f",
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
		ToAddress:       "cosmos19thxsunl9lzywglsndth5a278wtavawzzpv44q",
		ContractAddress: "",
		Amount:          10000,
		GasLimit:        137674,
		FeeAmount:       10000,
		Sequence:        1,
		AccountNumber:   3014650,
		Decimal:         6,
		Memo:            "10086",
		PubKey:          "032d535553c70dfbb9c13f32cb6d1002a4b421beff39009670e29a7e51fb88ec3f",
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
		PublicKey: "032d535553c70dfbb9c13f32cb6d1002a4b421beff39009670e29a7e51fb88ec3f",
	}
	response, err := chainAdaptor.BuildSignedTransaction(request)
	assert.NoError(t, err)
	fmt.Println("response", response)
}

func TestCosmos_SendTx(t *testing.T) {
	chainAdaptor, err := getChainAdaptor()
	assert.NoError(t, err)

	rawTx := "0A98010A8E010A1C2F636F736D6F732E62616E6B2E763162657461312E4D736753656E64126E0A2D636F736D6F73317167617338787070746E7030396C796C33326B66703630686C64676573366775753238716D6B122D636F736D6F73313974687873756E6C396C7A7977676C736E64746835613237387774617661777A7A70763434711A0E0A057561746F6D120531303030301205313030383712670A500A460A1F2F636F736D6F732E63727970746F2E736563703235366B312E5075624B657912230A21032D535553C70DFBB9C13F32CB6D1002A4B421BEFF39009670E29A7E51FB88EC3F12040A020801180312130A0D0A057561746F6D12043130303010CAB3081A4020F0753ED637D1125011466E78085BEFA9842F8591C5B92012A5E9737D8E63943BACD8046712306E0E7E821AC4C731AE6928C35FAA414D23C361186DDBFD0D8D"
	req := &account.SendTxRequest{
		RawTx: rawTx,
	}
	response, err := chainAdaptor.SendTx(req)
	assert.NoError(t, err)
	fmt.Printf("response TxHash=%s \n", response.TxHash)
	txbytes, err := hex.DecodeString(rawTx)
	assert.NoError(t, err)
	fmt.Printf("base64 RawTx=%s\n", base64.StdEncoding.EncodeToString(txbytes))
}
