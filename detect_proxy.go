package goalchemysdk

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	//"sync"
)

type ProxyResult struct {
	address string
	err     error
}

type ProxyDetectorFunc func(*AlchemyClient, string, BlockTag, chan ProxyResult)

func createJob(jobCounter *uint, detector ProxyDetectorFunc, c *AlchemyClient, addr string, bt BlockTag, out chan ProxyResult) {
	*jobCounter++
	go detector(c, addr, bt, out)
}

func (c *AlchemyClient) DetectProxyTarget(proxyAddress string, blockTag BlockTag) (address string, err error) {
	if blockTag == "" {
		blockTag = LATEST
	}
	address = "0x"
	err = errors.New("no Proxy detected")

	res := make(chan ProxyResult)
	done := make(chan bool)
	jobs := uint(0)

	detectors := []ProxyDetectorFunc{checkEIP1167, checkEIP1967Direct,
		checkEIP1967Beacon, checkOpenZeppelin, checkEIP1822,
		checkEIP897, checkGnosisSafe, checkComptroller}

	for _, f := range detectors {
		createJob(&jobs, f, c, proxyAddress, blockTag, res)
	}

	// exit on valid result routine
	go func(res chan ProxyResult, done chan bool) {
		counter := uint(0)
		for {
			val := <-res
			counter++
			if val.err == nil {
				address = val.address
				err = val.err
				done <- true
				break
			}
			if counter >= jobs {
				address = "0x"
				err = errors.New("no proxy found")
				done <- true
				break
			}
		}
	}(res, done)

	// fmt.Println("waiting for done!!!")
	<-done
	// fmt.Println("returned !!!")
	return address, err
}

func readAddress(address string) (string, error) {
	if address == "0x" || address == "" {
		return "0x", errors.New("invalid address 0x")
	}

	if len(address) == 66 {
		// fmt.Printf("66 address ======> %s\n", address)
		address = fmt.Sprintf("0x%s", address[len(address)-40:])

	}

	zeroAddress := fmt.Sprintf("0x%s", strings.Repeat("0", 40))
	if address == zeroAddress {
		return "0x", errors.New("zero address")
	}
	return address, nil
}

// storage based detection
func checkWithStorage(c *AlchemyClient, proxyAddress string, blockTag BlockTag, res chan ProxyResult, slot string) {
	resp, err := c.Eth_getStorageAt(proxyAddress, slot, blockTag)
	
	if err != nil {
		res <- ProxyResult{
			address: "0x",
			err:     err,
		}
		return
	}
	address, err := readAddress(resp.Result)

	if err != nil {
		res <- ProxyResult{
			address: "0x",
			err:     err,
		}
		return
	}

	res <- ProxyResult{
		address: address,
		err:     nil,
	}
}

// OpenZeppelin proxy pattern
func checkOpenZeppelin(c *AlchemyClient, proxyAddress string, blockTag BlockTag, res chan ProxyResult) {
	checkWithStorage(c, proxyAddress, blockTag, res, OPEN_ZEPPELIN_IMPLEMENTATION_SLOT)
}

// EIP-1822 Universal Upgradeable Proxy Standard
func checkEIP1822(c *AlchemyClient, proxyAddress string, blockTag BlockTag, res chan ProxyResult) {
	checkWithStorage(c, proxyAddress, blockTag, res, EIP_1822_LOGIC_SLOT)
}

// EIP-897 DelegateProxy pattern
func checkEIP897(c *AlchemyClient, proxyAddress string, blockTag BlockTag, res chan ProxyResult) {
	address, err := getAddressFromBeacon(c, proxyAddress, EIP_897_INTERFACE[0])
	if err != nil {
		res <- ProxyResult{
			address: "0x",
			err:     err,
		}
		return
	}
	res <- ProxyResult{
		address: address,
		err:     err,
	}
}

// GnosisSafeProxy contract
func checkGnosisSafe(c *AlchemyClient, proxyAddress string, blockTag BlockTag, res chan ProxyResult) {
	address, err := getAddressFromBeacon(c, proxyAddress, GNOSIS_SAFE_PROXY_INTERFACE[0])
	if err != nil {
		res <- ProxyResult{
			address: "0x",
			err:     err,
		}
		return
	}
	res <- ProxyResult{
		address: address,
		err:     err,
	}
}

