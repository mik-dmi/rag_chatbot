package store

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

type Storage struct {
	Vectors interface {
		CreateVectors(context.Context, *RagData) (*VectorCreatedResponse, error)
		GetClosestVectors(context.Context, string) ([]Document, error)
		DeleteChapterVectors(context.Context, string) error
		chapterExists(context.Context, string) (bool, error)
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
	}
}
