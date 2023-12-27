package goalchemysdk

// getCode Params
// String - 20 Bytes - Address
// String - Either the hex value of a block number OR a block hash OR One of the following block tags:
// pending - A sample next block built by the client on top of latest and containing the set of transactions usually taken from local mempool. Intuitively, you can think of these as blocks that have not been mined yet.
// latest - The most recent block in the canonical chain observed by the client, this block may be re-orged out of the canonical chain even under healthy/normal conditions.
// earliest - The lowest numbered block the client has available. Intuitively, you can think of this as the first block created.
type GetCodeParam = string

type GetCodeResult = string

func (c *AlchemyClient) Eth_getCode(address string, blocktag BlockTag) (*AlchemyResponse[GetCodeResult], error) {
	j := JsonParams[GetCodeParam]{
		Id:      1,
		Jsonrpc: "2.0",
		Method:  "eth_getCode",
		Params:  []string{address, string(blocktag)},
	}
	return executePost[GetCodeParam, GetCodeResult](c, j)
}
