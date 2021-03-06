package main

import (
	"io"
	"log"
	"net/http"

	httplogger "github.com/gleicon/go-httplogger"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

func main() {
	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/", helloHandler)
	srv := http.Server{
		Addr:    ":8080",
		Handler: httplogger.HTTPLogger(serveMux),
	}
	log.Fatal(srv.ListenAndServe())
}
