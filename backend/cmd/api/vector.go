package main

import (
	"net/http"

	"github.com/mik-dmi/rag_chatbot/backend/internal/store"
)

type CreateDocumentsPayload struct {
	UserID    string           `json:"userid"`
	Documents []store.Document `json:"document"`
}

func (app *application) createVectorHandler(w http.ResponseWriter, r *http.Request) {
	var documents CreateDocumentsPayload
	if err := readJSON(w, r, &documents); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	userId := "1"

	vector := &store.RagData{
		UserID: userId,
	}

	for _, doc := range documents.Documents {
		for _, subsection := range doc.Subsections {
			vector.Documents = append(vector.Documents, store.Document{
				Chapter: doc.Chapter,
				Subsections: []store.Subsection{
					{
						Title:   subsection.Title,
						Content: subsection.Content,
					},
				},
			})
		}
	}

	ctx := r.Context()

	if err := app.store.Vectors.CreateVectors(ctx, vector); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusCreated, vector); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

}
