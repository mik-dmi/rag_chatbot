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

	vectorObject, err := app.store.Vectors. (ctx, query.UserMessage)





	
}
func (app *application) deleteVectorObjectByIdHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print(r.PathValue("id"))

}