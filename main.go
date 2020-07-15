package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"slugger/db"
	"syscall"
	"time"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

func main() {

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()
	sm.HandleFunc("/{slug}", DecryptSlug)
	// public api no auth
	apiRoutesPublic := sm.PathPrefix("/api/v1").Subrouter()
	apiRoutesPublic.HandleFunc("/slug/new", urlShortener).Methods("GET")

	// CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	// create a new server
	s := http.Server{
		Addr:         ":9001",           // configure the bind address
		Handler:      ch(sm),            // set the default handler
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		log.Println("Starting orbiter service on port 9001")

		err := s.ListenAndServe()
		if err != nil {
			log.Fatal("Error starting orbiter service: %s\n", err)
			os.Exit(1)
		}
	}()

	sigs := make(chan os.Signal, 1)
	sigdone := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Signals
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		sigdone <- true
	}()

	<-sigdone
	log.Println("Got signal, Exiting orbiter service")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}

// DecryptSlug - Check if slug exists in db
func DecryptSlug(w http.ResponseWriter, r *http.Request) {
	Connections := db.DBConnection()
	mysql := Connections["mysql"].(db.Mysql)
	mysqlConn := mysql.Connect()
	client := mysqlConn.MYSQL
	defer client.Close()

	vars := mux.Vars(r)
	slug := vars["slug"]

	query := `SELECT url.original_url FROM url WHERE slug = ? `
	rows, err := client.Query(query, slug)

	if err != nil {
		log.Println(err.Error())
		return
	}

	var originalURL string
	for rows.Next() {

		err = rows.Scan(&originalURL)

		if err != nil {
			log.Println(err.Error())
		}

	}

	if originalURL == "" {
		fmt.Fprintf(w, "bad slug")
	} else {
		fmt.Fprintf(w, originalURL)
	}

}

// urlShortener - Shorten url
func urlShortener(w http.ResponseWriter, r *http.Request) {
	Connections := db.DBConnection()
	mysql := Connections["mysql"].(db.Mysql)
	mysqlConn := mysql.Connect()
	client := mysqlConn.MYSQL
	defer client.Close()

	slug := RandStringRunes(6)

	for {

		if !checkQuery(client, slug) {
			break
		}

	}

	a := r.URL.Query()

	query := `INSERT INTO url (original_url, slug) VALUES (?, ?);`
	stmt, err := client.Prepare(query)

	if err != nil {
		log.Println(err.Error())
		return
	}

	_, err = stmt.Exec(a.Get("original_url"), slug)
	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Fprintf(w, r.Host+"/"+slug)
}

func checkQuery(client *sql.DB, slug string) bool {

	queryCheck := `SELECT url.id 
	FROM url 
	WHERE url.slug = ?`

	rows, err := client.Query(queryCheck, slug)

	if err != nil {
		log.Println(err.Error())
		return false
	}

	temp := 0

	for rows.Next() {

		err = rows.Scan(&temp)
		log.Println(temp)
		if err != nil {
			log.Println(err.Error())
			return false
		}

	}

	if temp == 0 {
		return false
	}
	return true

}

func RandStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
