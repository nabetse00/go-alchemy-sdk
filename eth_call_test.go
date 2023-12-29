package goalchemysdk

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestAlchemyClient_Eth_call(t *testing.T) {
	initEnvs()
	type args struct {
		txn CallTxn
		blk CallBlk
	}
	tests := []struct {
		name    string
		c       *AlchemyClient
		args    args
		want    *AlchemyResponse[CallResult]
		wantErr bool
	}{
		{
			name: "test raw eth call - address result",
			c: &AlchemyClient{
				ApiKey:  ALCHEMY_API_KEY_TEST,
				Network: ETH_MAINNET,
				netClient: &http.Client{
					Timeout: time.Second * 10,
				},
			},
			args: args{
				txn: CallTxn{
					//From:             "",
					To:   "0x4976fb03C32e5B8cfe2b6cCB31c09Ba78EBaBa41",
					Data: "0x3b3b57debf074faa138b72c65adbdcfb329847e4f2c04bde7f7dd7fcad5a52d2f395a558",
					// addr(0xbf074faa138b72c65adbdcfb329847e4f2c04bde7f7dd7fcad5a52d2f395a558)
					// 0x5555763613a12D8F3e73be831DFf8598089d3dCa => ricmoo.eth
					// Gas:              "0x00",
					// GasPrice:         "0x09184e72a000",
					// Value:            "0x00",
					// BlockHash:        "",
					//RequireCanonical: false,
				},
				blk: LATEST,
			},
			want: &AlchemyResponse[CallResult]{
				Id:      1,
				Jsonrpc: "2.0",
				Result:  "0x0000000000000000000000005555763613a12d8f3e73be831dff8598089d3dca",
			},
			wantErr: false,
		},
		{
			name: "test raw eth call - not result",
			c: &AlchemyClient{
				ApiKey:  ALCHEMY_API_KEY_TEST,
				Network: ARB_MAINNET, 
				netClient: &http.Client{
					Timeout: time.Second * 10,
				},
			},
			args: args{
				txn: CallTxn{
					//From:             "",
					To:   "0xdeadbeefBd5D07dd0CeCc66161FC93D7c9000da1",
					Data: "0x70a082310000000000000000000000006E0d01A76C3Cf4288372a29124A26D4353EE51BE",
					// balanceOf()
					// Gas:              "0x00",
					// GasPrice:         "0x09184e72a000",
					// Value:            "0x00",
					// BlockHash:        "",
					//RequireCanonical: false,
				},
				blk: LATEST,
			},
			want: &AlchemyResponse[CallResult]{
				Id:      1,
				Jsonrpc: "2.0",
				Result:  "0x",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Eth_call(tt.args.txn, tt.args.blk)
			if (err != nil) != tt.wantErr {
				t.Errorf("AlchemyClient.Eth_call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AlchemyClient.Eth_call() = %v, want %v", got, tt.want)
			}
		})
	}
}
