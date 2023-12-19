package goalchemysdk

type TransactionByHashParam = string

type TransactionByHashResult = TransactionJson

type TransactionJson struct {
	BlockHash            string   `json:"blockHash,omitempty"`
	BlockNumber          string   `json:"blockNumber,omitempty"`
	Hash                 string   `json:"hash,omitempty"`
	AccessList           []string `json:"accessList,omitempty"`
	ChainId              string   `json:"chainId,omitempty"`
	From                 string   `json:"from,omitempty"`
	Gas                  string   `json:"gas,omitempty"`
	GasPrice             string   `json:"gasPrice,omitempty"`
	Input                string   `json:"input,omitempty"`
	MaxFeePerGas         string   `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas string   `json:"maxPriorityFeePerGas,omitempty"`
	Nonce                string   `json:"nonce,omitempty"`
	R                    string   `json:"r,omitempty"`
	S                    string   `json:"s,omitempty"`
	To                   string   `json:"to,omitempty"`
	TransactionIndex     string   `json:"transactionIndex,omitempty"`
	Type                 string   `json:"type,omitempty"`
	V                    string   `json:"v,omitempty"`
	Value                string   `json:"value,omitempty"`
}

func (c *AlchemyClient) eth_getTransactionByHash(ths []TransactionByHashParam) (*AlchemyResponse[TransactionByHashResult], error) {
	j := JsonParams[TransactionByHashParam]{
		Id:      1,
		Jsonrpc: "2.0",
		Method:  "eth_getTransactionByHash",
		Params:  ths,
	}
	return executePost[TransactionByHashParam, TransactionByHashResult](c, j)
}
