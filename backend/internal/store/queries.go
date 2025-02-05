package store

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

type Query struct {
	UserID   string `json:"user_id"`
	Question string `json:"question"`
	SendAt   string `json:"send_at"`
}

type QueriesStore struct {
	client *weaviate.Client
}

func (q *QueriesStore) GetRAGResponse(ctx context.Context, userQuery string) (string, error) {
	return "", nil
}
