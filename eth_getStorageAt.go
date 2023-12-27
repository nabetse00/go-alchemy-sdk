package goalchemysdk

type GetStorageAtParam = string

type GetStorageAtResult = string

func (c *AlchemyClient) Eth_getStorageAt(address string, id string, blocktag BlockTag) (*AlchemyResponse[GetCodeResult], error) {
	j := JsonParams[GetCodeParam]{
		Id:      1,
		Jsonrpc: "2.0",
		Method:  "eth_getStorageAt",
		Params:  []GetCodeParam{address, id, string(blocktag)},
	}
	return executePost[GetStorageAtParam, GetStorageAtResult](c, j)
}