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

type RedisUser struct {
	UserID      string    `json:"user_id"`
	ChatHistory []Message `json:"chat_history"`
	IP          string    `json:"ip"`
}

type ChatHistoryStore struct {
	client *redis.Client
}

// gets User Chat History if exists
func (c *ChatHistoryStore) GetChatHistory(ctx context.Context, clientID string) (map[string]any, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	chatHistory, err := redis_chat_history.New(clientID, 300, c.client)
	if err != nil {
		return nil, err
	}
	// memory buffer only has 4 slots
	memoryBuffer := memory.NewConversationWindowBuffer(4, func(b *memory.ConversationBuffer) {
		b.ChatHistory = chatHistory
	})

	memoryLoad, err := memoryBuffer.LoadMemoryVariables(ctx, map[string]any{})
	if err != nil {
		return nil, err
	}

	return memoryLoad, nil
}

func (client *ChatHistoryStore) PostChatData(ctx context.Context) error {
	return nil
}
