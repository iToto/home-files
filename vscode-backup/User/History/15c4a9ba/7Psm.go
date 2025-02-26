package main

import (
	"database/sql"
	"yield-mvp/wlog"
	// "encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

var db *sql.DB

type Foo struct {
	Bar string
}

// func DBConnection() *sql.DB {
// 	dbURL := os.Getenv("DATABASE_URL")

// 	if dbURL == "" {
// 		log.Fatal("$DATABASE_URL environment variable must be set")
// 	} else {
// 		fmt.Println("Connected to database: " + dbURL)
// 	}

// 	db, err := sql.Open("postgres", dbURL)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return db
// }

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT environment variable must be set")
	} else {
		fmt.Println("Running on port: " + port)
	}

	wl, err := wlog.NewBasicLogger()
	if err != nil {
		return fmt.Errorf("error configuring logger")
	}

	// db = DBConnection()
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", Index)
	router.HandleFunc("/foo", FooIndex)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello there and welcome to your service!")
}

func FooIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// jsonString, _ := json.Marshal(foos)
	jsonString := []byte(`{"foo":"bar"}`)

	w.Write(jsonString)
}
