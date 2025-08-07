package main

import (
	"database/sql"
	"embed"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	var (
		port        int
		databaseURL string
	)
	flag.IntVar(&port, "port", 8080, "HTTP server port")
	flag.StringVar(&databaseURL, "database-url", "postgres://postgres:@localhost:5432/postgres?sslmode=disable", "Database URL")
	flag.Parse()

	var repo *EntryRepo

	log.Print("Initializing database connection")
	if db, err := sql.Open("postgres", databaseURL); err == nil {
		repo = &EntryRepo{db}
		defer db.Close()
	} else {
		log.Fatal(err)
	}

	var (
		listHandler = &ListHandler{repo}
		formHandler = &FormHandler{repo}
		fileHandler = http.FileServer(http.FS(staticFiles))
	)

	http.Handle("GET /{$}", listHandler)
	http.Handle("POST /{$}", formHandler)
	http.Handle("GET /static/", fileHandler)

	log.Printf("Starting web server on http://localhost:%d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}
}
