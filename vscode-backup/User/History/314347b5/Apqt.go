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
	URL   string
	Token string
}

type Client struct {
	conf       Config
	httpClient client.Client
}

type ExchangeAccountResponse struct {
	Name                string   `json:"name"`
	Exchange            string   `json:"exchange"`
	Strategies          []string `json:"strategies"`
	Enabled             bool     `json:"enabled"`
	MakerFee            string   `json:"maker_fee"`
	TakerFee            string   `json:"taker_fee"`
	ClientID            string   `json:"client_id"`
	LastUpdatedBalances string   `json:"last_updated_balances"`
}

func New(conf Config, httpClient client.Client) *Client {
	return &Client{
		conf:       conf,
		httpClient: httpClient,
	}
}

func (c *Client) GetExchangeAccounts(ctx context.Context) (*[]ExchangeAccountResponse, error) {
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
			"coinroutesclient: error calling endppoint %s: %s",
			req.URL,
			err,
		)
	}
	defer r.Body.Close()

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
