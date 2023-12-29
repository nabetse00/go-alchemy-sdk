package goalchemysdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/avast/retry-go/v4"

	// "strings"
	//"io/ioutil"
)

const BASE_API_URL_V2 = ".g.alchemy.com/v2"
const MAX_RETRY_DEFAULT = 3
const DELAY_DEFAULT = 1

type AlchemyClient struct {
	ApiKey       string
	Network      Network
	MaxRetry     uint
	Delay        uint
	BaseUrlApiV2 string
	netClient  *http.Client
}

type AlchemyClientError struct {
	method string
	message string
}

func (e *AlchemyClientError) Error() string {
	return fmt.Sprintf("Alchemy client error on method [%s] reason %s", e.method, e.message)
}

func (c *AlchemyClient) Init(apiKey string, network Network, maxRetry uint, delay uint, baseUrlApiV2 string, timeout time.Duration) error {
	c.ApiKey = apiKey
	c.Network = network
	c.MaxRetry = maxRetry
	c.Delay = delay
	c.BaseUrlApiV2 = baseUrlApiV2
	if c.ApiKey == "" {
		return &AlchemyClientError{"Init", "Empty Alchemy key"}
	}
	if c.Network == "" {
		return &AlchemyClientError{"Init", "Empty Alchemy Network"}
	}
	if c.BaseUrlApiV2 == "" {
		c.BaseUrlApiV2 = BASE_API_URL_V2
	}
	if c.MaxRetry == 0 {
		c.MaxRetry = MAX_RETRY_DEFAULT
	}
	if delay == 0 {
		c.Delay = DELAY_DEFAULT
	}
	c.netClient = &http.Client{
		Timeout: timeout,
	  }
	return nil
}

func (c *AlchemyClient) getApiUrl() (string,error){
	if c.ApiKey == "" {
		return "", &AlchemyClientError{"getApiUrl()","Empty Alchemy key" }
	}
	if c.Network == "" {
		return c.BaseUrlApiV2, &AlchemyClientError{"getApiUrl()","Empty Alchemy Network" }
	}
	if c.BaseUrlApiV2 == ""  {
		c.BaseUrlApiV2 = BASE_API_URL_V2
	}
	if c.Delay == 0  {
		c.Delay = DELAY_DEFAULT
	}
	if c.MaxRetry == 0 {
		c.MaxRetry = MAX_RETRY_DEFAULT
	}
	return "https://" + string(c.Network) + c.BaseUrlApiV2 + "/" + c.ApiKey, nil
}

var _ error = (*RetriableError)(nil)

func executePost[P any, R any](client *AlchemyClient, jsonP JsonParams[P]) (*AlchemyResponse[R], error) {
	url, _ := client.getApiUrl()
	body, err := json.Marshal(jsonP)
	if err != nil {
		method := fmt.Sprintf("executePost - %s", jsonP.Method)
		return &AlchemyResponse[R]{}, &AlchemyClientError{method, err.Error()}
	}

	var resp *http.Response
	var data AlchemyResponse[R]

	data, err = retry.DoWithData(
		func() (AlchemyResponse[R], error) {
			resp, err = client.netClient.Post(
				url,
				"application/json",
				bytes.NewBuffer(body))

			if err == nil {
				defer func() {
					if err := resp.Body.Close(); err != nil {
						panic(err)
					}
				}()

				if resp.StatusCode != http.StatusOK {
					err = fmt.Errorf("HTTP %d for: %s", resp.StatusCode, string(body))
					if resp.StatusCode == http.StatusTooManyRequests {
						// fmt.Printf("inside too many %s -- %s \n", resp.Header.Get("Retry-After"), resp.Header.Get("retryAfter"))
						// check Retry-After header if it contains seconds to wait for the next retry
						if retryAfter, e := strconv.ParseInt(resp.Header.Get("Retry-After"), 10, 32); e == nil {
							// the server returns 0 to inform that the operation cannot be retried
							if retryAfter <= 0 {
								return AlchemyResponse[R]{}, retry.Unrecoverable(err)
							}
							return AlchemyResponse[R]{}, &RetriableError{
								Err:        err,
								RetryAfter: time.Duration(retryAfter) * time.Second,
							}
						}
						if retryAfter, e := strconv.ParseInt(resp.Header.Get("retryAfter"), 10, 32); e == nil {
							// the server returns 0 to inform that the operation cannot be retried
							if retryAfter <= 0 {
								return AlchemyResponse[R]{}, retry.Unrecoverable(err)
							}
							return AlchemyResponse[R]{}, &RetriableError{
								Err:        err,
								RetryAfter: time.Duration(retryAfter) * time.Second,
							}
						}
						return AlchemyResponse[R]{},err
					}
				}
				err = json.NewDecoder(resp.Body).Decode(&data)
				// warning ReadAll consumes resp.Body !
				// read_body, _ := ioutil.ReadAll(resp.Body)
				// fmt.Printf("response body is: %#v\n", data)
				// fmt.Printf("error on decode is %s\n", err)
				return data, err
			}

			return AlchemyResponse[R]{}, err
		},
		retry.Attempts(client.MaxRetry),
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			fmt.Printf("[%d]Server fails with: %s\n", n, err.Error())
			if retriable, ok := err.(*RetriableError); ok {
				fmt.Printf("Client follows server recommendation to retry after %v\n", retriable.RetryAfter)
				return retriable.RetryAfter
			}
			// apply a default exponential back off strategy
			return retry.BackOffDelay(n, err, config)
		}),
	)

	return &data, err
}


