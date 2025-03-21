package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/mik-dmi/rag_chatbot/backend/internal/store"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/prompts"
)

type CreateDocumentsPayload struct {
	UserID    string           `json:"user_id"`
	Documents []store.Document `json:"document"`
}
type UserQuery struct {
	UserID      string `json:"user_id"`
	UserMessage string `json:"user_message"`
}

func (app *application) createVectorHandler(w http.ResponseWriter, r *http.Request) {
	var documents CreateDocumentsPayload
	if err := readJSON(w, r, &documents); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	userId := "1"

	var aggregatedDocuments []store.Document

	for _, doc := range documents.Documents {
		// Create a slice to hold all subsections for the current chapter.
		var subsections []store.Subsection
		for _, subsection := range doc.Subsections {
			subsections = append(subsections, store.Subsection{
				Title:   subsection.Title,
				Content: subsection.Content,
			})
		}

		// Append the aggregated document.
		aggregatedDocuments = append(aggregatedDocuments, store.Document{
			Chapter:     doc.Chapter,
			Subsections: subsections,
		})
	}

	vector := &store.RagData{
		UserID:    userId,
		Documents: aggregatedDocuments,
	}

	ctx := r.Context()

	vectorsCreated, err := app.weaviateStore.Vectors.CreateVectors(ctx, vector)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, vectorsCreated); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) userQuestionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var query UserQuery

	if err := readJSON(w, r, &query); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	similarDocs, err := app.weaviateStore.Vectors.GetClosestVectors(ctx, query.UserMessage)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	// put all the docs in a array of string
	// may in ht future change the GetClosestVectors return value
	var docs []string
	for _, doc := range similarDocs {
		jsonData, err := json.Marshal(doc)
		if err != nil {
			log.Printf("Error marshaling document: %v", err)
			continue
		}
		docs = append(docs, string(jsonData))
	}

	// it will be changed in the future
	uniqueUserID := r.Header.Get("X-User-ID")

	memory, err := app.redisStore.ChatHistory.GetChatHistory(ctx, uniqueUserID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Normalize the user's question (if needed)
	questionUser := strings.ReplaceAll(strings.TrimSpace(query.UserMessage), "\n", " ")

	//check if chat_history exists in redis, if it does the users question and history are to make a standalone question
	if chatHist, ok := memory["chat_history"].(string); ok && chatHist != "" {
		// If there is chat history, create a standalone question based on history
		questionUser, err = app.standaloneQuestion(memory, questionUser)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
	}

	log.Println("Question used for the main chain ", questionUser)

	finalPrompt := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		finalPromptTemplate,
		prompts.NewHumanMessagePromptTemplate(
			`CHAT HISTORY: {{.chat_history}}
			CONTEXT: {{.context}}
			Question:{{.question}}`,
			[]string{"chat_history", "context", "question"},
		)})

	finalChain := chains.NewLLMChain(app.openaiClients.mainChainClient, finalPrompt)

	input := map[string]any{
		"chat_history": memory["chat_history"],
		"context":      strings.Join(docs, "\n"),
		"question":     questionUser,
	}
	formattedPrompt, err := finalPrompt.Format(input)
	if err != nil {
		log.Printf("Error formatting prompt: %v", err)
	} else {
		log.Printf("Final Rendered Prompt:\n%s", formattedPrompt)
	}

	finalRagAnswer, err := chains.Call(ctx, finalChain, input)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, finalRagAnswer); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type GetChapterNameIDBody struct {
	ChapterName string `json:"chapter_name"`
}

func (app *application) getObjectIDByChapterHandler(w http.ResponseWriter, r *http.Request) {
	var chapterName GetChapterNameIDBody
	if err := readJSON(w, r, &chapterName); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	ctx := r.Context()
	objectIDRetrieved, err := app.weaviateStore.Vectors.GetObjectIDByChapter(ctx, chapterName.ChapterName)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := writeJSON(w, http.StatusOK, objectIDRetrieved); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deleteVectorObjectByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Print(r.PathValue(id))

	ctx := r.Context()
	response, err := app.weaviateStore.Vectors.DeleteObjectWithID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)

		}
		return
	}

	if err := writeJSON(w, http.StatusOK, response); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
