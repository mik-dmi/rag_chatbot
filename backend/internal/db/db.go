package db

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

func NewWeaviateClient(host string, addr string) (*weaviate.Client, error) {
	config := weaviate.Config{
		Host:   fmt.Sprintf("%s%s", host, addr),
		Scheme: "http",
	}
	client, err := weaviate.NewClient(config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	liveStatus, err := client.Misc().LiveChecker().Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping Weaviate: %w", err)
	}
	if !liveStatus {
		return nil, fmt.Errorf("weaviate is not live: status=%t", liveStatus)
	}

	// creating a Book Class if it does not exist yet
	classObj := &models.Class{
		Class:      "Book",
		Vectorizer: "text2vec-transformers",
	}
	if err := client.Schema().ClassCreator().WithClass(classObj).Do(ctx); err != nil {
		// Depending on your needs, you might want to check if the class already exists
		// and ignore the error if so.
		fmt.Println("failed to create class:", err)
	}

	fmt.Println("weaviate is running at ", addr)
	return client, nil

}

func NewRedisClient(host string, addr string, password string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s%s", host, addr),
		Password: password,
		DB:       0, // Use default DB
		//Protocol: 2, // Connection protocol
	})

	// ping the db
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}
	fmt.Println("Redis is running at:", addr)
	return client, nil

}
