package main

import (
	"log"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/prompts"
)

func (app *application) standaloneChain(memoryLoad map[string]any, questionUser string) (*chains.LLMChain, error) {

	standalonePrompt := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		standalonePromptTemplate,
		prompts.NewHumanMessagePromptTemplate(
			`Chat History: {{.chat_history}}
			Follow-up Question: {{.question}}
			Independent Question:`,
			[]string{"chat_history", "question"},
		)})

	standaloneChain := chains.NewLLMChain(app.openaiClients.standaloneChainClient, standalonePrompt)

	standaloneChain.OutputKey
	log.Println("Standalond prompt: ", standalonePrompt)

	return "", nil

}
