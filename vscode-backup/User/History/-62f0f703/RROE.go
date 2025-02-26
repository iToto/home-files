package main

import (
	"context"
	"flag"
	"strconv"
	"yield-mvp/internal/balancelogger"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/exchangeDAL"
	"yield-mvp/internal/handler"
	"yield-mvp/internal/orderlogger"
	"yield-mvp/internal/service/signalsvc"
	"yield-mvp/internal/signallogger"
	"yield-mvp/internal/tickr"
	"yield-mvp/internal/tradelogger"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/client"
	"yield-mvp/pkg/coinroutesapi"
	"yield-mvp/pkg/coinroutespriceconsumer"
	"yield-mvp/pkg/signalapi"
	"yield-mvp/pkg/utils"

	// "encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	pathRealPrice    = "/api/streaming/real_price/"
	pathCbbo         = "/api/streaming/cbbo/"
	defaultSignalRes = int64(30)
)

func DBConnection() (*sqlx.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPwd := os.Getenv("DB_PASS")
	instanceConnectionName := os.Getenv("INSTANCE_CONNECTION_NAME")
	dbName := os.Getenv("DB_NAME")
	dbURL := os.Getenv("DATABASE_URL")

	// used for local
	if dbURL != "" {
		db, err := sqlx.Open("postgres", dbURL)
		if err != nil {
			return nil, err
		}

		return db, nil
	}

	if dbUser == "" {
		return nil, fmt.Errorf("missing required env var DB_USER")
	}
	if dbPwd == "" {
		return nil, fmt.Errorf("missing required env var DB_PASS")
	}
	if instanceConnectionName == "" {
		return nil, fmt.Errorf("missing required env var INSTANCE_CONNECTION_NAME")
	}
	if dbName == "" {
		return nil, fmt.Errorf("missing required env var DB_NAME")
	}

	socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
	if !isSet {
		socketDir = "/cloudsql"
	}

	dbURI := fmt.Sprintf("user=%s password=%s database=%s host=%s/%s",
		dbUser,
		dbPwd,
		dbName,
		socketDir,
		instanceConnectionName,
	)

	// dbPool is the pool of database connections.
	db, err := sqlx.Open("postgres", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	return db, nil
}

func WSConnection(path string, wl wlog.Logger) (*websocket.Conn, error) {
	baseURL := os.Getenv("CR_WS_BASE_URL")
	token := os.Getenv("CR_WS_TOKEN")

	if baseURL == "" {
		return nil, fmt.Errorf("missing required env var CR_WS_BASE_URL")
	}
	if token == "" {
		return nil, fmt.Errorf("missing required env var CR_WS_TOKEN")
	}
	if path == "" {
		return nil, fmt.Errorf("invalid path for websocket connection: %s", path)
	}

	u := url.URL{Scheme: "wss", Host: baseURL, Path: path}
	tokenHeader := fmt.Sprintf("Token %s", token)
	header := http.Header{}
	header.Set("Authorization", tokenHeader)
	wl.Debugf("connecting to websocket: %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Fatal("dial:", err)
	}

	return c, nil
}

func main() {
	var env string
	var local bool
	flag.StringVar(&env, "env", "", "path to env file")
	flag.StringVar(&env, "e", "", "shorthand for env")
	flag.BoolVar(&local, "local", false, "local mode")

	flag.Parse()

	if env != "" {
		if err := utils.PrimeEnv(env, local); err != nil {
			log.Fatalf("error priming eng: %s", err)
		}
	}

	wl, err := wlog.NewBasicLogger()
	if err != nil {
		log.Fatal("error configuring logger")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT environment variable must be set")
	}

	// setup db
	db, err := DBConnection()
	if err != nil {
		log.Fatal("unable to setup db: %w", err)
	}
	defer db.Close()

	// setup ws for Real Price for btc and eth
	btcWS, err := WSConnection(pathRealPrice, wl)
	if err != nil {
		log.Fatal("unable to setup btc ws: %w", err)
	}
	defer btcWS.Close()

	ethWS, err := WSConnection(pathRealPrice, wl)
	if err != nil {
		log.Fatal("unable to setup eth ws: %w", err)
	}
	defer ethWS.Close()

	// setup websocket for real price consuming from coinroutes
	btcPriceConsumerConfig := coinroutespriceconsumer.Config{
		ResourcePath: pathRealPrice,
		Payload: coinroutespriceconsumer.RealPriceRequest{
			CurrencyPair: "BTC-USD",
			Quantity:     "1.00",
		},
	}
	ethPriceConsumerConfig := coinroutespriceconsumer.Config{
		ResourcePath: pathRealPrice,
		Payload: coinroutespriceconsumer.RealPriceRequest{
			CurrencyPair: "ETH-USD",
			Quantity:     "1.00",
		},
	}

	btcPriceConsumer, err := coinroutespriceconsumer.NewConsumer(btcPriceConsumerConfig)
	if err != nil {
		log.Fatal("could not build btc consumer with error: ", err)
	}

	ethPriceConsumer, err := coinroutespriceconsumer.NewConsumer(ethPriceConsumerConfig)
	if err != nil {
		log.Fatal("could not build eth consumer with error: ", err)
	}

	btcWL := wlog.WithChain(wl, string(entities.BTC))
	ethWL := wlog.WithChain(wl, string(entities.ETH))

	btcPriceConsumer.Start(btcWS, btcWL)
	ethPriceConsumer.Start(ethWS, ethWL)

	// setup signal config
	btcURL := os.Getenv("BTC_URL")
	if btcURL == "" {
		log.Fatal("$BTC_UTL environment variable must be set")
	}
	ethURL := os.Getenv("ETH_URL")
	if ethURL == "" {
		log.Fatal("$ETH_UTL environment variable must be set")
	}
	stubSignal := false
	stubSignalEnv := os.Getenv("STUB_SIGNAL")
	if stubSignalEnv == "true" {
		wl.Info("stubbing out signal API")
		stubSignal = true
	}

	// setup signal sources
	signalSources := map[entities.ChainType]entities.SignalSource{
		entities.BTC: {
			Type:          entities.BTC,
			IP:            btcURL,
			Enabled:       true,
			SignalVersion: entities.V2,
		},
		entities.ETH: {
			Type:          entities.ETH,
			IP:            ethURL,
			Enabled:       true,
			SignalVersion: entities.V2,
		},
	}

	simulateTrade := false
	simEnv := os.Getenv("SIMULATE_TRADE")
	if simEnv == "true" {
		wl.Info("simulating trades")
		simulateTrade = true
	}

	// setup coinroutes config
	crURL := os.Getenv("CR_URL")
	if crURL == "" {
		log.Fatal("$CR_URL environment variable must be set")
	}
	crToken := os.Getenv("CR_TOKEN")
	if crToken == "" {
		log.Fatal("$CR_TOKEN environment variable must be set")
	}

	ftxUsdmBtcLs := os.Getenv("FTX_USDM_BTC_LS")
	if ftxUsdmBtcLs == "" {
		log.Fatal("$FTX_USDM_BTC_LS environment variable must be set")
	}
	ftxUsdmEthLs := os.Getenv("FTX_USDM_ETH_LS")
	if ftxUsdmEthLs == "" {
		log.Fatal("$FTX_USDM_ETH_LS environment variable must be set")
	}
	ftxUsdmBtcComp := os.Getenv("FTX_USDM_BTC_COMP")
	if ftxUsdmBtcComp == "" {
		log.Fatal("$FTX_USDM_BTC_COMP environment variable must be set")
	}
	ftxUsdmEthComp := os.Getenv("FTX_USDM_ETH_COMP")
	if ftxUsdmEthComp == "" {
		log.Fatal("$FTX_USDM_ETH_COMP environment variable must be set")
	}

	binanceUsdmBtcLs := os.Getenv("BINANCE_USDM_BTC_LS")
	if binanceUsdmBtcLs == "" {
		log.Fatal("$BINANCE_USDM_BTC_LS environment variable must be set")
	}
	binanceUsdmEthLs := os.Getenv("BINANCE_USDM_ETH_LS")
	if binanceUsdmEthLs == "" {
		log.Fatal("$BINANCE_USDM_ETH_LS environment variable must be set")
	}
	binanceUsdmBtcLs2X := os.Getenv("BINANCE_USDM_BTC_LS2X")
	if binanceUsdmBtcLs2X == "" {
		log.Fatal("$BINANCE_USDM_BTC_LS2X environment variable must be set")
	}
	binanceUsdmEthLs2X := os.Getenv("BINANCE_USDM_ETH_LS2X")
	if binanceUsdmEthLs2X == "" {
		log.Fatal("$BINANCE_USDM_ETH_LS2X environment variable must be set")
	}

	// build Chain/Strategy struct

	ftxUsdmBtcLsStrategy := entities.Strategy{
		Type:             entities.LongNeutral,
		Name:             ftxUsdmBtcLs,
		Exchange:         entities.FTX,
		Margin:           entities.USDTM,
		Leverage:         entities.OneX,
		TradeStrategy:    entities.Fixed,
		FixedTradeAmount: 500.00,
	}
	ftxUsdmEthLsStrategy := entities.Strategy{
		Type:             entities.LongNeutral,
		Name:             ftxUsdmEthLs,
		Exchange:         entities.FTX,
		Margin:           entities.USDTM,
		Leverage:         entities.OneX,
		TradeStrategy:    entities.Fixed,
		FixedTradeAmount: 500.00,
	}
	ftxUsdmBtcCompStrategy := entities.Strategy{
		Type:          entities.LongNeutral,
		Name:          ftxUsdmBtcComp,
		Exchange:      entities.FTX,
		Margin:        entities.USDTM,
		Leverage:      entities.OneX,
		TradeStrategy: entities.Compound,
	}
	ftxUsdmEthCompStrategy := entities.Strategy{
		Type:          entities.LongNeutral,
		Name:          ftxUsdmEthComp,
		Exchange:      entities.FTX,
		Margin:        entities.USDTM,
		Leverage:      entities.OneX,
		TradeStrategy: entities.Compound,
	}

	binanceUsdmBtcLsStrategy := entities.Strategy{
		Type:          entities.LongNeutral,
		Name:          binanceUsdmBtcLs,
		Exchange:      entities.Binance,
		Margin:        entities.USDM,
		Leverage:      entities.OneX,
		TradeStrategy: entities.Compound,
	}
	binanceUsdmEthLsStrategy := entities.Strategy{
		Type:          entities.LongNeutral,
		Name:          binanceUsdmEthLs,
		Exchange:      entities.Binance,
		Margin:        entities.USDM,
		Leverage:      entities.OneX,
		TradeStrategy: entities.Compound,
	}
	binanceUsdmBtcLs2XStrategy := entities.Strategy{
		Type:             entities.LongNeutral,
		Name:             binanceUsdmBtcLs2X,
		Exchange:         entities.Binance,
		Margin:           entities.USDM,
		Leverage:         entities.OneX,
		TradeStrategy:    entities.Fixed,
		FixedTradeAmount: 500.00,
	}
	binanceUsdmEthLs2XStrategy := entities.Strategy{
		Type:             entities.LongNeutral,
		Name:             binanceUsdmEthLs2X,
		Exchange:         entities.Binance,
		Margin:           entities.USDM,
		Leverage:         entities.OneX,
		TradeStrategy:    entities.Fixed,
		FixedTradeAmount: 500.00,
	}

	// ETH
	ethLN := entities.Strategy{
		Type:     entities.LongNeutral,
		Name:     ethLongNeutralStrategy,
		Exchange: entities.Binance,
		Margin:   entities.CoinM,
		Leverage: entities.OneX,
	}
	ethUSDT := entities.Strategy{
		Type:             entities.USDT,
		Name:             ethUSDTStrategy,
		Exchange:         entities.Binance,
		Margin:           entities.USDM,
		Leverage:         entities.OneX,
		TradeStrategy:    entities.Fixed,
		FixedTradeAmount: 500.00,
	}
	binUsdmEthLs := entities.Strategy{
		Type:             entities.LongNeutral,
		Name:             binanceUsdmEthLs,
		Exchange:         entities.Binance,
		Margin:           entities.USDM,
		Leverage:         entities.OneX,
		TradeStrategy:    entities.Fixed,
		FixedTradeAmount: 19000.00,
	}
	binUsdmEthLs2x := entities.Strategy{
		Type:             entities.LongNeutral,
		Name:             binanceUsdmEthLs2X,
		Exchange:         entities.Binance,
		Margin:           entities.USDM,
		Leverage:         entities.TwoX,
		TradeStrategy:    entities.Fixed,
		FixedTradeAmount: 10000.00,
	}
	ftxEthComp := entities.Strategy{
		Type:          entities.USD,
		Name:          ftxEthCompound,
		Exchange:      entities.FTX,
		Margin:        entities.USDTM,
		Leverage:      entities.OneX,
		TradeStrategy: entities.Compound,
	}
	ethChain := entities.Chain{
		Type:       entities.ETH,
		Strategies: []entities.Strategy{ethLN, ethUSDT, binUsdmEthLs, binUsdmEthLs2x, ftxEthComp},
	}
	// BTC
	btcLN := entities.Strategy{
		Type:     entities.LongNeutral,
		Name:     btcLongNeutralStrategy,
		Exchange: entities.Binance,
		Margin:   entities.CoinM,
		Leverage: entities.OneX,
	}
	btcLS := entities.Strategy{
		Type:     entities.LongShort,
		Name:     btcLongShortStrategy,
		Exchange: entities.Binance,
		Margin:   entities.CoinM,
		Leverage: entities.OneX,
	}
	btcFTXLS := entities.Strategy{
		Type:             entities.USD,
		Name:             btcFTXLongShortStrategy,
		Exchange:         entities.FTX,
		Margin:           entities.USDTM,
		Leverage:         entities.OneX,
		TradeStrategy:    entities.Fixed,
		FixedTradeAmount: 10000.00,
	}
	binUsdmBtcLs := entities.Strategy{
		Type:          entities.LongNeutral,
		Name:          binanceUsdmBtcLs,
		Exchange:      entities.Binance,
		Margin:        entities.USDM,
		Leverage:      entities.OneX,
		TradeStrategy: entities.Compound,
		// TradeStrategy:    entities.Fixed,
		// FixedTradeAmount: 19000.00,
	}
	binUsdmBtcLs2x := entities.Strategy{
		Type:             entities.LongNeutral,
		Name:             binanceUsdmBtcLs2X,
		Exchange:         entities.Binance,
		Margin:           entities.USDM,
		Leverage:         entities.TwoX,
		TradeStrategy:    entities.Fixed,
		FixedTradeAmount: 10000.00,
	}
	btcChain := entities.Chain{
		Type:       entities.BTC,
		Strategies: []entities.Strategy{btcLN, btcLS, btcFTXLS, binUsdmBtcLs, binUsdmBtcLs2x},
	}

	// build map of chains and their associated strategies
	chains := map[entities.ChainType]entities.Chain{
		entities.BTC: btcChain,
		entities.ETH: ethChain,
	}

	wl.Info("setting up router")

	// load trade logger for BQ
	dl, err := tradelogger.NewDataLogger(wl)
	if err != nil {
		log.Fatal("unable to load trade data logger: %w", err)
	}

	// load signal logger for BQ
	sl, err := signallogger.NewDataLogger(wl)
	if err != nil {
		log.Fatal("unable to load signal data logger: %w", err)
	}

	// load balance logger for BQ
	bl, err := balancelogger.NewDataLogger(wl)
	if err != nil {
		log.Fatal("unable to load balance data logger: %w", err)
	}

	// load order logger for BQ
	ol, err := orderlogger.NewDataLogger(wl)
	if err != nil {
		log.Fatal("unable to load balance data logger: %w", err)
	}

	// load HTTP clients
	// signal client
	signalConf := signalapi.Config{
		STUBOUT: stubSignal,
	}
	signalClient := signalapi.New(signalConf, client.NewHTTPClient())

	// coinroutes client
	crConf := coinroutesapi.Config{
		URL:      crURL,
		Token:    crToken,
		Simulate: simulateTrade,
	}
	coinRoutesClient := coinroutesapi.New(crConf, client.NewHTTPClient())

	// setup DALs
	exchangeDAL, err := exchangeDAL.New(coinRoutesClient)
	if err != nil {
		log.Fatal("unable to init exchange DAL %w", err)
	}

	// setup services
	signalService, err := signalsvc.New(
		signalClient,
		coinRoutesClient,
		db,
		dl,
		sl,
		exchangeDAL,
		ethPriceConsumer,
		btcPriceConsumer,
	)
	if err != nil {
		log.Fatal("unable to init signal service: %w", err)
	}

	// setup signal listener to ping every Xs
	signalResolutionSeconds := defaultSignalRes
	signalResolution := os.Getenv("SIGNAL_RESOLUTION")

	if signalResolution == "" {
		wl.Infof("no signal resolution found, defaulting to: %d", signalResolutionSeconds)
	} else {
		signalResolutionSeconds, err = strconv.ParseInt(signalResolution, 10, 64)
		if err != nil {
			log.Fatal("unable to parse env var SIGNAL_RESOLUTION %w", err)
		}
		wl.Infof("signal resolution: %d", signalResolutionSeconds)
	}

	btcTicker, err := tickr.New(
		signalResolutionSeconds,
		entities.BTC,
		chains[entities.BTC].Strategies,
		signalService)
	if err != nil {
		log.Fatal("unable to setup btc ticker")
	}

	ethTicker, err := tickr.New(
		signalResolutionSeconds,
		entities.ETH,
		chains[entities.ETH].Strategies,
		signalService)
	if err != nil {
		log.Fatal("unable to setup eth ticker")
	}

	btcTicker.Run(context.Background(), wl, signalSources)
	ethTicker.Run(context.Background(), wl, signalSources)

	// setup router and handlers
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc(
		"/a",
		handler.GetBTCSignal(
			wl,
			chains[entities.BTC].Strategies,
			signalService,
			signalSources)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc(
		"/b",
		handler.GetETHSignal(
			wl,
			chains[entities.ETH].Strategies,
			signalService,
			signalSources)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc(
		"/accounts",
		handler.GetCoinRoutesExchangeAccounts(
			wl,
			coinRoutesClient)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc(
		"/update_balance",
		handler.GetBalance(
			wl,
			coinRoutesClient,
			btcPriceConsumer,
			ethPriceConsumer,
			chains,
			bl, exchangeDAL),
	).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc(
		"/update_order",
		handler.UpdateOpenOrders(
			wl,
			coinRoutesClient,
			db,
			ol,
		)).Methods(http.MethodGet, http.MethodOptions)
	wl.Debugf("running on port: %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
