package store

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

type Message struct {
	Time          string `json:"time"`
	Content       string `json:"content"`
	IsUserMessage string `json:"is_user_message"`
}

type User struct {
	UserID      string    `json:"user_id"`
	ChatHistory []Message `json:"chat_history"`
	IP          string    `json:"ip"`
}

type UsersStore struct {
	client *weaviate.Client
}

func (s *UsersStore) CreateSession(ctx context.Context) error {
	return nil
}

func (s *UsersStore) GetChatHistory(ctx context.Context) error {
	return nil
}
