package main

import (
	"database/sql"
	// "encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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

	// db = DBConnection()
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", Index)
	router.HandleFunc("/foo", FooIndex)
	router.HandleFunc("/foo/{id}", GetFoo).Methods("GET")
	router.HandleFunc("/foo/{id}", PostFoo).Methods("POST")
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

func GetFoo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 0)

	var foo string

	err := db.QueryRow("SELECT * FROM foo WHERE id = $1", id).Scan(&foo)

	if err != nil {
		log.Fatal(err)
	}

	// fmt.Print("foo: ", foo)
	fmt.Fprintln(w, "Foo show:", foo)
}

func PostFoo(w http.ResponseWriter, r *http.Request) {
	// TODO
}
