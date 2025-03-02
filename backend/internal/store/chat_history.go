package store

import (
	"context"
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

type ChatHistoryStore struct {
	client *redis.Client
}

func (s *ChatHistoryStore) CreateSession(ctx context.Context) error {
	return nil
}

func (s *ChatHistoryStore) PostChatData(ctx context.Context) error {
	return nil
}

func (s *ChatHistoryStore) GetChatHistory(ctx context.Context) error {
	return nil
}
