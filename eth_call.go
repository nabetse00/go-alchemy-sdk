package goalchemysdk

// getCode Params
// String - 20 Bytes - Address
// String - Either the hex value of a block number OR a block hash OR One of the following block tags:
// pending - A sample next block built by the client on top of latest and containing the set of transactions usually taken from local mempool. Intuitively, you can think of these as blocks that have not been mined yet.
// latest - The most recent block in the canonical chain observed by the client, this block may be re-orged out of the canonical chain even under healthy/normal conditions.
// earliest - The lowest numbered block the client has available. Intuitively, you can think of this as the first block created.
type CallTxn struct{
	From string `json:"from,omitempty"`
	To string `json:"to"`
	Data string `json:"data"`
	Gas string `json:"gas,omitempty"`
	GasPrice string `json:"gasPrice,omitempty"`
	Value string `json:"value,omitempty"`
	BlockHash string `json:"blockHash,omitempty"`
	RequireCanonical bool `json:"requireCanonical,omitempty"`
}

type CallBlk = BlockTag


type CallResult = string

func (c *AlchemyClient) Eth_call(txn CallTxn,  blk CallBlk) (*AlchemyResponse[CallResult], error) {
	j := JsonParams[CallTxn]{
		Id:      1,
		Jsonrpc: "2.0",
		Method:  "eth_call",
		Params:  []CallTxn{txn},
	}
	return executePost[CallTxn, CallResult](c, j)
}
