package main

import (
	"log"

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

	log.Println("Standalond prompt: ", standalonePrompt)

	return "", nil

}
