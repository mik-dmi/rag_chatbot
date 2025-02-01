package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/mik-dmi/rag_chatbot/backend/internal/env"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/vectorstores/weaviate"
)

func GetSchema() {
	cfg := weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	schema, err := client.Schema().Getter().Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(schema)
}

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
	wvStore, err := weaviate.New(
		weaviate.WithEmbedder(emb),
		weaviate.WithScheme("http"),
		weaviate.WithHost("8080"),
		weaviate.WithIndexName("Document"),
	)

	if err != nil {
		return nil, err
	}

	return &wvStore, nil

}

func main() {

	err := godotenv.Load()

	weaviate, err := setupWeaviate(ctx, cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to setup weave")
	}
	GetSchema()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cfg := config{
		addr: env.GetString("ADDR", ":8000"),
	}
	app := &application{
		config: cfg,
	}
	mux := app.mount()

	log.Fatal(app.Run(mux))

}
