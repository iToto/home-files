// Package signalapi is the interface used for the client that calls the Signal API
package signalapi

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
	"yield/signal-logger/internal/entities"
	"yield/signal-logger/internal/wlog"
	"yield/signal-logger/pkg/client"
)

// {
// 	"signal": "null",
// 	"last_data": "2022-05-06 00:55:40",
// 	"current_time": "2022-05-06 00:55:45",
// 	"last_trade": "SELL",
// 	"last_trade_time": "2022-05-05 13:59:04"
// }
type rawSignalResponseV1 struct {
	Signal        string   `json:"signal"`
	LastData      Datetime `json:"last_data"`
	CurrentTime   Datetime `json:"current_time"`
	LastTrade     string   `json:"last_trade"`
	LastTradeTime Datetime `json:"last_trade_time"`
}
type rawSignalResponseV2 struct {
	FetchResultStatus   string `json:"fetch_result_status"`
	FetchType           string `json:"fetch_type"`
	StrategyState       string `json:"strategy_state"`
	StrategyVersion     string `json:"strategy_version"`
	LastChecked         string `json:"last_checked"`
	LastTradeSignal     string `json:"last_trade_signal"`
	LastTradeSignalTime string `json:"last_trade_signal_time"`
}

// V2
// "fetch_result_status": "SUCCESS",
// "fetch_type": "DIRECT",
// "strategy_state": "LONG",
// "strategy_version": "eth-r15",
// "last_checked": "2022-07-20 01:03:03",
// "last_trade_signal": "LONG",
// "last_trade_signal_time": "2022-07-16 16:15:04"

type SignalType string

const (
	Short SignalType = "short"
	Sell  SignalType = "sell"
	Null  SignalType = "null"
	Long  SignalType = "long"
	Buy   SignalType = "buy"
)

type SignalResponseV1 struct {
	Chain         string     `json:"chain"`
	Signal        SignalType `json:"signal"`
	LastData      time.Time  `json:"last_data"`
	CurrentTime   time.Time  `json:"current_time"`
	LastTrade     SignalType `json:"last_trade"`
	LastTradeTime time.Time  `json:"last_trade_time"`
}

type SignalResponseV2 struct {
	Chain               string     `json:"chain"`
	FetchResultStatus   string     `json:"fetch_result_status"`
	FetchType           string     `json:"fetch_type"`
	StrategyState       SignalType `json:"strategy_state"`
	StrategyVersion     string     `json:"strategy_version"`
	LastChecked         time.Time  `json:"last_checked"`
	LastTradeSignal     SignalType `json:"last_trade_signal"`
	LastTradeSignalTime time.Time  `json:"last_trade_signal_time"`
}

// Client is the instance variable of the signal api
type Client struct {
	httpClient client.Client
}

func New(httpClient client.Client) *Client {
	return &Client{
		httpClient: httpClient,
	}
}

func (c *Client) GetSignalFromIPV1(
	ctx context.Context,
	wl wlog.Logger,
	signalSource entities.SignalSource,
) (*SignalResponseV1, error) {
	url := fmt.Sprintf("http://%s/test-api", signalSource.IP)
	req, err := client.NewJSONRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error building signal request: %w", err)
	}

	var res rawSignalResponseV1
	if err := c.httpClient.DoJSON(req, &res); err != nil {
		return nil, fmt.Errorf("error making request to ip:%s signal: %w", signalSource.IP, err)
	}

	signal := &SignalResponseV1{
		Chain:         string(signalSource.Type),
		Signal:        SignalType(strings.ToLower(res.Signal)),
		LastData:      res.LastData.InUTC(),
		CurrentTime:   res.CurrentTime.InUTC(),
		LastTrade:     SignalType(strings.ToLower(res.LastTrade)),
		LastTradeTime: res.LastTradeTime.InUTC(),
	}

	return signal, nil

}

func (c *Client) GetSignalFromIPV2(
	ctx context.Context,
	wl wlog.Logger,
	signalSource entities.SignalSource,
) (*SignalResponseV2, error) {
	url := fmt.Sprintf("https://%s", signalSource.IP)
	req, err := client.NewJSONRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error building signal request: %w", err)
	}

	var res rawSignalResponseV2
	if err := c.httpClient.DoJSON(req, &res); err != nil {
		return nil, fmt.Errorf("error making request to signal: %w", err)
	}

	// parse timestamps if they look like timestamps
	lastCheckedTime, err := time.Parse(dateTimeFormat, res.LastChecked)
	if err != nil {
		wl.Debugf("invalid last_checked_time: %s", res.LastChecked)
	}
	lastTradeSignalTime, err := time.Parse(dateTimeFormat, res.LastTradeSignalTime)
	if err != nil {
		wl.Debugf("invalid last_trade_signal_time: %s", res.LastTradeSignalTime)
	}

	signal := &SignalResponseV2{
		Chain:               string(signalSource.Type),
		FetchResultStatus:   res.FetchResultStatus,
		FetchType:           res.FetchType,
		StrategyState:       SignalType(strings.ToLower(res.StrategyState)),
		StrategyVersion:     res.StrategyVersion,
		LastChecked:         lastCheckedTime,
		LastTradeSignal:     res.LastTradeSignal,
		LastTradeSignalTime: lastTradeSignalTime,
	}

	return signal, nil

}
