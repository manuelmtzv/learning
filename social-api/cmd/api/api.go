package main

import (
	"log"
	"net/http"
	"time"
)

type application struct {
	config config
}

type config struct {
	addr string
}

func (app *application) mount() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/health", app.healthCheckHandler)

	return mux
}

func (app *application) run(mux *http.ServeMux) error {

	server := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 45,
		ReadTimeout:  time.Second * 20,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at %s", app.config.addr)

	return server.ListenAndServe()
}
