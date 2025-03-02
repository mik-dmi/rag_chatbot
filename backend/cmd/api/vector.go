package main

import (
	"fmt"
	"net/http"

	"github.com/mik-dmi/rag_chatbot/backend/internal/store"
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
		writeJSONError(w, http.StatusBadRequest, err.Error())
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

	vectorsCreated, err := app.store.Vectors.CreateVectors(ctx, vector)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusCreated, vectorsCreated); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

}

func (app *application) userQuestionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var query UserQuery

	if err := readJSON(w, r, &query); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := app.store.Vectors.GetClosestVectors(ctx, query.UserMessage)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var prompt = prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		promptTemplate,
		prompts.NewHumanMessagePromptTemplate(
			`CHAT HISTORY: {{.chat_history}}
	CONTEXT: {{.context}}
	QUESTION: {{.question}}`,
			[]string{"chat_history", "context", "question"},
		),
	})

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusOK, result); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

}

type GetChapterNameIDBody struct {
	ChapterName string `json:"chapter_name"`
}

func (app *application) getObjectIDByChapterHandler(w http.ResponseWriter, r *http.Request) {
	var chapterName GetChapterNameIDBody
	if err := readJSON(w, r, &chapterName); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	ctx := r.Context()
	objectIDRetrieved, err := app.store.Vectors.GetObjectIDByChapter(ctx, chapterName.ChapterName)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := writeJSON(w, http.StatusOK, objectIDRetrieved); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

}

func (app *application) deleteVectorObjectByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Print(r.PathValue(id))

	ctx := r.Context()
	response, err := app.store.Vectors.DeleteObjectWithID(ctx, id)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusOK, response); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
