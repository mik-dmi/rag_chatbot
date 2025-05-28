package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mik-dmi/rag_chatbot/backend/internal/auth"
	"github.com/mik-dmi/rag_chatbot/backend/internal/store"
	"github.com/mik-dmi/rag_chatbot/backend/utils/middleware"
	"github.com/tmc/langchaingo/llms/openai"
	"go.uber.org/zap"
)

type application struct {
	config        config
	weaviateStore store.WeaviateStorage
	redisStore    store.RedisStorage
	postgreStore  store.PostgreStorage
	openaiClients OpenaiClients
	logger        *zap.SugaredLogger
	authenticator auth.Authenticator
}
type OpenaiClients struct {
	standaloneChainClient *openai.LLM
	mainChainClient       *openai.LLM
}
type config struct {
	addr               string
	weaviateDB         weaviateDBConfig
	redisDB            redisDBConfig
	postgresDB         postgresDBConfig
	env                string
	standaloneLLMModel llmConfig
	mainLLMModel       llmConfig
	authCredencials    authConfig
}

type authConfig struct {
	authCredencials authCredencialsConfig
	token           tokenConfig
}
type authCredencialsConfig struct {
	clientID string
	password string
}
type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type llmConfig struct {
	token string
	model string
}
type weaviateDBConfig struct {
	addr string
	host string
}

type postgresDBConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}
type redisDBConfig struct {
	addr     string
	host     string
	password string
}

func (app *application) mount() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /jwt-token-auth", app.jwtTokenHandler)
	router.HandleFunc("POST /query", app.userQuestionHandler)
	router.HandleFunc("POST /vector-db", app.createVectorHandler)
	router.HandleFunc("GET /vector-db/object", app.getObjectIDByChapterHandler)
	router.HandleFunc("DELETE /vector-db/object/{id}", app.deleteVectorObjectByIdHandler)
	router.HandleFunc("PATCH /vector-db/object/{id}", app.updateVectorObjectByIdHandler)

	router.HandleFunc("GET /{userId}", app.getUserHandler)

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

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      stack(mux),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Infow("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Infow("server has started", "addr", app.config.addr, "env", app.config.env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.logger.Infow("server has stopped", "addr", app.config.addr, "env", app.config.env)

	return nil

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
