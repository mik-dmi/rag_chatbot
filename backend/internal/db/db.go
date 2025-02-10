package db

import (
	"context"
	"fmt"
	"time"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

func NewWeaviateClient(host string, addr string) (*weaviate.Client, error) {
	config := weaviate.Config{
		Host:   fmt.Sprintf("%s%s", host, addr),
		Scheme: "http",
	}
	client, err := weaviate.New(config)
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
	fmt.Println("weaviate is running at ", addr)
	return client, nil

}
