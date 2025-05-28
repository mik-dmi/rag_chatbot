package main

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mik-dmi/rag_chatbot/backend/internal/auth"
	"github.com/mik-dmi/rag_chatbot/backend/internal/db"
	"github.com/mik-dmi/rag_chatbot/backend/internal/env"
	"github.com/mik-dmi/rag_chatbot/backend/internal/llm"
	"github.com/mik-dmi/rag_chatbot/backend/internal/store"
	lg "github.com/mik-dmi/rag_chatbot/backend/utils/logger"
	"go.uber.org/zap"
)

/*
func setupWeaviate(_ context.Context, cfg config) (any, error) {
	openaiClient, err := openai.New(
		openai.WithModel("gpt-3.5-turbo-0125"),
		openai.WithEmbeddingModel(embeddingModelName),
	)
	if err != nil {
		return nil, err
	}
	emb, err := embeddings.NewEmbedder(openaiClient)
	if err != nil {
		return nil, err
	}
	wvStore, err := langchain_weaviate.New(
		//langchain_weaviate.WithEmbedder(emb),
		langchain_weaviate.WithScheme("http"),
		langchain_weaviate.WithHost("8080"),
		langchain_weaviate.WithIndexName("Document"),
	)
	if err != nil {
		return nil, err
	}
	return &wvStore, nil
}
*/

const version = "0.0.1"

func main() {

	err := godotenv.Load() //"../../.env"
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		weaviateDB: weaviateDBConfig{
			addr: env.GetString("WEAVIATE_DB_PORT", ":8080"),
			host: env.GetString("WEAVIATE_DB_HOST", "localhost"),
		},
		redisDB: redisDBConfig{
			addr:     env.GetString("REDIS_DB_PORT", ":6379"),
			host:     env.GetString("REDIS_DB_HOST", "localhost"),
			password: env.GetString("REDIS_PASSWORD", "redis_password"),
		},
		postgresDB: postgresDBConfig{
			addr:         env.GetString("POSTGRES_ADDR", "postgres://admin:adminpassword@localhost:5492/postgres_rag?sslmode=disable"),
			maxOpenConns: env.GetInt("POSTGRES_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_CONNS", "15m"),
		},
		standaloneLLMModel: llmConfig{
			token: env.GetString("OPEN_AI_SECRET", "openai_key"),
			model: "gpt-3.5-turbo",
		},
		mainLLMModel: llmConfig{
			token: env.GetString("OPEN_AI_SECRET", "openai_key"),
			model: "gpt-3.5-turbo",
		},
		authCredencials: authConfig{
			authCredencials: authCredencialsConfig{
				clientID: env.GetString("CLIENT_ID", "test_id"),
				password: env.GetString("PASSWORD_CLIENT", "12345_test_client_password_!"),
			},
			token: tokenConfig{
				secret: env.GetString("SECRET", "secret_test"),
				exp:    time.Hour * 2,
				iss:    env.GetString("ISS", "rag_system"),
			},
		},

		env: env.GetString("ENV", "development"),
	}
	var logger *zap.SugaredLogger

	if env.GetString(cfg.env, "development") == "production" {
		logger, err = lg.NewProductionLogger()
		if err != nil {
			log.Fatalf("Failed to initialize logger: %s", err)
		}
	} else {
		logger, err = lg.NewDevelopmentLogger()
		if err != nil {
			log.Fatalf("Failed to initialize logger: %s", err)
		}
	}

	weaviateClient, err := db.NewWeaviateClient(cfg.weaviateDB.host, cfg.weaviateDB.addr)
	if err != nil {
		log.Fatal(err)
	}

	redisClient, err := db.NewRedisClient(cfg.redisDB.host, cfg.redisDB.addr, cfg.redisDB.password)
	if err != nil {
		log.Fatal(err)
	}
	postgreClient, err := db.NewPostgreClient(cfg.postgresDB.addr, cfg.postgresDB.maxOpenConns, cfg.postgresDB.maxIdleConns, cfg.postgresDB.maxIdleTime)
	if err != nil {
		log.Fatal(err)
	}

	standaloneChainOpenaiClient, mainChainOpenaiClient, err := llm.NewOpenaiClient(cfg.standaloneLLMModel.token, cfg.mainLLMModel.token, cfg.standaloneLLMModel.model, cfg.mainLLMModel.model)
	if err != nil {
		log.Fatal(err)
	}
	tokenHost := "rag_system"

	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.authCredencials.token.secret, tokenHost, tokenHost)

	weaviateStore := store.NewWeaviateStorage(weaviateClient)
	redisStore := store.NewRedisStorage(redisClient)
	postgreStore := store.NewPostgreStorage(postgreClient)
	app := &application{
		config:        cfg,
		weaviateStore: weaviateStore,
		redisStore:    redisStore,
		postgreStore:  postgreStore,
		openaiClients: OpenaiClients{
			standaloneChainClient: standaloneChainOpenaiClient,
			mainChainClient:       mainChainOpenaiClient,
		},
		logger:        logger,
		authenticator: jwtAuthenticator,
	}
	mux := app.mount()
	log.Fatal(app.Run(mux))
}
