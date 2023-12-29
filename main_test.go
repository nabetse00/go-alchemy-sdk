package goalchemysdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

var (
	ALCHEMY_API_KEY_TEST string
)

// Helpers
func initEnvs() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %s.\n Are env vars setup manualy?\n", err)
	}
	ALCHEMY_API_KEY_TEST = os.Getenv("ALCHEMY_API_KEY")
	if ALCHEMY_API_KEY_TEST == "" {
		panic("Init env test: ALCHEMY_API_KEY_TEST empty")
	}
	os.Setenv("APP_ENV", "test")
}

func initWrongKeyEnvs() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %s.\n Are env vars setup manualy?\n", err)
	}
	ALCHEMY_API_KEY_TEST = "123456"
	os.Setenv("APP_ENV", "test")
}

func fakeRetryServerRecoverable() *httptest.Server {
	const MAX_RETRY = 3
	var retryCount = 0
	responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if retryCount < MAX_RETRY {
			time.Sleep(20 * time.Millisecond)
			w.WriteHeader(http.StatusTooManyRequests)
			retryCount++
			return
		}

		if retryCount < MAX_RETRY+1 {
			time.Sleep(20 * time.Millisecond)
			w.Header().Add("retryAfter", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			retryCount++
			return
		}

		if retryCount < MAX_RETRY+2 {
			time.Sleep(20 * time.Millisecond)
			w.Header().Add("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			retryCount++
			return
		}

		resp := &AlchemyResponse[interface{}]{
			Id:      1,
			Jsonrpc: "2.0",
			Result:  "fake",
		}

		var result, _ = json.Marshal(resp)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result))

	})
	ts := httptest.NewServer(responseHandler)
	return ts
}

func fakeRetryServerUnRecoverable1() *httptest.Server {
	const MAX_RETRY = 3
	var retryCount = 0
	responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if retryCount < MAX_RETRY {
			time.Sleep(20 * time.Millisecond)
			w.WriteHeader(http.StatusTooManyRequests)
			retryCount++
			return
		}

		if retryCount < MAX_RETRY+1 {
			time.Sleep(20 * time.Millisecond)
			w.Header().Add("retryAfter", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			retryCount++
			return
		}

		if retryCount < MAX_RETRY+2 {
			time.Sleep(20 * time.Millisecond)
			w.Header().Add("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			retryCount++
			return
		}

		resp := &AlchemyResponse[interface{}]{
			Id:      1,
			Jsonrpc: "2.0",
			Result:  "fake",
		}

		var result, _ = json.Marshal(resp)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result))

	})
	ts := httptest.NewServer(responseHandler)
	return ts
}

func fakeRetryServerUnRecoverable2() *httptest.Server {
	const MAX_RETRY = 3
	var retryCount = 0
	responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if retryCount < MAX_RETRY {
			time.Sleep(20 * time.Millisecond)
			w.WriteHeader(http.StatusTooManyRequests)
			retryCount++
			return
		}

		if retryCount < MAX_RETRY+1 {
			time.Sleep(20 * time.Millisecond)
			w.Header().Add("retryAfter", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			retryCount++
			return
		}

		if retryCount < MAX_RETRY+2 {
			time.Sleep(20 * time.Millisecond)
			w.Header().Add("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			retryCount++
			return
		}

		resp := &AlchemyResponse[interface{}]{
			Id:      1,
			Jsonrpc: "2.0",
			Result:  "fake",
		}

		var result, _ = json.Marshal(resp)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result))

	})
	ts := httptest.NewServer(responseHandler)
	return ts
}

// tests

