package goalchemysdk

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestAlchemyClient_Eth_getStorageAt(t *testing.T) {
	initEnvs()
	type args struct {
		address  string
		id       string
		blocktag BlockTag
	}
	tests := []struct {
		name    string
		c       *AlchemyClient
		args    args
		want    *AlchemyResponse[GetCodeResult]
		wantErr bool
	}{
		{
			name: "test valid contract",
			c: &AlchemyClient{
				ApiKey:  ALCHEMY_API_KEY_TEST,
				Network: ARB_MAINNET,
				netClient: &http.Client{
					Timeout: time.Second *10,
				},
			},
			args: args{
				address:  "0x912CE59144191C1204E64559FE8253a0e49E6548",
				id: "0x0",
				blocktag: LATEST,
			},
			want: &AlchemyResponse[GetCodeResult]{
				Id:      1,
				Jsonrpc: "2.0",
				Result:  "0x0000000000000000000000000000000000000000000000000000000000000001",
			},
			wantErr: false,
		},
		{
			name: "test invalid address too short",
			c: &AlchemyClient{
				ApiKey:  ALCHEMY_API_KEY_TEST,
				Network: ARB_MAINNET,
				netClient: &http.Client{
					Timeout: time.Second *10,
				},
			},
			args: args{
				address:  "0xdeadbeef",
				id: "0x0",
				blocktag: LATEST,
			},
			want: &AlchemyResponse[GetCodeResult]{
				Id:      1,
				Jsonrpc: "2.0",
				Error:   ErrorTooShortAddress,
			},
			wantErr: false,
		},
		{
			name: "test invalid address too long",
			c: &AlchemyClient{
				ApiKey:  ALCHEMY_API_KEY_TEST,
				Network: ARB_MAINNET,
				netClient: &http.Client{
					Timeout: time.Second *10,
				},
			},
			args: args{
				address:  "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
				id: "0x0",
				blocktag: LATEST,
			},
			want: &AlchemyResponse[GetCodeResult]{
				Id:      1,
				Jsonrpc: "2.0",
				Error:   ErrorTooLongAddress,
			},
			wantErr: false,
		},
		{
			name: "test invalid address",
			c: &AlchemyClient{
				ApiKey:  ALCHEMY_API_KEY_TEST,
				Network: ARB_MAINNET,
				netClient: &http.Client{
					Timeout: time.Second *10,
				},
			},
			args: args{
				address:  "0xXaedbeefbeef067E90D5Cd1F8052B83562Ae670bA4A211a8",
				id: "0x0",
				blocktag: LATEST,
			},
			want: &AlchemyResponse[GetCodeResult]{
				Id:      1,
				Jsonrpc: "2.0",
				Error:   ErrorInvalidAddress,
			},
			wantErr: false,
		},
		{
			name: "test not a contract",
			c: &AlchemyClient{
				ApiKey:  ALCHEMY_API_KEY_TEST,
				Network: ARB_MAINNET,
				netClient: &http.Client{
					Timeout: time.Second *10,
				},
			},
			args: args{
				address:  "0xdeadbeedeadbeefdeadbeefdeadbeefdeadbeefd",
				id: "0x0",
				blocktag: LATEST,
			},
			want: &AlchemyResponse[GetCodeResult]{
				Id:      1,
				Jsonrpc: "2.0",
				Result:  "0x0000000000000000000000000000000000000000000000000000000000000000",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Eth_getStorageAt(tt.args.address, tt.args.id, tt.args.blocktag)
			if (err != nil) != tt.wantErr {
				t.Errorf("AlchemyClient.Eth_getStorageAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AlchemyClient.Eth_getStorageAt() = %v, want %v", got, tt.want)
			}
		})
	}
}
