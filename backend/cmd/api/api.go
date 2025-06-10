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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mik-dmi/rag_chatbot/backend/internal/auth"
	"github.com/mik-dmi/rag_chatbot/backend/internal/store"
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
	mail               mailConfig
}

type mailConfig struct {
	exp time.Duration
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

func (app *application) mount() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/v1/authentication", func(r chi.Router) {
		router.Post("/jwt-token-auth", app.jwtTokenHandler)
		r.Post("/user", app.registerUserHandler)

	})

	router.Route("/v1", func(r chi.Router) {

		r.Use(app.AuthTokenMiddleware)

		r.Post("/vector-db", app.createVectorHandler)
		r.Get("/vector-db/object", app.getObjectIDByChapterHandler)
		r.Delete("/vector-db/object/{id}", app.deleteVectorObjectByIdHandler)
		r.Patch("/vector-db/object/{id}", app.updateVectorObjectByIdHandler)

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/", app.getUserHandler)

				r.Post("/query", app.userQuestionHandler)
				//r.Post("/create-user", app.createUserHandler)

			})

		})
	})

	return router

}

func (app *application) Run(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
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