func TestAlchemyClient_getApiUrl(t *testing.T) {
	initEnvs()
	tests := []struct {
		name    string
		c       *AlchemyClient
		want    string
		wantErr bool
	}{
		{
			name:    "Error on empty key",
			c:       &AlchemyClient{ApiKey: "", Network: "whatever"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Error on empty network",
			c:       &AlchemyClient{ApiKey: "key", Network: ""},
			want:    "",
			wantErr: true,
		},
		{
			name:    "test defaults",
			c:       &AlchemyClient{ApiKey: "key", Network: "whatever"},
			want:    fmt.Sprintf("https://%s%s/key", "whatever", BASE_API_URL_V2),
			wantErr: false,
		},
		{
			name:    "test network arbritum main",
			c:       &AlchemyClient{ApiKey: "key", Network: ARB_MAINNET},
			want:    fmt.Sprintf("https://%s%s/key", ARB_MAINNET, BASE_API_URL_V2),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.getApiUrl()
			if tt.wantErr && err == nil {
				t.Errorf("AlchemyClient.getApiUrl() wants an error and none emited")
			}
			if !tt.wantErr && (got != tt.want || err != nil) {
				t.Errorf("AlchemyClient.getApiUrl() = %v, want %v, error is: %s", got, tt.want, err)
			}
		})
	}
}

