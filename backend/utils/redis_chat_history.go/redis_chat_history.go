package redis_chat_history

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

type RedisChatMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

// RedisChatMessageHistory implements the schema.ChatMessageHistory interface.
type RedisChatMessageHistory struct {
	sessionID  string
	sessionTTL time.Duration
	client     *redis.Client
}

// New creates a new RedisChatMessageHistory instance.
// sessionTTL is provided in seconds.
func New(sessionID string, sessionTTL int, client *redis.Client) (*RedisChatMessageHistory, error) {
	return &RedisChatMessageHistory{
		sessionID:  sessionID,
		sessionTTL: time.Duration(sessionTTL) * time.Second,
		client:     client,
	}, nil
}

// Ensure RedisChatMessageHistory implements schema.ChatMessageHistory.
var _ schema.ChatMessageHistory = &RedisChatMessageHistory{}

// AddMessage marshals the given llms.ChatMessage and adds it to Redis.
func (h *RedisChatMessageHistory) AddMessage(ctx context.Context, message llms.ChatMessage) error {
	return h.addMessage(ctx, redisChatMessageFromLLMS(message))
}

// AddUserMessage is a convenience method to add a human message.
func (h *RedisChatMessageHistory) AddUserMessage(ctx context.Context, text string) error {
	return h.addMessage(ctx, RedisChatMessage{Type: "human", Content: text})
}

// AddAIMessage is a convenience method to add an AI message.
func (h *RedisChatMessageHistory) AddAIMessage(ctx context.Context, text string) error {
	return h.addMessage(ctx, RedisChatMessage{Type: "ai", Content: text})
}

// Clear removes all messages associated with the session.
func (h *RedisChatMessageHistory) Clear(ctx context.Context) error {
	return h.client.Del(ctx, h.sessionID).Err()
}

// SetMessages clears the current history and adds the provided messages.
func (h *RedisChatMessageHistory) SetMessages(ctx context.Context, messages []llms.ChatMessage) error {
	// Optionally clear existing messages.
	if err := h.Clear(ctx); err != nil {
		return err
	}
	for _, message := range messages {
		if err := h.AddMessage(ctx, message); err != nil {
			return err
		}
	}
	return nil
}

// Messages retrieves all messages stored in Redis.
func (h *RedisChatMessageHistory) Messages(ctx context.Context) ([]llms.ChatMessage, error) {
	msgs, err := h.client.LRange(ctx, h.sessionID, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var chatMessages []llms.ChatMessage
	for _, msgStr := range msgs {
		var rMsg RedisChatMessage
		if err := json.Unmarshal([]byte(msgStr), &rMsg); err != nil {
			return nil, err
		}

		var chatMsg llms.ChatMessage
		switch rMsg.Type {
		case "human":
			chatMsg = llms.HumanChatMessage{Content: rMsg.Content}
		case "ai":
			chatMsg = llms.AIChatMessage{Content: rMsg.Content}
		default:
			// Optionally, handle unknown types.
			continue
		}
		chatMessages = append(chatMessages, chatMsg)
	}
	return chatMessages, nil
}

// addMessage is a helper function to push a RedisChatMessage into the Redis list.
func (h *RedisChatMessageHistory) addMessage(ctx context.Context, message RedisChatMessage) error {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	pipe := h.client.Pipeline()
	pipe.LPush(ctx, h.sessionID, string(msgBytes))
	pipe.Expire(ctx, h.sessionID, h.sessionTTL)
	_, err = pipe.Exec(ctx)
	return err
}

// redisChatMessageFromLLMS converts an llms.ChatMessage to a RedisChatMessage.
func redisChatMessageFromLLMS(message llms.ChatMessage) RedisChatMessage {
	return RedisChatMessage{
		Type:    string(message.GetType()),
		Content: message.GetContent(),
	}
}
