package main

import (
	"fmt"
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

	mux.HandleFunc(" /v1/ragchatbot", app.healthCheckHandler)
	return mux

}

func (app *application) Run(mux *http.ServeMux) error {

	svr := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	fmt.Printf("server is running at %s", app.config.addr)
	return svr.ListenAndServe()

}
