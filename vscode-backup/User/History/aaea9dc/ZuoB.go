package coinroutesapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"yield-mvp/pkg/client"
)

type Config struct {
	URL      string
	Token    string
	Simulate bool // will not send client order requests only
}

func (c *Config) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("missing required URL config property")
	}
	if c.Token == "" {
		return fmt.Errorf("missing required Token config property")
	}

	return nil
}

type Client struct {
	conf       Config
	httpClient client.Client
}

func New(conf Config, httpClient client.Client) *Client {
	return &Client{
		conf:       conf,
		httpClient: httpClient,
	}
}

// List all Exchange account that the current user has access to. Exchange Accounts are permissioned by strategy.
// http://downloads.coinroutes.com/docs/api/spec/3.14.14/index.html#operation/api_exchange_accounts_list
func (c *Client) GetExchangeAccounts(ctx context.Context) (
	*[]ExchangeAccountResponse,
	error,
) {
	// https://sor.yourcompany.com/api/exchange_accounts/
	url := fmt.Sprintf("%s/api/exchange_accounts", c.conf.URL)
	req, err := c.buildRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var resp []ExchangeAccountResponse
	if err := c.makeRequest(req, &resp, false); err != nil {
		return nil, err
	}

	return &resp, nil

}

// Creates a new client order and automatically creates and executes appropriate SOR orders.
// http://downloads.coinroutes.com/docs/api/spec/3.14.14/index.html#operation/api_client_orders_create
func (c *Client) CreateClientOrders(
	ctx context.Context,
	order *ClientOrderCreateRequest) (*ClientOrderCreateResponse, error) {
	// https://sor.yourcompany.com/api/client_orders/
	url := fmt.Sprintf("%s/api/client_orders/", c.conf.URL)
	req, err := c.buildRequest(ctx, http.MethodPost, url, order)
	if err != nil {
		return nil, err
	}

	var resp ClientOrderCreateResponse

	if err := c.makeRequest(req, &resp, c.conf.Simulate); err != nil {
		return nil, err
	}

	return &resp, nil
}

// List all Exchange Account Positions
// http://downloads.coinroutes.com/docs/api/spec/3.14.14/index.html#operation/api_positions_list
func (c *Client) GetPositions(ctx context.Context, str string) (*[]PositionResponse, error) {
	url := fmt.Sprintf("%s/api/positions?strategy=%s", c.conf.URL, str)
	req, err := c.buildRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var resp []PositionResponse
	if err := c.makeRequest(req, &resp, false); err != nil {
		return nil, err
	}

	return &resp, nil
}

// List all Exchange Currency Balances.
// http://downloads.coinroutes.com/docs/api/spec/3.14.14/index.html#operation/api_currency_balances_list
func (c *Client) GetBalances(ctx context.Context, str string) (*[]CurrencyBalanceResponse, error) {
	url := fmt.Sprintf("%s/api/currency_balances?strategy=%s", c.conf.URL, str)
	req, err := c.buildRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var resp []CurrencyBalanceResponse
	if err := c.makeRequest(req, &resp, false); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Get details for a client order
// http://downloads.coinroutes.com/docs/api/spec/3.14.14/index.html#operation/api_client_orders_read
func (c *Client) GetClientOrder(
	ctx context.Context,
	id string,
) (*ClientOrderGetResponse, error) {
	url := fmt.Sprintf("%s/api/client_orders/%s", c.conf.URL, id)
	req, err := c.buildRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var resp ClientOrderGetResponse
	if err := c.makeRequest(req, &resp, false); err != nil {
		return nil, err
	}

	return &resp, nil
}

// IsSimulate returns true if the client is configured to simulate
func (c *Client) IsSimulated() bool {
	return c.conf.Simulate
}

func (c *Client) buildRequest(
	ctx context.Context,
	method, url string,
	body interface{},
) (*http.Request, error) {
	req, err := client.NewJSONRequest(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf(
			"coinroutesclient: error building request %s: %s",
			url,
			err,
		)
	}

	req.Header.Add("Authorization", "Token "+c.conf.Token)

	return req, nil
}

func (c *Client) makeRequest(req *http.Request, res interface{}, sim bool) error {
	var body []byte
	if !sim {
		r, err := c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf(
				"coinroutesclient: error calling endpoint %s: %s",
				req.URL,
				err,
			)
		}
		defer r.Body.Close()

		// dump, err := httputil.DumpResponse(r, true)
		// if err != nil {
		// 	return fmt.Errorf("error dumping response: %s", err)
		// }

		// fmt.Printf("url: %s", req.URL)
		// fmt.Printf("response dump: %+s", dump)

		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf(
				"coinroutesclient: error parsing body in response %s: %s",
				req.URL,
				err,
			)
		}
	} else {
		fmt.Println("DEBUG: SIMULATING RESPONSE FROM COINROUTES")
		body = sampleOrderCreateResponse
	}

	json.Unmarshal(body, &res)

	return nil
}
