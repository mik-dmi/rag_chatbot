package store

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

type UsersStore struct {
	client *weaviate.Client
}

func (s *UsersStore) CreateSession(ctx context.Context) error {
	return nil
}

func (s *UsersStore) GetChatHistory(ctx context.Context) error {
	return nil
}
