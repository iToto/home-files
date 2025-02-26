package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"yield/signal-logger/internal/entities"
	"yield/signal-logger/internal/service/signalsvc"
	"yield/signal-logger/internal/signallogger"
	"yield/signal-logger/internal/signalloggerv2"
	"yield/signal-logger/internal/tickr"
	"yield/signal-logger/internal/wlog"
	"yield/signal-logger/pkg/client"
	"yield/signal-logger/pkg/signalapi"
	"yield/signal-logger/pkg/utils"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	defaultSignalRes = int64(5)
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

	// setup signal config
	// btcURL := os.Getenv("BTC_URL")
	// if btcURL == "" {
	// 	log.Fatal("$BTC_UTL environment variable must be set")
	// }
	// ethURL := os.Getenv("ETH_URL")
	// if ethURL == "" {
	// 	log.Fatal("$ETH_UTL environment variable must be set")
	// }

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

	// load signal logger for BQ
	sl, err := signallogger.NewDataLogger(wl)
	if err != nil {
		log.Fatal("unable to load signal data logger: %w", err)
	}
	slv2, err := signalloggerv2.NewDataLogger(wl)
	if err != nil {
		log.Fatal("unable to load signal v2 data logger: %w", err)
	}

	// signal client
	signalClient := signalapi.New(client.NewHTTPClient())

	// setup services
	signalService, err := signalsvc.New(
		db,
		signalClient,
		sl,
		slv2,
	)
	if err != nil {
		log.Fatal("unable to init signal service: %w", err)
	}

	signalSources := []entities.SignalSource{
		// {
		// 	Type:          entities.BTC,
		// 	IP:            "35.204.101.142",
		// 	Enabled:       true,
		// 	SignalVersion: entities.V1,
		// },
		// {
		// 	Type:          entities.BTC,
		// 	IP:            "35.204.251.164",
		// 	Enabled:       true,
		// 	SignalVersion: entities.V1,
		// },
		// {
		// 	Type:          entities.ETH,
		// 	IP:            "34.141.231.198",
		// 	Enabled:       true,
		// 	SignalVersion: entities.V1,
		// },
		// {
		// 	Type:          entities.ETH,
		// 	IP:            "34.90.130.174",
		// 	Enabled:       true,
		// 	SignalVersion: entities.V1,
		// },
		// {
		// 	Type:          entities.ETH,
		// 	IP:            "34.91.74.58",
		// 	Enabled:       true,
		// 	SignalVersion: entities.V1,
		// },
		// {
		// 	Type:          entities.BTC,
		// 	IP:            "api.beta.btc.xeohive.com/test-api/r17-old",
		// 	Enabled:       true,
		// 	SignalVersion: entities.V2,
		// },
		// {
		// 	Type:          entities.BTC,
		// 	IP:            "api.beta.btc.xeohive.com/test-api/r17-new",
		// 	Enabled:       true,
		// 	SignalVersion: entities.V2,
		// },
		// {
		// 	Type:          entities.BTC,
		// 	IP:            "api.beta.btc.xeohive.com/test-api/r17-neutral",
		// 	Enabled:       true,
		// 	SignalVersion: entities.V2,
		// },
		// {
		// 	Type:          entities.ETH,
		// 	IP:            "api.beta.eth.xeohive.com/test-api/r15-old",
		// 	Enabled:       true,
		// 	SignalVersion: entities.V2,
		// },
		// {
		// 	Type:          entities.ETH,
		// 	IP:            "api.beta.eth.xeohive.com/test-api/r15-new",
		// 	Enabled:       true,
		// 	SignalVersion: entities.V2,
		// },
		// {
		// 	Type:          entities.ETH,
		// 	IP:            "api.beta.eth.xeohive.com/test-api/r15-neutral",
		// 	Enabled:       true,
		// 	SignalVersion: entities.V2,
		// },
		// {
		// 	Type:          entities.BTC,
		// 	IP:            "api.beta.btc.xeohive.com/test-api/v3-beta",
		// 	Enabled:       true,
		// 	SignalVersion: entities.V2,
		// },
		// {
		// 	Type:          entities.ETH,
		// 	IP:            "api.beta.eth.xeohive.com/test-api/r15v100-beta",
		// 	Enabled:       true,
		// 	SignalVersion: entities.V2,
		// },
		// {
		// 	Type:          entities.BTC,
		// 	IP:            "api.beta.btc.xeohive.com/test-api/mix-v3-r17",
		// 	Enabled:       true,
		// 	SignalVersion: entities.V2,
		// },
		// {
		// 	Type:          entities.ETH,
		// 	IP:            "api.beta.eth.xeohive.com/test-api/mix-v2-r17",
		// 	Enabled:       true,
		// 	SignalVersion: entities.V2,
		// },
		// {
		// 	Type:          entities.ETH,
		// 	IP:            "api.beta.eth.xeohive.com/test-api/r18-beta",
		// 	Enabled:       true,
		// 	SignalVersion: entities.V2,
		// },
		// {
		// 	Type:          entities.ETH,
		// 	IP:            "api.beta.eth.xeohive.com/test-api/mix-v2-r18",
		// 	Enabled:       true,
		// 	SignalVersion: entities.V2,
		// },
		// {
		// 	Type:          entities.BTC,
		// 	IP:            "api.beta.btc.xeohive.com/test-api/r18-beta",
		// 	Enabled:       true,
		// 	SignalVersion: entities.V2,
		// },
		// {
		// 	Type:          entities.BTC,
		// 	IP:            "api.beta.btc.xeohive.com/test-api/mix-v3-r18",
		// 	Enabled:       true,
		// 	SignalVersion: entities.V2,
		// },
		// {
		// 	Type:          entities.BTC,
		// 	IP:            "api.beta.btc.xeohive.com/test-api/r18-stop-loss-4",
		// 	Enabled:       true,
		// 	SignalVersion: entities.V2,
		// },
		// {
		// 	Type:          entities.ETH,
		// 	IP:            "api.beta.eth.xeohive.com/test-api/r18-stop-loss-4",
		// 	Enabled:       true,
		// 	SignalVersion: entities.V2,
		// },
		// {
		// 	Type:          entities.BTC,
		// 	IP:            "api.beta.btc.xeohive.com/test-api/r18-stop-loss-8",
		// 	Enabled:       true,
		// 	SignalVersion: entities.V2,
		// },
		// {
		// 	Type:          entities.ETH,
		// 	IP:            "api.beta.eth.xeohive.com/test-api/r18-stop-loss-9",
		// 	Enabled:       true,
		// 	SignalVersion: entities.V2,
		// },
		{
			Type:          entities.ETH,
			IP:            "api.beta.eth.xeohive.com/test-api/r18-sopr",
			Enabled:       true,
			SignalVersion: entities.V2,
			TLS:           true,
			TimeFormat:    "2006-01-02 15:04:05",
		},
		{
			Type:          entities.BTC,
			IP:            "209.250.249.95/seasonalitybtc",
			Enabled:       true,
			SignalVersion: entities.V2,
			TLS:           false,
			TimeFormat:    "2006-01-02 15:04:15.999999999+07:00",
		},
		{
			Type:          entities.ETH,
			IP:            "209.250.249.95/seasonalityeth",
			Enabled:       true,
			SignalVersion: entities.V2,
			TLS:           false,
			TimeFormat:    "2006-01-02 15:04:15.999999999Z07:00",
		},
		{
			Type:          entities.BTC,
			IP:            "209.250.249.95:8888/get_liquidation_signals",
			Enabled:       true,
			SignalVersion: entities.V2,
			TLS:           false,
			TimeFormat:    "2006-01-02 15:04:15.999999999Z07:00",
		},
	}

	// setup ticker
	signalTicker, err := tickr.New(
		signalResolutionSeconds,
		signalSources,
		signalService)
	if err != nil {
		log.Fatal("unable to setup btc ticker")
	}

	signalTicker.Run(context.Background(), wl)

	// setup router and handlers
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", Index)
	wl.Debugf("running on port: %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello there and welcome to your service!")
}
