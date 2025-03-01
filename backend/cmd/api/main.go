package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mik-dmi/rag_chatbot/backend/internal/db"
	"github.com/mik-dmi/rag_chatbot/backend/internal/env"
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
		vectorDB: vectorDBConfig{
			addr: env.GetString("VECTOR_DB_PORT", "8080"),
			host: env.GetString("VECTOR_DB_HOST", "http://localhost"),
		},
		env: env.GetString("ENV", "development"),
	}

	weaviateClient, err := db.NewWeaviateClient(cfg.vectorDB.host, cfg.vectorDB.addr)
	if err != nil {
		log.Fatal(err)
	}

	redisClient, err := db.NewWeaviateClient(cfg.vectorDB.host, cfg.vectorDB.addr)
	if err != nil {
		log.Fatal(err)
	}

	weaviateStore := store.NewWeaviateStorage(weaviateClient)
	redisStore := store.RedisStorage(redisClient)
	app := &application{
		config:        cfg,
		weaviateStore: weaviateStore,
		redisStore:    redisStore,
	}
	mux := app.mount()
	log.Fatal(app.Run(mux))
}
