package db

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

func NewWeaviateClient(_ context.Context, host string, addr string) (*weaviate.Client, error) {
	config := weaviate.Config{
		Host:   fmt.Sprintf("%s%s", host, addr),
		Scheme: "http",
	}
	client, err := weaviate.NewClient(config)
	if err != nil {
		return nil, err
	}
	return client, nil

}
