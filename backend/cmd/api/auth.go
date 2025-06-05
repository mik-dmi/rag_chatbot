package main

import (
	"net/http"

	"github.com/mik-dmi/rag_chatbot/backend/internal/store"
)

func(app *application )( w http.ResponseWriter, r *http.Request){
	var payload RegiterUserPayload
	if err := readJSON(w,r, payload ); err != nil{
		app.badRequestResponse(w,r,err)
		return 
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w,r,err)
		return 
	}

	user := &store.PostgreUser{
		Username: payload.Username,
		Email: payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w,r , err)
		return 
	}

	ctx := r.Context()

	err := app.postgreStore.Users.CreateAndInvite(ctx, user, "token-123")
	if err := app.jsonResponse(w,http.StatusCreated, nil); err != nil{
		app.internalServerError(w,r,err)
		
	}

}