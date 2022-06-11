package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"


	"git.eletrotupi.com/momento/api"
	"git.eletrotupi.com/momento/database"
)

func main() {
	// TODO: Allow to configure bind address and the connection string to be
	// used

	db, err := sql.Open("postgres", os.Args[1])
	if err != nil {
		log.Fatalf("Failed to open a database connection: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/api/", api.New())
	serve := database.Middleware(db)(mux)
	serve = WithLogging(serve)

	server := &http.Server{
		Addr:		":8000",
		Handler:	serve,
	}

	log.Printf("Listening on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error: Couldn't serve %s", err.Error())
	}
}
