package exchangesvc

import (
	"context"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/exchangeclient"
)

// ReportData with json tag
type ReportData struct {
	Position                     string `json:"position,omitempty"`                        // pos
	NotionalUSD                  string `json:"notional_usd,omitempty"`                    // notionalUsd
	UnrealizedPL                 string `json:"unrealized_pl,omitempty"`                   // upl
	MarketPrice                  string `json:"market_price,omitempty"`                    // markPx
	EquityUSD                    string `json:"equity_usd,omitempty"`                      // eqUsd
	BTCAvailableBalance          string `json:"btc_available_balance,omitempty"`           // availbal (btc ccy)
	BTCInitialMarginRequirement  string `json:"btc_initial_margin_requirement,omitempty"`  // imr (btc ccy)
	BTCEquityOfCurrency          string `json:"btc_equity_of_currency,omitempty"`          // eq (btc ccy)
	ETHAvailableBalance          string `json:"eth_available_balance,omitempty"`           // availbal (eth ccy)
	ETHInitialMarginRequirement  string `json:"eth_initial_margin_requirement,omitempty"`  // imr (eth ccy)
	ETHEquityOfCurrency          string `json:"eth_equity_of_currency,omitempty"`          // eqA (eth ccy)
	USDTAvailableBalance         string `json:"usdt_available_balance,omitempty"`          // availbal (USDT ccy)
	USDTInitialMarginRequirement string `json:"usdt_initial_margin_requirement,omitempty"` // imr (USDT ccy)
	USDTEquityOfCurrency         string `json:"usdt_equity_of_currency,omitempty"`         // eq (USDT ccy)
	USDCAvailableBalance         string `json:"usdc_available_balance,omitempty"`          // availbal (USDC ccy)
	USDCInitialMarginRequirement string `json:"usdc_initial_margin_requirement,omitempty"` // imr (USDC ccy)
	USDCEquityOfCurrency         string `json:"usdc_equity_of_currency,omitempty"`         // eq (USDC ccy)
}

type SVC interface {
	GenereateReport(
		ctx context.Context,
		wl wlog.Logger,
	) (*ReportData, error)
}

type exchangeService struct {
	exchangeClient exchangeclient.Client
}

func New(
	exchangeClient exchangeclient.Client,
) (SVC, error) {
	return &exchangeService{
		exchangeClient: exchangeClient,
	}, nil
}

func (s *exchangeService) GenereateReport(
	ctx context.Context,
	wl wlog.Logger,
) (*ReportData, error) {
	b, err := s.exchangeClient.GetStrategyBalances(ctx, wl)
	if err != nil {
		return nil, err
	}

	p, err := s.exchangeClient.GetStrategyPositions(ctx, wl)
	if err != nil {
		return nil, err
	}

	wl.Debugf("balances: %v", b)
	wl.Debugf("positions: %v", p)

	var report []ReportData

	for strategy, record := range b {
		ReportData := ReportData{
			Position:                     p[strategy].Position,
			NotionalUSD:                  p[strategy].NotionalUSD,
			UnrealizedPL:                 p[strategy].UnrealizedPL,
			MarketPrice:                  p[strategy].MarketPrice,
			EquityUSD:                    b[strategy].EquityUSD,
			BTCAvailableBalance:          b[strategy].BTCAvailableBalance,
			BTCInitialMarginRequirement:  b[strategy].BTCInitialMarginRequirement,
			BTCEquityOfCurrency:          b[strategy].BTCEquityOfCurrency,
			ETHAvailableBalance:          b[strategy].ETHAvailableBalance,
			ETHInitialMarginRequirement:  b[strategy].ETHInitialMarginRequirement,
			ETHEquityOfCurrency:          b[strategy].ETHEquityOfCurrency,
			USDTAvailableBalance:         b[strategy].USDTAvailableBalance,
			USDTInitialMarginRequirement: b[strategy].USDTInitialMarginRequirement,
			USDTEquityOfCurrency:         b[strategy].USDTEquityOfCurrency,
			USDCAvailableBalance:         b[strategy].USDCAvailableBalance,
			USDCInitialMarginRequirement: b[strategy].USDCInitialMarginRequirement,
			USDCEquityOfCurrency:         b[strategy].USDCEquityOfCurrency,
		}

		report = append(report, ReportData)

	}

	return &report, nil
}
