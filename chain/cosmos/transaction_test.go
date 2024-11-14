package cosmos

import (
	"fmt"
	"testing"
)

// success
func TestCosmos_CreateUnSignTransaction(t *testing.T) {
	txStruct := &TxStructure{
		FromAddress:   "cosmos1qgas8xpptnp09lyl32kfp60hldges6guu28qmk",
		ToAddress:     "cosmos1l6vul20q74gw56fped8srkjq2x8d9m305gnxr2",
		Amount:        100000,
		GasLimit:      137674,
		FeeAmount:     1000,
		Sequence:      10,
		accountNumber: "2424228",
	}
	txBytes, _ := BuildUnSignTransaction(txStruct)
	fmt.Printf("txBytes=%X\n", txBytes)
}
