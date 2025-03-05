package store

import (
	"context"

	"github.com/mik-dmi/rag_chatbot/backend/utils/redis_chat_history.go"
	"github.com/redis/go-redis/v9"
	"github.com/tmc/langchaingo/memory"
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

func (c *ChatHistoryStore) CreateChatHistory(ctx context.Context) (*memory.ConversationWindowBuffer, error) {
	// Create a conversation buffer with 4 slots
	memory := memory.NewConversationWindowBuffer(4, func(b *memory.ConversationBuffer) {
		chatHistory, err := redis_chat_history.New(c.client.ClientID(ctx).String(), 300, c.client)
		if err != nil {
			// Handle error appropriately (you may choose to log or propagate this error)
		}
		b.ChatHistory = chatHistory
	})

	return memory, nil
}

func (client *ChatHistoryStore) PostChatData(ctx context.Context) error {
	return nil
}

func (client *ChatHistoryStore) GetChatHistory(ctx context.Context) error {
	return nil
}
