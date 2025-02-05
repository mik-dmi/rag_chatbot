package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/mik-dmi/rag_chatbot/backend/internal/store"
	middleware "github.com/mik-dmi/rag_chatbot/backend/utils"
)

type application struct {
	config config
	store  store.Storage
}
type config struct {
	addr     string
	vectorDB vectorDBConfig
}

type vectorDBConfig struct {
	addr string
	host string
}

func (app *application) mount() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /query", app.healthCheckHandler)
	router.HandleFunc("POST /vector-db", app.healthCheckHandler)

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", router)) //dealing with subrouting
	return v1

}

func (app *application) Run(mux *http.ServeMux) error {

	stack := middleware.CreateStack(
		middleware.Timeout(90*time.Second),
		middleware.Recovery, // Added first (but runs last)
		middleware.Logging,
	)

	svr := &http.Server{
		Addr:         app.config.addr,
		Handler:      stack(mux),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	fmt.Println("server is running at", app.config.addr)
	return svr.ListenAndServe()

}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		logger.Info("Request Info",
			slog.String("method", r.Method),
			slog.String("path", r.RequestURI),
			slog.String("host", r.Host),
		)
		next.ServeHTTP(w, r)
	})
}
