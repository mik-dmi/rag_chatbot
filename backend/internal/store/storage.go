package store

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

type Storage struct {
	Queries interface {
		GetResponse(string, context.Context) (string, error)
	}
	Vectors interface {
		CreateVectors([]string, context.Context) error
		DeleteVectors(string, context.Context) error
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
