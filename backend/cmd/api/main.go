package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mik-dmi/rag_chatbot/backend/internal/db"
	"github.com/mik-dmi/rag_chatbot/backend/internal/env"
	"github.com/mik-dmi/rag_chatbot/backend/internal/llm"
	"github.com/mik-dmi/rag_chatbot/backend/internal/store"
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
		standaloneLLMModel: llmConfig{
			token: env.GetString("OPEN_AI_SECRET", "openai_key"),
			model: "gpt-4",
		},
		mainLLMModel: llmConfig{
			token: env.GetString("OPEN_AI_SECRET", "openai_key"),
			model: "gpt-4",
		},

		env: env.GetString("ENV", "development"),
	}

	weaviateClient, err := db.NewWeaviateClient(cfg.weaviateDB.host, cfg.weaviateDB.addr)
	if err != nil {
		log.Fatal(err)
	}

	redisClient, err := db.NewRedisClient(cfg.redisDB.host, cfg.redisDB.addr, cfg.redisDB.password)
	if err != nil {
		log.Fatal(err)
	}

	standaloneChainOpenaiClient, mainChainOpenaiClient, err := llm.NewOpenaiClient(cfg.standaloneLLMModel.token, cfg.mainLLMModel.token, cfg.standaloneLLMModel.model, cfg.mainLLMModel.model)
	if err != nil {
		log.Fatal(err)
	}

	weaviateStore := store.NewWeaviateStorage(weaviateClient)
	redisStore := store.NewRedisStorage(redisClient)
	app := &application{
		config:        cfg,
		weaviateStore: weaviateStore,
		redisStore:    redisStore,
		openaiClients: OpenaiClients{
			standaloneChainClient: standaloneChainOpenaiClient,
			mainChainClient:       mainChainOpenaiClient,
		},
	}
	mux := app.mount()
	log.Fatal(app.Run(mux))
}
