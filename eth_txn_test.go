package goalchemysdk

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)


func TestAlchemyClient_eth_getTransactionByHash(t *testing.T) {
	initEnvs()
	type args struct {
		ths []TransactionByHashParam
	}
	tests := []struct {
		name    string
		c       *AlchemyClient
		args    args
		want    *AlchemyResponse[TransactionByHashResult]
		wantErr bool
	}{
		{
			name: "test empty request",
			c: &AlchemyClient{
				ApiKey:  ALCHEMY_API_KEY_TEST,
				Network: ARB_MAINNET,
				netClient: &http.Client{
					Timeout: time.Second *10,
				},
			},
			args: args{
				ths: nil,
			},

			want: &AlchemyResponse[TransactionByHashResult]{
				Id:      1,
				Jsonrpc: "2.0",
				Result:  TransactionJson{},
				Error:   ErrorExpectedAtLeastOneArgument,
			},
			wantErr: false,
		},
		{
			name: "test with results",
			c: &AlchemyClient{
				ApiKey:  ALCHEMY_API_KEY_TEST,
				Network: ARB_MAINNET,
				netClient: &http.Client{
					Timeout: time.Second *10,
				},
			},
			args: args{
				ths: []string{"0x2d6da6ea7d7d7d1ca72576dc457a2b6f59fb798566fb97492b3f2835b0a59178"},
			},

			want: &AlchemyResponse[TransactionByHashResult]{
				Id:      1,
				Jsonrpc: "2.0",
				Result: TransactionJson{
					BlockHash:            "0x8fabe002a1d4f368ac26435cf998e6d3fe408843ae6d193b7be4f49931d9ea97",
					BlockNumber:          "0x4116290",
					Hash:                 "0x2d6da6ea7d7d7d1ca72576dc457a2b6f59fb798566fb97492b3f2835b0a59178",
					AccessList:           []string{},
					ChainId:              "0xa4b1",
					From:                 "0xeba9a3b3664ce4c950cba62ed372c7815cbbfd75",
					Gas:                  "0xe5398",
					GasPrice:             "0x5f5e100",
					Input:                "0x287ad99a000000000000000000000000a6e249ffb81cf6f28ab021c3bd97620283c7335f000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000490a1",
					MaxFeePerGas:         "0x80befc0",
					MaxPriorityFeePerGas: "0x0",
					Nonce:                "0x52",
					R:                    "0xea114e9bd51b1d0cd3bf21f8c12dfb3cb2c9e9421c693d00784861ba50222f45",
					S:                    "0x424eaebe2be7963bb739a45c1b318d8f0316e116ef31450f20e2fdacd9307ac7",
					To:                   "0x5957582f020301a2f732ad17a69ab2d8b2741241",
					TransactionIndex:     "0x1",
					Type:                 "0x2",
					V:                    "0x1",
					Value:                "0x0",
				},
			},
			wantErr: false,
		},
		{
			name: "test multiples hashes error",
			c: &AlchemyClient{
				ApiKey:  ALCHEMY_API_KEY_TEST,
				Network: ARB_MAINNET,
				netClient: &http.Client{
					Timeout: time.Second *10,
				},
			},
			args: args{
				ths: []string{"0xdeada6ea7d7d7d1ca72576dc457a2b6f59fb798566fb97492b3f2835b0a59178", "0xbeefa6ea7d7d7d1ca72576dc457a2b6f59fb798566fb97492b3f2835b0a59178"},
			},

			want: &AlchemyResponse[TransactionByHashResult]{
				Id:      1,
				Jsonrpc: "2.0",
				Error:   ErrorTooManyArguments,
			},
			wantErr: false,
		},
		{
			name: "test wrong hash",
			c: &AlchemyClient{
				ApiKey:  ALCHEMY_API_KEY_TEST,
				Network: ARB_MAINNET,
				netClient: &http.Client{
					Timeout: time.Second *10,
				},
			},
			args: args{
				ths: []string{"0xdeadbeef"},
			},

			want: &AlchemyResponse[TransactionByHashResult]{
				Id:      1,
				Jsonrpc: "2.0",
				Error:   ErrorInvalidTxnHash,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Eth_getTransactionByHash(tt.args.ths)
			if (err != nil) != tt.wantErr {
				t.Errorf("AlchemyClient.eth_getTransactionByHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AlchemyClient.eth_getTransactionByHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
