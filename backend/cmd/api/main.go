package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/mik-dmi/rag_chatbot/backend/internal/env"
	"github.com/mik-dmi/rag_chatbot/backend/internal/store"
	langchain_weaviate "github.com/tmc/langchaingo/vectorstores/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

func GetWeaviateClient(_ context.Context, cfg config) (*weaviate.Client, error) {
	config := weaviate.Config{
		Host:   fmt.Sprintf("%s%s", cfg.vectorDB.host, cfg.vectorDB.addr),
		Scheme: "http",
	}
	client, err := weaviate.NewClient(config)
	if err != nil {
		return nil, err
	}
	return client, nil

}

func setupWeaviate(_ context.Context, cfg config) (any, error) {

	/*openaiClient, err := openai.New(
		openai.WithModel("gpt-3.5-turbo-0125"),
		openai.WithEmbeddingModel(embeddingModelName),
	)
	if err != nil {
		return nil, err
	}
	emb, err := embeddings.NewEmbedder(openaiClient)
	if err != nil {
		return nil, err
	}*/
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

func main() {

	err := godotenv.Load()

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		vectorDB: vectorDBConfig{
			addr: env.GetString("VECTOR_DB_PORT", "8080"),
			host: env.GetString("VECTOR_DB_HOST", "http://localhost"),
		},
	}

	weaviateClient, err := GetWeaviateClient(nil, cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	store := store.NewWeaviateStorage(weaviateClient)
	app := &application{
		config: cfg,
		store:  store,
	}
	mux := app.mount()
	log.Fatal(app.Run(mux))
}
