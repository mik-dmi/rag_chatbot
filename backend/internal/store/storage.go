package store

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

type WeaviateStorage struct {
	Vectors interface {
		CreateVectors(context.Context, *RagData) (*VectorCreatedResponse, error)
		GetClosestVectors(context.Context, string) ([]Document, error)

		chapterExists(context.Context, string) (bool, error)
		GetObjectIDByChapter(context.Context, string) (*IDResponse, error)
		DeleteChapterWithChapterName(context.Context, string) (*SuccessfullyDeleted, error)
		DeleteObjectWithID(context.Context, string) (*SuccessfullyDeleted, error)
	}
}
type RedisStorage struct {
	ChatHistory interface {
		CreateSession(context.Context) error
		PostChatData(context.Context) error
		GetChatHistory(context.Context) error
	}
}

func NewWeaviateStorage(client *weaviate.Client) WeaviateStorage {
	return WeaviateStorage{
		Vectors: &VectorsStore{client},
	}
}

func NewRedisStorage(client *reddis.Client) RedisStorage {
	return RedisStorage{
		ChatHistory: &ChatHistoryStore{client},
	}
}
