package main

import (
	"boilerplate-go-api/cmd/my-app/internal/handler"
	"boilerplate-go-api/internal/service/hellosvc"
	"boilerplate-go-api/internal/wlog"
	"boilerplate-go-api/pkg/utils"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

	// setup services
	helloService, err := hellosvc.New()
	if err != nil {
		log.Fatal("unable to init hello service: %w", err)
	}

	// setup router and handlers
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", Index)
	router.HandleFunc(
		"/hello",
		handler.HelloWorld(wl, helloService)).
		Methods(http.MethodGet, http.MethodOptions)
	wl.Debugf("running on port: %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello there and welcome to your service!")
}
