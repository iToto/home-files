package okxapi

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/client"
	"yield-mvp/pkg/exchangeclient"
)

type StrategyCredentials struct {
	Name       string
	APIKey     string
	SecretKey  string
	Passphrase string
}

func (s *StrategyCredentials) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("missing required Name config property")
	}
	if s.APIKey == "" {
		return fmt.Errorf("missing required APIKey config property")
	}
	if s.Passphrase == "" {
		return fmt.Errorf("missing required Passphrase config property")
	}
	if s.SecretKey == "" {
		return fmt.Errorf("missing required SecretKey config property")
	}

	return nil
}

type Config struct {
	URL       string
	Stategies []StrategyCredentials
}

func (c *Config) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("missing required URL config property")
	}
	if len(c.Stategies) == 0 {
		return fmt.Errorf("missing required Stategies config property")
	}

	return nil
}

type rawBalanceResponse struct {
	Code string `json:"code"`
	Data []struct {
		AdjEq      string `json:"adjEq"`
		BorrowFroz string `json:"borrowFroz"`
		Details    []struct {
			AvailBal      string `json:"availBal"`
			AvailEq       string `json:"availEq"`
			BorrowFroz    string `json:"borrowFroz"`
			CashBal       string `json:"cashBal"`
			Ccy           string `json:"ccy"`
			CrossLiab     string `json:"crossLiab"`
			DisEq         string `json:"disEq"`
			Eq            string `json:"eq"`
			EqUsd         string `json:"eqUsd"`
			FixedBal      string `json:"fixedBal"`
			FrozenBal     string `json:"frozenBal"`
			Imr           string `json:"imr"`
			Interest      string `json:"interest"`
			IsoEq         string `json:"isoEq"`
			IsoLiab       string `json:"isoLiab"`
			IsoUpl        string `json:"isoUpl"`
			Liab          string `json:"liab"`
			MaxLoan       string `json:"maxLoan"`
			MgnRatio      string `json:"mgnRatio"`
			Mmr           string `json:"mmr"`
			NotionalLever string `json:"notionalLever"`
			OrdFrozen     string `json:"ordFrozen"`
			RewardBal     string `json:"rewardBal"`
			SpotInUseAmt  string `json:"spotInUseAmt"`
			SpotIsoBal    string `json:"spotIsoBal"`
			StgyEq        string `json:"stgyEq"`
			Twap          string `json:"twap"`
			UTime         string `json:"uTime"`
			Upl           string `json:"upl"`
			UplLiab       string `json:"uplLiab"`
		} `json:"details"`
		Imr         string `json:"imr"`
		IsoEq       string `json:"isoEq"`
		MgnRatio    string `json:"mgnRatio"`
		Mmr         string `json:"mmr"`
		NotionalUsd string `json:"notionalUsd"`
		OrdFroz     string `json:"ordFroz"`
		TotalEq     string `json:"totalEq"`
		UTime       string `json:"uTime"`
		Upl         string `json:"upl"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type rawPositionResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Adl            string `json:"adl"`
		AvailPos       string `json:"availPos"`
		AvgPx          string `json:"avgPx"`
		CTime          string `json:"cTime"`
		Ccy            string `json:"ccy"`
		DeltaBS        string `json:"deltaBS"`
		DeltaPA        string `json:"deltaPA"`
		GammaBS        string `json:"gammaBS"`
		GammaPA        string `json:"gammaPA"`
		Imr            string `json:"imr"`
		InstID         string `json:"instId"`
		InstType       string `json:"instType"`
		Interest       string `json:"interest"`
		IdxPx          string `json:"idxPx"`
		UsdPx          string `json:"usdPx"`
		BePx           string `json:"bePx"`
		Last           string `json:"last"`
		Lever          string `json:"lever"`
		Liab           string `json:"liab"`
		LiabCcy        string `json:"liabCcy"`
		LiqPx          string `json:"liqPx"`
		MarkPx         string `json:"markPx"`
		Margin         string `json:"margin"`
		MgnMode        string `json:"mgnMode"`
		MgnRatio       string `json:"mgnRatio"`
		Mmr            string `json:"mmr"`
		NotionalUsd    string `json:"notionalUsd"`
		OptVal         string `json:"optVal"`
		PTime          string `json:"pTime"`
		Pos            string `json:"pos"`
		BaseBorrowed   string `json:"baseBorrowed"`
		BaseInterest   string `json:"baseInterest"`
		QuoteBorrowed  string `json:"quoteBorrowed"`
		QuoteInterest  string `json:"quoteInterest"`
		PosCcy         string `json:"posCcy"`
		PosID          string `json:"posId"`
		PosSide        string `json:"posSide"`
		SpotInUseAmt   string `json:"spotInUseAmt"`
		SpotInUseCcy   string `json:"spotInUseCcy"`
		BizRefID       string `json:"bizRefId"`
		BizRefType     string `json:"bizRefType"`
		ThetaBS        string `json:"thetaBS"`
		ThetaPA        string `json:"thetaPA"`
		TradeID        string `json:"tradeId"`
		UTime          string `json:"uTime"`
		Upl            string `json:"upl"`
		UplLastPx      string `json:"uplLastPx"`
		UplRatio       string `json:"uplRatio"`
		UplRatioLastPx string `json:"uplRatioLastPx"`
		VegaBS         string `json:"vegaBS"`
		VegaPA         string `json:"vegaPA"`
		RealizedPnl    string `json:"realizedPnl"`
		Pnl            string `json:"pnl"`
		Fee            string `json:"fee"`
		FundingFee     string `json:"fundingFee"`
		LiqPenalty     string `json:"liqPenalty"`
		CloseOrderAlgo []struct {
			AlgoID          string `json:"algoId"`
			SlTriggerPx     string `json:"slTriggerPx"`
			SlTriggerPxType string `json:"slTriggerPxType"`
			TpTriggerPx     string `json:"tpTriggerPx"`
			TpTriggerPxType string `json:"tpTriggerPxType"`
			CloseFraction   string `json:"closeFraction"`
		} `json:"closeOrderAlgo"`
	} `json:"data"`
}

type Client struct {
	conf       Config
	httpClient client.Client
}

func New(conf Config, httpClient client.Client) (*Client, error) {
	if conf.Validate() != nil {
		return nil, errors.New("invalid config")
	}
	for _, strategy := range conf.Stategies {
		if strategy.Validate() != nil {
			return nil, errors.New("invalid strategy")

		}
	}
	return &Client{
		conf:       conf,
		httpClient: httpClient,
	}, nil
}

// This needs to implement the exchangeclient.Client interface
func (c *Client) GetBalances(ctx context.Context, wl wlog.Logger) ([]*exchangeclient.Balance, error) {
	var balances []*exchangeclient.Balance
	endpoint := "/api/v5/account/balance"
	url := c.conf.URL + endpoint

	// get balances for each strategy we have a config fog
	for _, strategy := range c.conf.Stategies {
		// requestSignature is computed as follows:
		// Create a prehash string of timestamp + method + requestPath + body (where + represents String concatenation).
		// Prepare the SecretKey.
		// Sign the prehash string with the SecretKey using the HMAC SHA256.
		// Encode the signature in the Base64 format.

		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		requestSignature := timestamp + "GET" + endpoint
		secretKey := strategy.SecretKey
		mac := hmac.New(sha256.New, []byte(secretKey))
		mac.Write([]byte(requestSignature))
		requestSignature = base64.StdEncoding.EncodeToString(mac.Sum(nil))

		// get the balance for the strategy
		headers := http.Header{
			"OK-ACCESS-KEY":        strategy.APIKey,
			"OK-ACCESS-PASSPHRASE": strategy.Passphrase,
			"OK-ACCESS-SIGN":       requestSignature,
			"OK-ACCESS-TIMESTAMP":  timestamp,
		}

		// make the request
		req, err := client.NewJSONRequestWithHeaders(ctx, http.MethodGet, url, nil, headers)
		if err != nil {
			return nil, err
		}

		var res rawBalanceResponse
		if err := c.httpClient.DoJSON(req, &res); err != nil {
			return nil, err
		}

		// parse the response
		balance := &exchangeclient.Balance{
			EquityUSD: rawBalanceResponse.Data[0].TotalEq,
		}

		// fill out currency specific balances from rawBalanceResponse.Data[0].Details
		for _, currencyBalance := range rawBalanceResponse.Data[0].Details {
			switch currencyBalance.Ccy {
			case "BTC":
				balance.BTCAvailableBalance = currencyBalance.AvailBal
				balance.BTCInitialMarginRequirement = currencyBalance.Imr
				balance.BTCEquityOfCurrency = currencyBalance.Eq
			case "ETH":
				balance.ETHAvailableBalance = currencyBalance.AvailBal
				balance.ETHInitialMarginRequirement = currencyBalance.Imr
				balance.ETHEquityOfCurrency = currencyBalance.Eq
			case "USDT":
				balance.USDTAvailableBalance = currencyBalance.AvailBal
				balance.USDTInitialMarginRequirement = currencyBalance.Imr
				balance.USDTEquityOfCurrency = currencyBalance.Eq
			case "USDC":
				balance.USDCAvailableBalance = currencyBalance.AvailBal
				balance.USDCInitialMarginRequirement = currencyBalance.Imr
				balance.USDCEquityOfCurrency = currencyBalance.Eq
			}
		}

		balances = append(balances, balance)
	}

	return balances, nil
}

func (c *Client) GetPositions(ctx context.Context, wl wlog.Logger) ([]*exchangeclient.Position, error) {
	// TODO: add HTTP call here

	positionJSON := c.getTestPosition()

	var rawPositionResponse rawPositionResponse
	err := json.Unmarshal([]byte(positionJSON), &rawPositionResponse)
	if err != nil {
		return nil, err
	}

	position := &exchangeclient.Position{
		Position:     rawPositionResponse.Data[0].Pos,
		NotionalUSD:  rawPositionResponse.Data[0].NotionalUsd,
		UnrealizedPL: rawPositionResponse.Data[0].Upl,
		MarketPrice:  rawPositionResponse.Data[0].MarkPx,
	}

	return position, nil
}

func (c *Client) getTestPosition() string {
	return `{
  "code": "0",
  "msg": "",
  "data": [
    {
      "adl":"1",
      "availPos":"1",
      "avgPx":"2566.31",
      "cTime":"1619507758793",
      "ccy":"ETH",
      "deltaBS":"",
      "deltaPA":"",
      "gammaBS":"",
      "gammaPA":"",
      "imr":"",
      "instId":"ETH-USD-210430",
      "instType":"FUTURES",
      "interest":"0",
      "idxPx":"2566.13",
      "usdPx":"",
      "bePx":"2353.949",
      "last":"2566.22",
      "lever":"10",
      "liab":"",
      "liabCcy":"",
      "liqPx":"2352.8496681818233",
      "markPx":"2353.849",
      "margin":"0.0003896645377994",
      "mgnMode":"isolated",
      "mgnRatio":"11.731726509588816",
      "mmr":"0.0000311811092368",
      "notionalUsd":"2276.2546609009605",
      "optVal":"",
      "pTime":"1619507761462",
      "pos":"1",
      "baseBorrowed": "",
      "baseInterest": "",
      "quoteBorrowed": "",
      "quoteInterest": "",
      "posCcy":"",
      "posId":"307173036051017730",
      "posSide":"long",
      "spotInUseAmt": "",
      "spotInUseCcy": "",
      "bizRefId": "",
      "bizRefType": "",
      "thetaBS":"",
      "thetaPA":"",
      "tradeId":"109844",
      "uTime":"1619507761462",
      "upl":"-0.0000009932766034",
      "uplLastPx":"-0.0000009932766034",
      "uplRatio":"-0.0025490556801078",
      "uplRatioLastPx":"-0.0025490556801078",
      "vegaBS":"",
      "vegaPA":"",
      "realizedPnl":"0.001",
      "pnl":"0.0011",
      "fee":"-0.0001",
      "fundingFee":"0",
      "liqPenalty":"0",
      "closeOrderAlgo":[
          {
              "algoId":"123",
              "slTriggerPx":"123",
              "slTriggerPxType":"mark",
              "tpTriggerPx":"123",
              "tpTriggerPxType":"mark",
              "closeFraction":"0.6"
          },
          {
              "algoId":"123",
              "slTriggerPx":"123",
              "slTriggerPxType":"mark",
              "tpTriggerPx":"123",
              "tpTriggerPxType":"mark",
              "closeFraction":"0.4"
          }
      ]
    }
  ]
}
`
}

func (c *Client) getTestBalance() string {
	return `{
    "code": "0",
    "data": [
        {
            "adjEq": "55415.624719833286",
            "borrowFroz": "0",
            "details": [
                {
                    "availBal": "4834.317093622894",
                    "availEq": "4834.3170936228935",
                    "borrowFroz": "0",
                    "cashBal": "4850.435693622894",
                    "ccy": "USDT",
                    "crossLiab": "0",
                    "disEq": "4991.542013297616",
                    "eq": "4992.890093622894",
                    "eqUsd": "4991.542013297616",
                    "fixedBal": "0",
                    "frozenBal": "158.573",
                    "imr": "",
                    "interest": "0",
                    "isoEq": "0",
                    "isoLiab": "0",
                    "isoUpl": "0",
                    "liab": "0",
                    "maxLoan": "0",
                    "mgnRatio": "",
                    "mmr": "",
                    "notionalLever": "",
                    "ordFrozen": "0",
                    "rewardBal": "0",
                    "spotInUseAmt": "",
                    "spotIsoBal": "0",
                    "stgyEq": "150",
                    "twap": "0",
                    "uTime": "1705449605015",
                    "upl": "-7.545600000000006",
                    "uplLiab": "0"
                }
            ],
            "imr": "8.57068529",
            "isoEq": "0",
            "mgnRatio": "143682.59776662575",
            "mmr": "0.3428274116",
            "notionalUsd": "85.7068529",
            "ordFroz": "0",
            "totalEq": "55837.43556134779",
            "uTime": "1705474164160",
            "upl": "-7.543562688000006"
        }
    ],
    "msg": ""
}
`
}
