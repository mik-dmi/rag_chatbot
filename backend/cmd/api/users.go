package main

import (
	"net/http"

	"github.com/mik-dmi/rag_chatbot/backend/internal/store"
)

type RegiterUserPayload struct {
	Username string `json:"username" validate:"required,max=50,min=3"`
	Email    string `json:"email" validate:"required,max=50,email"`
	Password string `json:"password" validate:"required,max=50,min=3"`
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {

	userId := r.PathValue("userId")
	//validate the user input
	if err := Validate.Struct(userId); err != nil {
		app.badRequestError(w, r, err)
	}

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

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {

	var newUserData RegiterUserPayload
	if err := readJSON(w, r, &newUserData); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(newUserData); err != nil {
		app.badRequestError(w, r, err)
	}

	ctx := r.Context()

	app.postgreStore.CreateUser(ctx)

}
