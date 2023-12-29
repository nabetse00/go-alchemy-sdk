package goalchemysdk

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestAlchemyClient_eth_getLogs(t *testing.T) {
	initEnvs()
	type args struct {
		lps []LogsParam
	}
	tests := []struct {
		name    string
		c       *AlchemyClient
		args    args
		want    *AlchemyResponse[LogsResults]
		wantErr bool
	}{
		{
			name: "test empty results",
			c: &AlchemyClient{
				ApiKey:  ALCHEMY_API_KEY_TEST,
				Network: ARB_MAINNET,
				netClient: &http.Client{
					Timeout: time.Second *10,
				},
			},
			args: args{
				lps: nil,
			},

			want: &AlchemyResponse[LogsResults]{
				Id:      1,
				Jsonrpc: "2.0",
				Result:  nil,
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
				lps: []LogsParam{
					{
						BlockHash: "0x3ff6a0c14a272c9379838543735edf677fbe718df12ae52e921fc20f499f6feb",
						Address:   "0x2cde9919e81b20b4b33dd562a48a84b54c48f00c",
					},
				},
			},
			want: &AlchemyResponse[LogsResults]{
				Id:      1,
				Jsonrpc: "2.0",
				Result: LogsResults{
					{
						Address: "0x2cde9919e81b20b4b33dd562a48a84b54c48f00c",
						Topics: []string{
							"0xa6faee2246474597b6de7c76bf9a45d256737543cb0806e6e805b55b38c7663f",
							"0x000000000000000000000000000000000000000000000000000000000000012c"},
						Data:             "0x000000000000000000000000000000000000000000002d6077a3601d1b78000000000000000000000000000000000000000000000000b581c2cc130d1f1800000000000000000000000000000000000000000000000000000000000065692200",
						BlockNumber:      "0x9579fbd",
						TransactionHash:  "0xc37715b1d976e1133f1f62ab3dd856d41be3bb4fe8d4091f3c7849769d041ad3",
						TransactionIndex: fmt.Sprintf("0x%x", 1),
						BlockHash:        "0x3ff6a0c14a272c9379838543735edf677fbe718df12ae52e921fc20f499f6feb",
						LogIndex:         fmt.Sprintf("0x%x", 2),
						Removed:          false,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Eth_getLogs(tt.args.lps)
			if (err != nil) != tt.wantErr {
				t.Errorf("AlchemyClient.eth_getLogs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AlchemyClient.eth_getLogs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlchemyClient_wrong_api_keys_eth_getLogs(t *testing.T) {
	initWrongKeyEnvs()
	type args struct {
		lps []LogsParam
	}
	tests := []struct {
		name    string
		c       *AlchemyClient
		args    args
		want    *AlchemyResponse[LogsResults]
		wantErr bool
	}{
		{
			name: "test empty results",
			c: &AlchemyClient{
				ApiKey:  ALCHEMY_API_KEY_TEST,
				Network: ARB_MAINNET,
				BaseUrlApiV2: BASE_API_URL_V2,
				netClient: &http.Client{
					Timeout: time.Second *10,
				},
			},
			args: args{
				lps: nil,
			},

			want: &AlchemyResponse[LogsResults]{
				Id:      1,
				Jsonrpc: "2.0",
				Result:  nil,
				Error:   ErrorMustBeAuthenticated,
			},
			wantErr: false,
		},
		{
			name: "test with results",
			c: &AlchemyClient{
				ApiKey:  ALCHEMY_API_KEY_TEST,
				Network: ARB_MAINNET,
				BaseUrlApiV2: BASE_API_URL_V2,
				netClient: &http.Client{
					Timeout: time.Second *10,
				},
			},
			args: args{
				lps: []LogsParam{
					{
						BlockHash: "0x3ff6a0c14a272c9379838543735edf677fbe718df12ae52e921fc20f499f6feb",
						Address:   "0x2cde9919e81b20b4b33dd562a48a84b54c48f00c",
					},
				},
			},
			want: &AlchemyResponse[LogsResults]{
				Id:      1,
				Jsonrpc: "2.0",
				Result:  nil,
				Error:   ErrorMustBeAuthenticated,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Eth_getLogs(tt.args.lps)
			if (err != nil) != tt.wantErr {
				t.Errorf("AlchemyClient.eth_getLogs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AlchemyClient.eth_getLogs() = %v, want %v", got, tt.want)
			}
		})
	}
}

