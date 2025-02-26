// Package signalapi is the interface used for the client that calls the Signal API
package signalapi

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
	"yield-mvp/pkg/client"
)

// Config represents the configuration parameters for the client
type Config struct {
	BTCURL  string
	ETHURL  string
	STUBOUT bool
}

// {
// 	"signal": "null",
// 	"last_data": "2022-05-06 00:55:40",
// 	"current_time": "2022-05-06 00:55:45",
// 	"last_trade": "SELL",
// 	"last_trade_time": "2022-05-05 13:59:04"
// }
type rawSignalResponse struct {
	Signal        string   `json:"signal"`
	LastData      Datetime `json:"last_data"`
	CurrentTime   Datetime `json:"current_time"`
	LastTrade     string   `json:"last_trade"`
	LastTradeTime Datetime `json:"last_trade_time"`
}

type SignalType string

const (
	Short SignalType = "short"
	Sell  SignalType = "sell"
	Null  SignalType = "null"
	Long  SignalType = "long"
	Buy   SignalType = "buy"
)

type SignalResponse struct {
	Chain         string     `json:"chain"`
	Signal        SignalType `json:"signal"`
	LastData      time.Time  `json:"last_data"`
	CurrentTime   time.Time  `json:"current_time"`
	LastTrade     SignalType `json:"last_trade"`
	LastTradeTime time.Time  `json:"last_trade_time"`
}

// Client is the instance variable of the signal api
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

func (c *Client) GetBTCSignal(ctx context.Context) (*SignalResponse, error) {
	// check for stub
	if c.conf.STUBOUT {
		stubTradeTime, err := time.Parse("2006-01-02 15:04:05", "2022-05-20 00:58:28")
		if err != nil {
			return nil, err
		}
		return &SignalResponse{
			Chain:         "btc",
			Signal:        Null,
			LastData:      time.Now(),
			CurrentTime:   time.Now(),
			LastTrade:     Sell,
			LastTradeTime: stubTradeTime,
		}, nil
	}

	url := fmt.Sprintf("%s/test-api", c.conf.BTCURL)
	req, err := client.NewJSONRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error building btc signal request: %w", err)
	}

	var res rawSignalResponse
	if err := c.httpClient.DoJSON(req, &res); err != nil {
		return nil, fmt.Errorf("error making request to btc signal: %w", err)
	}

	signal := &SignalResponse{
		Chain:         "btc",
		Signal:        SignalType(strings.ToLower(res.Signal)),
		LastData:      res.LastData.AsTime(),
		CurrentTime:   res.CurrentTime.AsTime(),
		LastTrade:     SignalType(strings.ToLower(res.LastTrade)),
		LastTradeTime: res.LastTradeTime.AsTime(),
	}

	return signal, nil

}

func (c *Client) GetETHSignal(ctx context.Context) (*SignalResponse, error) {
	// check for stub
	if c.conf.STUBOUT {
		stubTradeTime, err := time.Parse("2006-01-02 15:04:05", "2022-05-20 00:58:28")
		if err != nil {
			return nil, err
		}
		return &SignalResponse{
			Chain:         "eth",
			Signal:        Null,
			LastData:      time.Now(),
			CurrentTime:   time.Now(),
			LastTrade:     Buy,
			LastTradeTime: stubTradeTime,
		}, nil
	}

	url := fmt.Sprintf("%s/test-api", c.conf.ETHURL)
	req, err := client.NewJSONRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error building eth signal request: %w", err)
	}

	var res rawSignalResponse
	if err := c.httpClient.DoJSON(req, &res); err != nil {
		return nil, fmt.Errorf("error making request to eth signal: %w", err)
	}

	signal := &SignalResponse{
		Chain:         "eth",
		Signal:        SignalType(strings.ToLower(res.Signal)),
		LastData:      res.LastData.AsTime(),
		CurrentTime:   res.CurrentTime.AsTime(),
		LastTrade:     SignalType(strings.ToLower(res.LastTrade)),
		LastTradeTime: res.LastTradeTime.AsTime(),
	}

	return signal, nil
}
