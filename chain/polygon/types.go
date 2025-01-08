package polygon

type Tx struct {
	ChainId     string `json:"chain_id"`
	Nonce       uint64 `json:"nonce"`
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
	GasLimit    uint64 `json:"gas_limit"`
	Gas         uint64 `json:"Gas"`

	MaxFeePerGas         string `json:"max_fee_per_gas"`
	MaxPriorityFeePerGas string `json:"max_priority_fee_per_gas"`
	Signature            string `json:"signature,omitempty"`

	// eth/ethereum amount
	Amount string `json:"amount"`
	// ethereum erc721 erc1155 contract_address
	ContractAddress string `json:"contract_address"`
}
