package cosmos

import (
	"encoding/hex"
	"fmt"
	"testing"
)

// success
func TestCosmos_BuildUnSignTransaction(t *testing.T) {
	txStruct := &TxStructure{
		ChainId:       "cosmoshub-4",
		FromAddress:   "cosmos1qgas8xpptnp09lyl32kfp60hldges6guu28qmk",
		ToAddress:     "cosmos1l6vul20q74gw56fped8srkjq2x8d9m305gnxr2",
		Amount:        100000,
		GasLimit:      137674,
		FeeAmount:     1000,
		Sequence:      10,
		AccountNumber: 2424228,
		Decimal:       6,
		Memo:          "10086",
		PubKey:        "03f16c9160c81b806a04da3c27d9200fe684aad79b9fcdcccfac8aa60ad7f0a56a",
	}
	txBytes, _ := BuildUnSignTransaction(txStruct)
	fmt.Printf("txBytes=%X\n", txBytes)
}

func TestCosmos_CreateSignTransaction(t *testing.T) {
	txStruct := &TxStructure{
		ChainId:       "cosmoshub-4",
		FromAddress:   "cosmos1qgas8xpptnp09lyl32kfp60hldges6guu28qmk",
		ToAddress:     "cosmos1l6vul20q74gw56fped8srkjq2x8d9m305gnxr2",
		Amount:        100000,
		GasLimit:      137674,
		FeeAmount:     1000,
		Sequence:      10,
		AccountNumber: 2424228,
		Decimal:       6,
		Memo:          "10086",
		PubKey:        "03f16c9160c81b806a04da3c27d9200fe684aad79b9fcdcccfac8aa60ad7f0a56a",
	}
	signature := "0A9F010A95010A1C2F636F736D6F732E62616E6B2E763162657461312E4D736753656E6412750A2D636F736D6F73317167617338787070746E7030396C796C33326B66703630686C64676573366775753238716D6B122D636F736D6F73316C3676756C323071373467773536667065643873726B6A7132783864396D333035676E7872321A150A057561746F6D120C31303030303030303030303012053130303836126D0A500A460A1F2F636F736D6F732E63727970746F2E736563703235366B312E5075624B657912230A2103F16C9160C81B806A04DA3C27D9200FE684AAD79B9FCDCCCFAC8AA60AD7F0A56A12040A020801180A12190A130A057561746F6D120A3130303030303030303010CAB3081A0B636F736D6F736875622D3420A4FB9301"
	signBytes, _ := hex.DecodeString(signature)
	txBytes, _ := BuildSignTransaction(txStruct, signBytes)
	fmt.Printf("txBytes=%X\n", txBytes)
}
