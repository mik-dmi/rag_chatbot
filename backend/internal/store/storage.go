package store

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

type Storage struct {
	Queries interface {
		GetRAGResponse(context.Context, string) (string, error)
	}
	Vectors interface {
		CreateVectors(context.Context, []string) error
		DeleteVectors(context.Context, string) error
	}
	Users interface {
		CreateSession(context.Context) error
		GetChatHistory(context.Context) error
	}
}

func NewWeaviateStorage(client *weaviate.Client) Storage {
	return Storage{
		Vectors: &VectorsStore{client},
		Users:   &UsersStore{client},
		Queries: &QueriesStore{client},
	}
}
