package goalchemysdk

// types

type LogsParam struct {
	BlockHash string      `json:"blockHash,omitempty"`
	Address   string      `json:"address,omitempty"`
	FromBlock BlockFilter `json:"fromBlock,omitempty"`
	ToBlock   BlockFilter `json:"toBlock,omitempty"`
	Topics    []string    `json:"topics,omitempty"`
}

type LogsResult struct {
	Address          string   `json:"address,omitempty"`
	Topics           []string `json:"topics,omitempty"`
	Data             string   `json:"data,omitempty"`
	BlockNumber      string   `json:"blockNumber,omitempty"`
	TransactionHash  string   `json:"transactionHash,omitempty"`
	TransactionIndex string   `json:"transactionIndex,omitempty"`
	BlockHash        string   `json:"BlockHash,omitempty"`
	LogIndex         string   `json:"logIndex,omitempty"`
	Removed          bool     `json:"removed,omitempty"`
}
type LogsResults = []LogsResult

//queries
func (c *AlchemyClient) Eth_getLogs(lp []LogsParam) (*AlchemyResponse[LogsResults], error) {
	j := JsonParams[LogsParam]{
		Id:      1,
		Jsonrpc: "2.0",
		Method:  "eth_getLogs",
		Params:  lp,
	}
	return executePost[LogsParam, LogsResults](c, j)
}