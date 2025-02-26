package main

import (
	"context"
	"flag"
	"strconv"
	"yield-mvp/internal/balancelogger"
	"yield-mvp/internal/emailer"
	"yield-mvp/internal/entities"
	"yield-mvp/internal/exchangeDAL"
	"yield-mvp/internal/handler"
	"yield-mvp/internal/orderDAL"
	"yield-mvp/internal/orderlogger"
	"yield-mvp/internal/service/exchangesvc"
	"yield-mvp/internal/service/reportsvc"
	"yield-mvp/internal/service/signalsvc"
	"yield-mvp/internal/signalSourceDAL"
	"yield-mvp/internal/signallogger"
	"yield-mvp/internal/strategyDAL"
	"yield-mvp/internal/tickr"
	"yield-mvp/internal/tradelogger"
	"yield-mvp/internal/userDAL"
	"yield-mvp/internal/wlog"
	"yield-mvp/pkg/client"
	"yield-mvp/pkg/coinroutesapi"
	"yield-mvp/pkg/coinroutespriceconsumer"
	"yield-mvp/pkg/exchangeclient/okxapi"
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

	// get user_id
	userID := os.Getenv("USER_ID")
	if userID == "" {
		log.Fatal("$USER_ID environment variable must be set")
	}

	userDAL, err := userDAL.New(db)
	if err != nil {
		log.Fatal("could not create user DAL")
	}

	signalDAL, err := signalSourceDAL.New(db)
	if err != nil {
		log.Fatal("could not create signal source DAL")
	}

	strategyDAL, err := strategyDAL.New(db)
	if err != nil {
		log.Fatal("could not create strategy DAL")
	}

	orderDAL, err := orderDAL.New(db)
	if err != nil {
		log.Fatal("could not create order DAL", err)
	}

	stubSignal := false
	stubSignalEnv := os.Getenv("STUB_SIGNAL")
	if stubSignalEnv == "true" {
		wl.Info("stubbing out signal API")
		stubSignal = true
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

	// email client
	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		log.Fatal("$SMTP_PORT environment variable must be set")
	}
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		log.Fatal("$SMTP_HOST environment variable must be set")
	}
	smtpUser := os.Getenv("SMTP_USER")
	if smtpUser == "" {
		log.Fatal("$SMTP_USER environment variable must be set")
	}
	smtpPass := os.Getenv("SMTP_PASS")
	if smtpPass == "" {
		log.Fatal("$SMTP_PASS environment variable must be set")
	}
	csvPath := os.Getenv("PATH_TO_CSV")
	if csvPath == "" {
		log.Fatal("$PATH_TO_CSV environment variable must be set")
	}

	p, _ := strconv.ParseInt(smtpPort, 10, 64)

	wl.Debugf("%s, %s, %s, %s, %d", smtpHost, smtpUser, smtpPass, smtpPort, p)

	ec, err := emailer.New(smtpUser, smtpPass, smtpHost, int(p))
	if err != nil {
		log.Fatal(err)
	}

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
		signalDAL,
		strategyDAL,
		userDAL,
		ethPriceConsumer,
		btcPriceConsumer,
	)
	if err != nil {
		log.Fatal("unable to init signal service: %w", err)
	}

	reportService, err := reportsvc.New(
		db,
		orderDAL,
		strategyDAL,
		ethPriceConsumer,
		btcPriceConsumer,
		ec,
		csvPath,
	)
	if err != nil {
		log.Fatal("unable to init report service: %w", err)
	}

	okxConf := okxapi.Config{
		URL:        "https://www.okx.com/api/v5",
		APIKey:     "foo",
		Passphrase: "bar",
	}
	okxExchangeClient := okxapi.New(okxConf, client.NewHTTPClient())
	exchangesvc, err := exchangesvc.New(okxExchangeClient)
	if err != nil {
		log.Fatal("unable to init exchange service: %w", err)
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

	ticker, err := tickr.New(
		signalResolutionSeconds,
		entities.BTC,
		userID,
		signalService)
	if err != nil {
		log.Fatal("unable to setup btc ticker")
	}

	ticker.Run(context.Background(), wl)

	// setup router and handlers
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc(
		"/proc-user",
		handler.ProcessSignalsForUser(
			wl,
			signalService,
			userID)).Methods(http.MethodGet, http.MethodOptions)
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
			strategyDAL,
			bl,
			exchangeDAL),
	).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc(
		"/update_order",
		handler.UpdateOpenOrders(
			wl,
			coinRoutesClient,
			db,
			ol,
		)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc(
		"/disable_strategy/{name}",
		handler.DisableStrategy(
			wl,
			signalService,
		)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc(
		"/enable_strategy/{name}",
		handler.EnableStrategy(
			wl,
			signalService,
		)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc(
		"/strategy",
		handler.GetStrategies(
			wl,
			signalService,
		)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc(
		"/strategy",
		handler.CreateStrategy(
			wl,
			signalService,
		)).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc(
		"/strategy",
		handler.UpdateStrategy(
			wl,
			signalService,
		)).Methods(http.MethodPatch, http.MethodOptions)
	router.HandleFunc(
		"/signal",
		handler.CreateSignal(
			wl,
			signalService,
		)).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc(
		"/signal",
		handler.GetSignals(
			wl,
			signalService,
		)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc(
		"/report",
		handler.GenerateReport(
			wl,
			reportService,
		)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc(
		"/okx_report",
		handler.GenerateExchangeReport(
			wl,
			exchangesvc,
		)).Methods(http.MethodGet, http.MethodOptions)
	wl.Debugf("running on port: %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
