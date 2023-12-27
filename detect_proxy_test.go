package goalchemysdk

import (
	"errors"
	"testing"
)

func Test_parse1167Bytecode(t *testing.T) {
	type args struct {
		byteCode string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Parse eip 1167 - normal address",
			args: args{
				byteCode: "0x363d3d373d3d3d363d73f62849f9a0b5bf2913b396098f7c7019b51a820a5af43d82803e903d91602b57fd5bf3000000000000000000000000000000000000000000000000000000000000007a6900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			},
			want:    "0xf62849f9a0b5bf2913b396098f7c7019b51a820a",
			wantErr: false,
		},
		{
			name: "Parse eip 1167 - vanity address padding",
			args: args{
				byteCode: "0x363d3d373d3d3d363d6f10fd301be3200e67978e3cc67c962f485af43d82803e903d91602757fd5bf3",
			},
			want:    "0x0000000010fd301be3200e67978e3cc67c962f48",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse1167Bytecode(tt.args.byteCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse1167Bytecode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parse1167Bytecode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetectProxyTarget(t *testing.T) {
	initEnvs()
	type args struct {
		c            *AlchemyClient
		proxyAddress string
		blockTag     BlockTag
	}
	tests := []struct {
		name        string
		args        args
		wantAddress string
		wantErr     bool
	}{

		{
			name: "test no proxy",
			args: args{
				c: &AlchemyClient{
					ApiKey:   ALCHEMY_API_KEY_TEST,
					Network:  ETH_MAINNET,
					MaxRetry: 10,
				},
				proxyAddress: "0xdead3fd06D57D140f6ad8C2913DbE87fdecDd5F",
				blockTag:     LATEST,
			},
			wantAddress: "0x",
			wantErr:     true,
		},
		{
			name: "test detect proxy EIP1167",
			args: args{
				c: &AlchemyClient{
					ApiKey:   ALCHEMY_API_KEY_TEST,
					Network:  ETH_MAINNET,
					MaxRetry: 10,
				},
				proxyAddress: "0xa81043fd06D57D140f6ad8C2913DbE87fdecDd5F",
				blockTag:     LATEST,
			},
			wantAddress: "0x0000000010fd301be3200e67978e3cc67c962f48",
			wantErr:     false,
		},
		{
			name: "test detect proxy EIP1967 Direct Proxy",
			args: args{
				c: &AlchemyClient{
					ApiKey:   ALCHEMY_API_KEY_TEST,
					Network:  ETH_MAINNET,
					MaxRetry: 10,
				},
				proxyAddress: "0xA7AeFeaD2F25972D80516628417ac46b3F2604Af",
				blockTag:     LATEST,
			},
			wantAddress: "0x4bd844f72a8edd323056130a86fc624d0dbcf5b0",
			wantErr:     false,
		},
		{
			name: "test detects EIP1967 beacon proxies",
			args: args{
				c: &AlchemyClient{
					ApiKey:   ALCHEMY_API_KEY_TEST,
					Network:  ETH_MAINNET,
					MaxRetry: 10,
				},
				proxyAddress: "0xDd4e2eb37268B047f55fC5cAf22837F9EC08A881",
				blockTag:     LATEST,
			},
			wantAddress: "0xe5c048792dcf2e4a56000c8b6a47f21df22752d1",
			wantErr:     false,
		},
		{
			name: "test detects EIP1967 beacon variant proxies",
			args: args{
				c: &AlchemyClient{
					ApiKey:   ALCHEMY_API_KEY_TEST,
					Network:  ETH_MAINNET,
					MaxRetry: 10,
				},
				proxyAddress: "0x114f1388fAB456c4bA31B1850b244Eedcd024136",
				blockTag:     LATEST,
			},
			wantAddress: "0x36b799160cdc2d9809d108224d1967cc9b7d321c",
			wantErr:     false,
		},
		{
			name: "test detects OpenZeppelin proxies",
			args: args{
				c: &AlchemyClient{
					ApiKey:   ALCHEMY_API_KEY_TEST,
					Network:  ETH_MAINNET,
					MaxRetry: 10,
				},
				proxyAddress: "0x8260b9eC6d472a34AD081297794d7Cc00181360a",
				blockTag:     LATEST,
			},
			wantAddress: "0xe4e4003afe3765aca8149a82fc064c0b125b9e5a",
			wantErr:     false,
		},
		{
			name: "test detects EIP-897 delegate proxies",
			args: args{
				c: &AlchemyClient{
					ApiKey:   ALCHEMY_API_KEY_TEST,
					Network:  ETH_MAINNET,
					MaxRetry: 10,
				},
				proxyAddress: "0x8260b9eC6d472a34AD081297794d7Cc00181360a",
				blockTag:     LATEST,
			},
			wantAddress: "0xe4e4003afe3765aca8149a82fc064c0b125b9e5a",
			wantErr:     false,
		},
		{
			name: "test detects EIP-1167 minimal proxies",
			args: args{
				c: &AlchemyClient{
					ApiKey:   ALCHEMY_API_KEY_TEST,
					Network:  ETH_MAINNET,
					MaxRetry: 10,
				},
				proxyAddress: "0x6d5d9b6ec51c15f45bfa4c460502403351d5b999",
				blockTag:     LATEST,
			},
			wantAddress: "0x210ff9ced719e9bf2444dbc3670bac99342126fa",
			wantErr:     false,
		},
		{
			name: "test detects EIP-1167 minimal proxies with vanity addresses",
			args: args{
				c: &AlchemyClient{
					ApiKey:   ALCHEMY_API_KEY_TEST,
					Network:  ETH_MAINNET,
					MaxRetry: 10,
				},
				proxyAddress: "0xa81043fd06D57D140f6ad8C2913DbE87fdecDd5F",
				blockTag:     LATEST,
			},
			wantAddress: "0x0000000010fd301be3200e67978e3cc67c962f48",
			wantErr:     false,
		},
		{
			name: "test detects Gnosis Safe proxies",
			args: args{
				c: &AlchemyClient{
					ApiKey:   ALCHEMY_API_KEY_TEST,
					Network:  ETH_MAINNET,
					MaxRetry: 10,
				},
				proxyAddress: "0x0DA0C3e52C977Ed3cBc641fF02DD271c3ED55aFe",
				blockTag:     LATEST,
			},
			wantAddress: "0xd9db270c1b5e3bd161e8c8503c55ceabee709552",
			wantErr:     false,
		},
		{
			name: "test detects Compound's custom proxy",
			args: args{
				c: &AlchemyClient{
					ApiKey:   ALCHEMY_API_KEY_TEST,
					Network:  ETH_MAINNET,
					MaxRetry: 10,
				},
				proxyAddress: "0x3d9819210A31b4961b30EF54bE2aeD79B9c9Cd3B",
				blockTag:     LATEST,
			},
			wantAddress: "0xbafe01ff935c7305907c33bf824352ee5979b526",
			wantErr:     false,
		},
		{
			name: "test detects ARB token proxy [EIP1967]",
			args: args{
				c: &AlchemyClient{
					ApiKey:   ALCHEMY_API_KEY_TEST,
					Network:  ARB_MAINNET,
					MaxRetry: 10,
				},
				proxyAddress: "0x912ce59144191c1204e64559fe8253a0e49e6548",
				blockTag:     LATEST,
			},
			wantAddress: "0xc4ed0a9ea70d5bcc69f748547650d32cc219d882",
			wantErr:     false,
		},
		{
			name: "test detects Abracadabra [GnosisSafeProxy]",
			args: args{
				c: &AlchemyClient{
					ApiKey:   ALCHEMY_API_KEY_TEST,
					Network:  ARB_MAINNET,
					MaxRetry: 10,
				},
				proxyAddress: "0xfBDf75866904767dE1Caa8B64eb18a7562517F5A",
				blockTag:     LATEST,
			},
			wantAddress: "0x3e5c63644e683549055b9be8653de26e0b4cd36e",
			wantErr:     false,
		},
		{
			name: "test detects MUX [GnosisSafeProxy]",
			args: args{
				c: &AlchemyClient{
					ApiKey:   ALCHEMY_API_KEY_TEST,
					Network:  ARB_MAINNET,
					MaxRetry: 10,
				},
				proxyAddress: "0x4Fa610DD115e790B8768A482Fc366803534e9Adc",
				blockTag:     LATEST,
			},
			wantAddress: "0x3e5c63644e683549055b9be8653de26e0b4cd36e",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAddress, err := tt.args.c.DetectProxyTarget(tt.args.proxyAddress, tt.args.blockTag)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectProxyTarget() error = %v, wantErr %v, got %v", err, tt.wantErr, gotAddress)
				return
			}
			if gotAddress != tt.wantAddress {
				t.Errorf("DetectProxyTarget() = %v, want %v", gotAddress, tt.wantAddress)
			}
		})
	}
}

func Test_checkWithStorage(t *testing.T) {
	initEnvs()
	type args struct {
		c            *AlchemyClient
		proxyAddress string
		blockTag     BlockTag
		res          chan ProxyResult
		slot         string
	}
	tests := []struct {
		name string
		args args
		want ProxyResult
	}{
		{
			name: "test wrong api error",
			args: args{
				c: &AlchemyClient{
					ApiKey:   ALCHEMY_API_KEY_TEST,
					Network:  "wrong",
					MaxRetry: 3,
				},
				proxyAddress: "0x4Fa610DD115e790B8768A482Fc366803534e9Adc",
				res:          make(chan ProxyResult),
				slot:         "FAKE_SLOT",
				blockTag:     LATEST,
			},
			want: ProxyResult{
				address: "0x",
				err:     errors.New("some error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go checkWithStorage(tt.args.c, tt.args.proxyAddress, tt.args.blockTag, tt.args.res, tt.args.slot)
			result := <-tt.args.res

			if result.address != tt.want.address {
				t.Errorf("checkWithStorage() invalid address = %v, want %v", result.address, tt.want.address)
			}

			if result.err.Error() == "" {
				t.Errorf("checkWithStorage() error wanted ! [%s]", result.err)
			}
		})
	}
}

func Test_checkEIP1167(t *testing.T) {
	initEnvs()
	type args struct {
		c            *AlchemyClient
		proxyAddress string
		blockTag     BlockTag
		res          chan ProxyResult
	}
	tests := []struct {
		name string
		args args
		want ProxyResult
	}{
		{
			name: "test wrong api error",
			args: args{
				c: &AlchemyClient{
					ApiKey:   ALCHEMY_API_KEY_TEST,
					Network:  "wrong",
					MaxRetry: 3,
				},
				proxyAddress: "0x4Fa610DD115e790B8768A482Fc366803534e9Adc",
				res:          make(chan ProxyResult),
				blockTag:     LATEST,
			},
			want: ProxyResult{
				address: "0x",
				err:     errors.New("some error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go checkEIP1167(tt.args.c, tt.args.proxyAddress, tt.args.blockTag, tt.args.res)
			result := <-tt.args.res

			if result.address != tt.want.address {
				t.Errorf("checkEIP1167(() invalid address = %v, want %v", result.address, tt.want.address)
			}

			if result.err.Error() == "" {
				t.Errorf("checkEIP1167(() error wanted ! [%s]", result.err)
			}
		})
	}
}

func Test_checkEIP1967Beacon(t *testing.T) {
	initEnvs()
	type args struct {
		c             *AlchemyClient
		proxyAddress  string
		blockTag      BlockTag
		res           chan ProxyResult
		beaconMethods []string
	}
	tests := []struct {
		name string
		args args
		want ProxyResult
	}{
		{
			name: "test wrong api error",
			args: args{
				c: &AlchemyClient{
					ApiKey:   ALCHEMY_API_KEY_TEST,
					Network:  "wrong",
					MaxRetry: 3,
				},
				proxyAddress:  "0x4Fa610DD115e790B8768A482Fc366803534e9Adc",
				res:           make(chan ProxyResult),
				beaconMethods: EIP_1167_BEACON_METHODS,
				blockTag:      LATEST,
			},
			want: ProxyResult{
				address: "0x",
				err:     errors.New("some error"),
			},
		},
		{
			name: "test alt method errrors",
			args: args{
				c: &AlchemyClient{
					ApiKey:   ALCHEMY_API_KEY_TEST,
					Network:  ETH_MAINNET,
					MaxRetry: 10,
				},
				proxyAddress:  "0x114f1388fAB456c4bA31B1850b244Eedcd024136",
				res:           make(chan ProxyResult),
				blockTag:      LATEST,
				beaconMethods: []string{EIP_1167_BEACON_METHODS[1], EIP_1167_BEACON_METHODS[0]},
			},
			want: ProxyResult{
				address: "0x36b799160cdc2d9809d108224d1967cc9b7d321c",
				err:     nil,
			},
		},
		{
			name: "test alt method errrors",
			args: args{
				c: &AlchemyClient{
					ApiKey:   ALCHEMY_API_KEY_TEST,
					Network:  ETH_MAINNET,
					MaxRetry: 10,
				},
				proxyAddress:  "0x114f1388fAB456c4bA31B1850b244Eedcd024136",
				res:           make(chan ProxyResult),
				blockTag:      LATEST,
				beaconMethods: []string{"WROMG", EIP_1167_BEACON_METHODS[1]},
			},
			want: ProxyResult{
				address: "0x36b799160cdc2d9809d108224d1967cc9b7d321c",
				err:     nil,
			},
		},

		{
			name: "test alt method errrors",
			args: args{
				c: &AlchemyClient{
					ApiKey:   ALCHEMY_API_KEY_TEST,
					Network:  ETH_MAINNET,
					MaxRetry: 10,
				},
				proxyAddress:  "0x114f1388fAB456c4bA31B1850b244Eedcd024136",
				res:           make(chan ProxyResult),
				blockTag:      LATEST,
				beaconMethods: []string{"WROMG", "STUFF"},
			},
			want: ProxyResult{
				address: "0x",
				err:     errors.New("some error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//EIP_1167_BEACON_METHODS[1] = EIP_1167_BEACON_METHODS[0]
			//EIP_1167_BEACON_METHODS[0] = "wrong"
			EIP_1167_BEACON_METHODS = tt.args.beaconMethods
			go checkEIP1967Beacon(tt.args.c, tt.args.proxyAddress, tt.args.blockTag, tt.args.res)
			result := <-tt.args.res

			if result.address != tt.want.address {
				t.Errorf("checkEIP1967Beacon() invalid address = %v, want %v", result.address, tt.want.address)
				return
			}

			if tt.want.err != nil && result.err == nil {
				t.Errorf("checkEIP1967Beacon() error wanted ! [%s]", result.err)
				return
			}
		})
	}
}
