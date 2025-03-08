package llm

import (
	"github.com/tmc/langchaingo/llms/openai"
)

func NewOpenaiClient(standaloneChainLLMToken string, mainChainLLMToken string, standaloneChainLLMModel string, mainChainLLMModel string) (*openai.LLM, *openai.LLM, error) {

	standaloneChainOpenaiClient, err := openai.New(
		openai.WithToken(standaloneChainLLMToken),
		openai.WithModel(standaloneChainLLMModel),
	)
	mainChainOpenaiClient, err := openai.New(
		openai.WithToken(mainChainLLMToken),
		openai.WithModel(mainChainLLMModel),
	)

	if err != nil {
		return nil, nil, err
	}

	return standaloneChainOpenaiClient, mainChainOpenaiClient, nil

}
