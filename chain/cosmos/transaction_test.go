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
		FromAddress:   "",
		ToAddress:     "",
		Amount:        10000,
		GasLimit:      137674,
		FeeAmount:     1000,
		Sequence:      4,
		AccountNumber: 3014650,
		Decimal:       6,
		Memo:          "10087",
		PubKey:        "",
	}
	txBytes, _ := BuildUnSignTransaction(txStruct)
	fmt.Printf("txBytes=%X\n", txBytes)
}

func TestCosmos_CreateSignTransaction(t *testing.T) {
	txStruct := &TxStructure{
		ChainId:       "cosmoshub-4",
		FromAddress:   "",
		ToAddress:     "",
		Amount:        10000,
		GasLimit:      137674,
		FeeAmount:     1000,
		Sequence:      4,
		AccountNumber: 3014650,
		Decimal:       6,
		Memo:          "10087",
		PubKey:        "",
	}
	signature := "20f0753ed637d1125011466e78085befa9842f8591c5b92012a5e9737d8e63943bacd8046712306e0e7e821ac4c731ae6928c35faa414d23c361186ddbfd0d8d"
	signBytes, _ := hex.DecodeString(signature)
	txBytes, _ := BuildSignTransaction(txStruct, signBytes)
	fmt.Printf("txBytes=%X\n", txBytes)
}