func TestAlchemyClient_executePost(t *testing.T) {
	initEnvs()
	type args struct {
		j JsonParams[LogsParam]
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
					Timeout: time.Second * 10,
				},
			},
			args: args{
				j: JsonParams[LogsParam]{
					Id:      1,
					Jsonrpc: "2.0",
					Method:  "eth_getLogs",
					Params:  nil,
				},
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
					Timeout: time.Second * 10,
				},
			},
			args: args{
				j: JsonParams[LogsParam]{
					Id:      1,
					Jsonrpc: "2.0",
					Method:  "eth_getLogs",
					Params: []LogsParam{
						{
							BlockHash: "0x3ff6a0c14a272c9379838543735edf677fbe718df12ae52e921fc20f499f6feb",
							Address:   "0x2cde9919e81b20b4b33dd562a48a84b54c48f00c",
						},
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
			got, err := executePost[LogsParam, LogsResults](tt.c, tt.args.j)
			if (err != nil) != tt.wantErr {
				t.Errorf("AlchemyClient.executePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, *tt.want) {
				// fmt.Printf("=====> %s -- %s \n", reflect.TypeOf(got.Jsonrpc), reflect.TypeOf(tt.want.Jsonrpc))
				t.Errorf("AlchemyClient.executePost() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestAlchemyClient_executePost_wrongJson(t *testing.T) {
	initEnvs()
	wrong := make([]interface{}, 1)
	wrong[0] = make(chan int)
	type args struct {
		j JsonParams[interface{}]
	}
	tests := []struct {
		name    string
		c       *AlchemyClient
		args    args
		want    *AlchemyResponse[interface{}]
		wantErr bool
	}{
		{
			name: "test empty results",
			c: &AlchemyClient{
				ApiKey:  ALCHEMY_API_KEY_TEST,
				Network: ARB_MAINNET,
				netClient: &http.Client{
					Timeout: time.Second * 10,
				},
			},
			args: args{
				j: JsonParams[interface{}]{
					Id:      1,
					Jsonrpc: "2.0",
					Method:  "eth_getLogs",
					Params:  wrong,
				},
			},
			want: &AlchemyResponse[interface{}]{
				Id:      1,
				Jsonrpc: "2.0",
				Result:  nil,
				Error:   ErrorExpectedAtLeastOneArgument,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := executePost[interface{}, interface{}](tt.c, tt.args.j)
			if (err != nil) != tt.wantErr {
				t.Errorf("AlchemyClient.executePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(*got, *tt.want) {
				// fmt.Printf("=====> %s -- %s \n", reflect.TypeOf(got.Jsonrpc), reflect.TypeOf(tt.want.Jsonrpc))
				t.Errorf("AlchemyClient.executePost() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestAlchemyClient_executePost_Retry_Recoverable(t *testing.T) {
	initEnvs()
	ts := fakeRetryServerRecoverable()
	defer ts.Close()

	empty := make([]interface{}, 0)
	type args struct {
		j JsonParams[interface{}]
	}
	tests := []struct {
		name    string
		c       *AlchemyClient
		args    args
		want    *AlchemyResponse[interface{}]
		wantErr bool
	}{
		{
			name: "test results after 3 retrys",
			c: &AlchemyClient{
				ApiKey:       ALCHEMY_API_KEY_TEST,
				Network:      "",
				BaseUrlApiV2: ts.URL,
				MaxRetry:     6,
				Delay:        1,
				netClient: &http.Client{
					Timeout: time.Second * 10,
				},
			},
			args: args{
				j: JsonParams[interface{}]{
					Id:      1,
					Jsonrpc: "2.0",
					Method:  "eth_fake",
					Params:  empty,
				},
			},
			want: &AlchemyResponse[interface{}]{
				Id:      1,
				Jsonrpc: "2.0",
				Result:  "fake",
				//Error:   ErrorExpectedAtLeastOneArgument,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := executePost[interface{}, interface{}](tt.c, tt.args.j)
			if (err != nil) != tt.wantErr {
				t.Errorf("AlchemyClient.executePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(*got, *tt.want) {
				t.Errorf("AlchemyClient.executePost() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestAlchemyClient_executePost_Retry_UnrecovarableAfter1(t *testing.T) {
	initEnvs()
	ts := fakeRetryServerUnRecoverable1()
	defer ts.Close()

	empty := make([]interface{}, 0)
	type args struct {
		j JsonParams[interface{}]
	}
	tests := []struct {
		name    string
		c       *AlchemyClient
		args    args
		want    *AlchemyResponse[interface{}]
		wantErr bool
	}{
		{
			name: "test results after 3 retrys",
			c: &AlchemyClient{
				ApiKey:       ALCHEMY_API_KEY_TEST,
				Network:      "",
				BaseUrlApiV2: ts.URL,
				MaxRetry:     6,
				Delay:        1,
				netClient: &http.Client{
					Timeout: time.Second * 10,
				},
			},
			args: args{
				j: JsonParams[interface{}]{
					Id:      1,
					Jsonrpc: "2.0",
					Method:  "eth_fake",
					Params:  empty,
				},
			},
			want: &AlchemyResponse[interface{}]{
				Id:      1,
				Jsonrpc: "2.0",
				Result:  "fake",
				//Error:   ErrorExpectedAtLeastOneArgument,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := executePost[interface{}, interface{}](tt.c, tt.args.j)
			if (err != nil) != tt.wantErr {
				t.Errorf("AlchemyClient.executePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(*got, *tt.want) {
				t.Errorf("AlchemyClient.executePost() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestAlchemyClient_executePost_Retry_UnrecovarableAfter2(t *testing.T) {
	initEnvs()
	ts := fakeRetryServerUnRecoverable2()
	defer ts.Close()

	empty := make([]interface{}, 0)
	type args struct {
		j JsonParams[interface{}]
	}
	tests := []struct {
		name    string
		c       *AlchemyClient
		args    args
		want    *AlchemyResponse[interface{}]
		wantErr bool
	}{
		{
			name: "test results after 3 retrys",
			c: &AlchemyClient{
				ApiKey:       ALCHEMY_API_KEY_TEST,
				Network:      "",
				BaseUrlApiV2: ts.URL,
				MaxRetry:     6,
				Delay:        1,
				netClient: &http.Client{
					Timeout: time.Second * 10,
				},
			},
			args: args{
				j: JsonParams[interface{}]{
					Id:      1,
					Jsonrpc: "2.0",
					Method:  "eth_fake",
					Params:  empty,
				},
			},
			want: &AlchemyResponse[interface{}]{
				Id:      1,
				Jsonrpc: "2.0",
				Result:  "fake",
				//Error:   ErrorExpectedAtLeastOneArgument,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := executePost[interface{}, interface{}](tt.c, tt.args.j)
			if (err != nil) != tt.wantErr {
				t.Errorf("AlchemyClient.executePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(*got, *tt.want) {
				t.Errorf("AlchemyClient.executePost() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestAlchemyClient_executePost_wrong_method(t *testing.T) {
	initEnvs()
	type args struct {
		j JsonParams[LogsParam]
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
					Timeout: time.Second * 10,
				},
			},
			args: args{
				j: JsonParams[LogsParam]{
					Id:      1,
					Jsonrpc: "2.0",
					Method:  "eth_wrong_method",
					Params:  nil,
				},
			},
			want: &AlchemyResponse[LogsResults]{
				Id:      1,
				Jsonrpc: "2.0",
				Result:  nil,
				Error:   ErrorWrongMethod("eth_wrong_method"),
			},
			wantErr: false,
		},
		{
			name: "test with results",
			c: &AlchemyClient{
				ApiKey:  ALCHEMY_API_KEY_TEST,
				Network: ARB_MAINNET,
				netClient: &http.Client{
					Timeout: time.Second * 10,
				},
			},
			args: args{
				j: JsonParams[LogsParam]{
					Id:      1,
					Jsonrpc: "2.0",
					Method:  "eth_wrong_method",
					Params: []LogsParam{
						{
							BlockHash: "0x3ff6a0c14a272c9379838543735edf677fbe718df12ae52e921fc20f499f6feb",
							Address:   "0x2cde9919e81b20b4b33dd562a48a84b54c48f00c",
						},
					},
				},
			},
			want: &AlchemyResponse[LogsResults]{
				Id:      1,
				Jsonrpc: "2.0",
				Result:  nil,
				Error:   ErrorWrongMethod("eth_wrong_method"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := executePost[LogsParam, LogsResults](tt.c, tt.args.j)
			if (err != nil) != tt.wantErr {
				t.Errorf("AlchemyClient.executePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, *tt.want) {
				// fmt.Printf("=====> %s -- %s \n", reflect.TypeOf(got.Jsonrpc), reflect.TypeOf(tt.want.Jsonrpc))
				t.Errorf("AlchemyClient.executePost() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestAlchemyClient_executePost_wrong_url(t *testing.T) {
	initEnvs()
	type args struct {
		j JsonParams[LogsParam]
	}
	tests := []struct {
		name    string
		c       *AlchemyClient
		args    args
		want    *AlchemyResponse[LogsResults]
		wantErr bool
	}{
		{
			name: "test wrong api url",
			c: &AlchemyClient{
				ApiKey:       ALCHEMY_API_KEY_TEST,
				Network:      ARB_MAINNET,
				BaseUrlApiV2: "wrong" + BASE_API_URL_V2,
				netClient: &http.Client{
					Timeout: time.Second * 10,
				},
			},
			args: args{
				j: JsonParams[LogsParam]{
					Id:      1,
					Jsonrpc: "2.0",
					Method:  "eth_wrong_method",
					Params:  nil,
				},
			},
			want: &AlchemyResponse[LogsResults]{
				Id:      0,
				Jsonrpc: "",
				Result:  nil,
			},
			wantErr: true,
		},
		{
			name: "test with wrong method",
			c: &AlchemyClient{
				ApiKey:  ALCHEMY_API_KEY_TEST,
				Network: ARB_MAINNET,
				netClient: &http.Client{
					Timeout: time.Second * 10,
				},
			},
			args: args{
				j: JsonParams[LogsParam]{
					Id:      1,
					Jsonrpc: "2.0",
					Method:  "eth_wrong_method",
					Params: []LogsParam{
						{
							BlockHash: "0x3ff6a0c14a272c9379838543735edf677fbe718df12ae52e921fc20f499f6feb",
							Address:   "0x2cde9919e81b20b4b33dd562a48a84b54c48f00c",
						},
					},
				},
			},
			want: &AlchemyResponse[LogsResults]{
				Id:      1,
				Jsonrpc: "2.0",
				Result:  nil,
				Error:   ErrorWrongMethod("eth_wrong_method"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := executePost[LogsParam, LogsResults](tt.c, tt.args.j)
			if (err != nil) != tt.wantErr {
				t.Errorf("AlchemyClient.executePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, *tt.want) {
				// fmt.Printf("=====> %s -- %s \n", reflect.TypeOf(got.Jsonrpc), reflect.TypeOf(tt.want.Jsonrpc))
				t.Errorf("AlchemyClient.executePost() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestAlchemyClient_Init(t *testing.T) {
	initEnvs()
	type args struct {
		apiKey       string
		network      Network
		maxRetry     uint
		delay        uint
		baseUrlApiV2 string
	}
	tests := []struct {
		name    string
		c       *AlchemyClient
		args    args
		wantErr bool
		wants   *AlchemyClient
	}{
		{
			name: "Incorrect Init: empty key",
			c:    &AlchemyClient{},
			args: args{
				apiKey:       "",
				network:      "",
				maxRetry:     0,
				delay:        0,
				baseUrlApiV2: "",
			},
			wants:   nil,
			wantErr: true,
		},
		{
			name: "Incorrect Init: empty network",
			c:    &AlchemyClient{},
			args: args{
				apiKey:       ALCHEMY_API_KEY_TEST,
				network:      "",
				maxRetry:     0,
				delay:        0,
				baseUrlApiV2: "",
			},
			wants:   nil,
			wantErr: true,
		},
		{
			name: "Correct Init default values",
			c:    &AlchemyClient{},
			args: args{
				apiKey:       ALCHEMY_API_KEY_TEST,
				network:      ARB_MAINNET,
				maxRetry:     0,
				delay:        0,
				baseUrlApiV2: "",
			},
			wants: &AlchemyClient{
				ApiKey:       ALCHEMY_API_KEY_TEST,
				Network:      ARB_MAINNET,
				MaxRetry:     MAX_RETRY_DEFAULT,
				Delay:        DELAY_DEFAULT,
				BaseUrlApiV2: BASE_API_URL_V2,
				netClient: &http.Client{
					Timeout: time.Second * 10,
				},
			},
			wantErr: false,
		},
		{
			name: "Correct Init custom values values",
			c:    &AlchemyClient{},
			args: args{
				apiKey:       ALCHEMY_API_KEY_TEST,
				network:      ARB_MAINNET,
				maxRetry:     10,
				delay:        2,
				baseUrlApiV2: "someurl.here.com",
			},
			wants: &AlchemyClient{
				ApiKey:       ALCHEMY_API_KEY_TEST,
				Network:      ARB_MAINNET,
				MaxRetry:     10,
				Delay:        2,
				BaseUrlApiV2: "someurl.here.com",
				netClient: &http.Client{
					Timeout: time.Second * 10,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Init(tt.args.apiKey, tt.args.network, tt.args.maxRetry, tt.args.delay, tt.args.baseUrlApiV2, time.Second*10); (err != nil) != tt.wantErr {
				t.Errorf("AlchemyClient.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.wants, tt.c) {
				t.Errorf("AlchemyClient.Init() error got = %v, wants %v", tt.c, tt.wants)
			}
		})
	}
}

func TestAlchemyClientError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *AlchemyClientError
		want string
	}{
		{
			name: "Test client error",
			e:    &AlchemyClientError{"method", "message"},
			want: "Alchemy client error on method [method] reason message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("AlchemyClientError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlchemyClient_Close(t *testing.T) {
	tests := []struct {
		name string
		c    *AlchemyClient
	}{
		{
			name: "test client close",
			c: &AlchemyClient{
				ApiKey:       ALCHEMY_API_KEY_TEST,
				Network:      ARB_MAINNET,
				MaxRetry:     10,
				Delay:        2,
				BaseUrlApiV2: "someurl.here.com",
				netClient: &http.Client{
					Timeout: time.Second * 10,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Close()
		})
	}
}
