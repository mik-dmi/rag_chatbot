package store

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

type QueriesStore struct {
	client *weaviate.Client
}

func (q *QueriesStore) GetResponse(userQuery string, ctx context.Context) (string, error) {
	return "", nil
}
