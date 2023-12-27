package goalchemysdk

import (
	"fmt"
	"time"
)

type Network string

const (
	ETH_MAINNET   Network = "eth-mainnet"
	ETH_GOERLI    Network = "eth-goerli"
	ETH_SEPOLIA   Network = "eth-sepolia"
	MATIC_MAINNET Network = "polygon-mainnet"
	MATIC_MUMBAI  Network = "polygon-mumbai"
	OPT_MAINNET   Network = "opt-mainnet"
	OPT_GOERLI    Network = "opt-goerli"
	OPT_KOVAN     Network = "opt-kovan"
	ARB_MAINNET   Network = "arb-mainnet"
	ARB_GOERLI    Network = "arb-goerli"
	ASTAR_MAINNET Network = "astar-mainnet"
)

type BlockTag string

const (
	PENDING   BlockTag = "pending"   // - A sample next block built by the client on top of latest and containing the set of transactions usually taken from local mempool. Intuitively, you can think of these as blocks that have not been mined yet.
	LATEST    BlockTag = "latest"    // - The most recent block in the canonical chain observed by the client, this block may be re-orged out of the canonical chain even under healthy/normal conditions.
	SAFE      BlockTag = "safe"      // - The most recent crypto-economically secure block, cannot be re-orged outside of manual intervention driven by community coordination. Intuitively, this block is “unlikely” to be re-orged.
	FINALIZED BlockTag = "finalized" // - The most recent crypto-economically secure block, that has been accepted by >2/3 of validators. Cannot be re-orged outside of manual intervention driven by community coordination. Intuitively, this block is very unlikely to be re-orged.
	EARLIEST  BlockTag = "earliest"  // - The lowest numbered block the client has available. Intuitively, you can think of this as the first block created.
)

type JsonParams[P any] struct {
	Id      uint   `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []P    `json:"params"`
}

type AlchemyApiError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// custom error type
func (ae *AlchemyApiError) Error() string {
	return fmt.Sprintf("Alchemy Api error code=%d, message=%s", ae.Code, ae.Message)
}

type AlchemyResponse[R any] struct {
	Id      uint            `json:"id"`
	Jsonrpc string          `json:"jsonrpc"`
	Result  R               `json:"result,omitempty"` // []map[string]any `json:"result"`
	Error   AlchemyApiError `json:"error,omitempty"`
}

type RetriableError struct {
	Err        error
	RetryAfter time.Duration
}

// Error returns error message and a Retry-After duration
func (e *RetriableError) Error() string {
	return fmt.Sprintf("%s (retry after %v)", e.Err.Error(), e.RetryAfter)
}

// Alchemy Errors
var (
	ErrorMustBeAuthenticated = AlchemyApiError{
		Code:    -32600,
		Message: "Must be authenticated!",
	}

	ErrorExpectedAtLeastOneArgument = AlchemyApiError{
		Code:    -32602,
		Message: "expected 'params' array of at least 1 argument",
	}

	ErrorInvalidTxnHash = AlchemyApiError{
		Code:    -32602,
		Message: "invalid 1st argument: transaction_hash value was too short",
	}

	ErrorTooManyArguments = AlchemyApiError{
		Code:    -32602,
		Message: "too many arguments, want at most 1",
	}
	ErrorTooShortAddress = AlchemyApiError{
		Code: -32602,
		Message: "invalid 1st argument: address value was too short",
	}
	ErrorTooLongAddress = AlchemyApiError{
		Code: -32602,
		Message: "invalid 1st argument: address value was too long",
	}
	ErrorInvalidAddress = AlchemyApiError{
		Code: -32602,
		Message: "invalid 1st argument: address value was not valid hexadecimal",
	}
)

// Alchemy Errors functions
func ErrorWrongMethod(methodName string) AlchemyApiError {
	return AlchemyApiError{
		Code:    -32600,
		Message: fmt.Sprintf("Unsupported method: %s. See available methods at https://docs.alchemy.com/alchemy/documentation/apis", methodName),
	}
}
