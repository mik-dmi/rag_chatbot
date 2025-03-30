package store

import (
	"context"

	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

var (
	ErrNotFound             = errors.New("vector not found")
	ErrChapterAlreadyExists = errors.New("already exists in weaviate")
)

type WeaviateStorage struct {
	Vectors interface {
		CreateVectors(context.Context, *RagData) (*VectorCreatedResponse, error)
		GetClosestVectors(context.Context, string) ([]*Document, error)

		chapterExists(context.Context, string) (bool, error)
		GetObjectIDByChapter(context.Context, string) (*IDResponse, error)
		DeleteChapterWithChapterName(context.Context, string) (*SuccessfullyAPIOperation, error)
		DeleteObjectWithID(context.Context, string) (*SuccessfullyAPIOperation, error)
		UpdateObjectWithID(context.Context, Document, string) (*SuccessfullyAPIOperation, error)
	}
}
type RedisStorage struct {
	ChatHistory interface {
		GetChatHistory(context.Context, string) (map[string]any, error)
		PostChatData(context.Context) error
	}
}

func NewWeaviateStorage(client *weaviate.Client) WeaviateStorage {
	return WeaviateStorage{
		Vectors: &VectorsStore{client},
	}
}

func NewRedisStorage(client *redis.Client) RedisStorage {
	return RedisStorage{
		ChatHistory: &ChatHistoryStore{client},
	}
}
