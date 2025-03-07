package llm

import (
	"github.com/tmc/langchaingo/llms/openai"
)

func NewOpenaiClient(token string, llmModel string) (*openai.LLM, error) {

	openaiClient, err := openai.New(
		openai.WithToken(token),
		openai.WithModel(llmModel),
	)

	if err != nil {
		return nil, err
	}

	return openaiClient, nil

}
