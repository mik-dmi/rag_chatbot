package main

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtTokenPayload struct {
	ClientID string `json:"user" validate:"required,max=50"`
	Password string `json:"password" validate:"required,max=70,min=3 "`
}

func (app *application) jwtTokenHandler(w http.ResponseWriter, r *http.Request) {
	var credentials jwtTokenPayload

	if err := readJSON(w, r, &credentials); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	//validate the user input
	if err := Validate.Struct(credentials); err != nil {
		app.badRequestError(w, r, err)
	}
	clientCredencials := app.config.authCredencials

	if clientCredencials.authCredencials.clientID != credentials.ClientID {
		app.unauthorizedErrorResponse(w, r, ErrorUserNotAuthorized)
		return
	}
	if clientCredencials.authCredencials.password != credentials.Password {
		app.unauthorizedErrorResponse(w, r, ErrorUserNotAuthorized)
		return
	}

	//get session_iD, that represents a unique user
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		app.badRequestResponse(w, r, ErrorMissingSessionIDHeader)
		return
	}

	claims := jwt.MapClaims{
		"suv": userID,
		"exp": time.Now().Add(app.config.authCredencials.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.authCredencials.token.iss,
		"aud": app.config.authCredencials.token.iss,
	}
	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, token); err != nil {
		app.internalServerError(w, r, err)

	}

}
