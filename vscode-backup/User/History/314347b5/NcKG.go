package coinroutesapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"yield-mvp/pkg/client"
)

type Config struct {
	URL   string
	Token string
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
	if err := c.makeRequest(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil

}

func (c *Client) CreateClientOrders(
	ctx context.Context,
	order *ClientOrderCreateRequest) (*ClientOrderCreateResponse, error) {
	// check if quantity has not been set or zero value
	if orderFloat, err := strconv.ParseFloat(order.Quantity, 32); err != nil {
		return nil, fmt.Errorf("cannot parse quantity, aborting: %w", err)
	}
	if orderFloat == 0 {
		return nil, fmt.Errorf(
			"quantity is zero value, aborting request with order: %+v",
			order,
		)
	}
	// https://sor.yourcompany.com/api/client_orders/
	url := fmt.Sprintf("%s/api/client_orders/", c.conf.URL)
	req, err := c.buildRequest(ctx, http.MethodPost, url, order)
	if err != nil {
		return nil, err
	}

	var resp ClientOrderCreateResponse
	if err := c.makeRequest(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) GetPositions(ctx context.Context, str StrategyType) (*[]PositionResponse, error) {
	url := fmt.Sprintf("%s/api/positions?strategy=%s", c.conf.URL, str)
	req, err := c.buildRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var resp []PositionResponse
	if err := c.makeRequest(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) GetBalances(ctx context.Context, str StrategyType) (*[]CurrencyBalanceResponse, error) {
	url := fmt.Sprintf("%s/api/currency_balances?strategy=%s", c.conf.URL, str)
	req, err := c.buildRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var resp []CurrencyBalanceResponse
	if err := c.makeRequest(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
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

func (c *Client) makeRequest(req *http.Request, res interface{}) error {
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf(
			"coinroutesclient: error parsing body in response %s: %s",
			req.URL,
			err,
		)
	}

	json.Unmarshal(body, &res)

	return nil
}
