package main

import (
	"context"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/prompts"
)

func (app *application) standaloneQuestion(memoryLoad map[string]any, questionUser string) (string, error) {

	standalonePrompt := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		standalonePromptTemplate,
		prompts.NewHumanMessagePromptTemplate(
			`Chat History: {{.chat_history}}
			Follow-up Question: {{.question}}
			Independent Question:`,
			[]string{"chat_history", "question"},
		)})

	standaloneChain := chains.NewLLMChain(app.openaiClients.standaloneChainClient, standalonePrompt)

	input := map[string]any{
		"chat_history": memoryLoad["chat_history"],
		"question":     questionUser,
	}

	res, err := chains.Run(context.Background(), standaloneChain, input)
	if err != nil {
		return "", nil
	}

	return res, nil

}