// Comptroller proxy
func checkComptroller(c *AlchemyClient, proxyAddress string, blockTag BlockTag, res chan ProxyResult) {
	address, err := getAddressFromBeacon(c, proxyAddress, COMPTROLLER_PROXY_INTERFACE[0])
	if err != nil {
		res <- ProxyResult{
			address: "0x",
			err:     err,
		}
		return
	}
	res <- ProxyResult{
		address: address,
		err:     err,
	}
}

// EIP-1967 direct proxy
func checkEIP1967Direct(c *AlchemyClient, proxyAddress string, blockTag BlockTag, res chan ProxyResult) {
	checkWithStorage(c, proxyAddress, blockTag, res, EIP_1967_LOGIC_SLOT)
}

// EIP-1967 beacon proxy
func checkEIP1967Beacon(c *AlchemyClient, proxyAddress string, blockTag BlockTag, res chan ProxyResult) {
	resp, err := c.Eth_getStorageAt(proxyAddress, EIP_1967_BEACON_SLOT, blockTag)
	if err != nil {
		res <- ProxyResult{
			address: "0x",
			err:     err,
		}
		return
	}
	beaconAddress, err := readAddress(resp.Result)
	if err != nil {
		res <- ProxyResult{
			address: "0x",
			err:     err,
		}
		return
	}

	address, err := getAddressFromBeacon(c, beaconAddress, EIP_1167_BEACON_METHODS[0])

	if err != nil {
		address, err = getAddressFromBeacon(c, beaconAddress, EIP_1167_BEACON_METHODS[1])
		if err != nil {
			res <- ProxyResult{
				address: "0x",
				err:     err,
			}
			return
		}
		res <- ProxyResult{
			address: address,
			err:     err,
		}
		return
	}

	res <- ProxyResult{
		address: address,
		err:     err,
	}
}

func getAddressFromBeacon(c *AlchemyClient, proxyAddress string, methodEncoded string) (string, error) {
	resp, err := c.Eth_call(CallTxn{To: proxyAddress, Data: methodEncoded}, LATEST)
	if err != nil {
		return "0x", err
	}
	address, err := readAddress(resp.Result)
	if err != nil {
		return "0x", err
	}
	return address, nil
}

func checkEIP1167(c *AlchemyClient, proxyAddress string, blockTag BlockTag, res chan ProxyResult) {
	resp, err := c.Eth_getCode(proxyAddress, blockTag)
	if err != nil {
		res <- ProxyResult{
			address: "0x",
			err:     err,
		}
		return
	}
	addr, err := parse1167Bytecode(resp.Result)

	if err != nil {
		res <- ProxyResult{
			address: "0x",
			err:     err,
		}
		return
	}

	address, err := readAddress(addr)
	if err != nil {
		res <- ProxyResult{
			address: "0x",
			err:     err,
		}
		return
	}

	res <- ProxyResult{
		address: address,
		err:     nil,
	}

}

// EIP-1167 Minimal Proxy Contract parse
func parse1167Bytecode(byteCode string) (string, error) {

	if !strings.HasPrefix(byteCode, EIP_1167_BYTECODE_PREFIX) {
		return "0x", errors.New("eip 1167 bytecode prefix not found, not an EIP-1167 bytecode")
	}

	// detect length of address (20 bytes non-optimized, 0 < N < 20 bytes for vanity addresses)
	pushNHex := byteCode[LEN_EIP_1167_BYTECODE_PREFIX : LEN_EIP_1167_BYTECODE_PREFIX+2]
	// push1 ... push20 use opcodes 0x60 ... 0x73
	pushNInt, err := strconv.ParseInt(pushNHex, 16, 32)
	if err != nil {
		return "0x", errors.New("invalid pushN integer, Not an EIP-1167 bytecode")
	}
	addressLength := int(pushNInt - 0x5f)

	if addressLength < 1 || addressLength > 20 {
		return "0x", errors.New("invalid address length, Not an EIP-1167 bytecode")
	}

	addressFromBytecode := byteCode[LEN_EIP_1167_BYTECODE_PREFIX+2 : LEN_EIP_1167_BYTECODE_PREFIX+2+addressLength*2] // address length is in bytes, 2 hex chars make up 1 byte

	suffix := byteCode[LEN_EIP_1167_BYTECODE_PREFIX+2+addressLength*2+SUFFIX_OFFSET_FROM_ADDRESS_END:]

	if !strings.HasPrefix(suffix, EIP_1167_BYTECODE_SUFFIX) {
		return "0x", errors.New("eip 1167 bytecode suffix not found, Not an EIP-1167 bytecode")
	}

	// padStart is needed for vanity addresses
	return fmt.Sprintf("0x%0*s", 40, addressFromBytecode), nil
}
