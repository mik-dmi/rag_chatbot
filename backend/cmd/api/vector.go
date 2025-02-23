package main

import (
	"fmt"
	"net/http"

	"github.com/mik-dmi/rag_chatbot/backend/internal/store"
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

func (app *application) userQueryHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := writeJSON(w, http.StatusCreated, result); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

}

func (app *application) getVectorObjectByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Print(id)

	//vectorObject, err := app.store.Vectors. (ctx, query.UserMessage)

}
func (app *application) deleteVectorObjectByIdHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print(r.PathValue("id"))

}
