package main

import (
	"net/http"

	"github.com/mik-dmi/rag_chatbot/backend/internal/store"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {

	userId := r.PathValue("userId")

	ctx := r.Context()
	user, err := app.postgreStore.Users.GetUserById(ctx, userId)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.badRequestError(w, r, err)
			return

		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)

	}
}
