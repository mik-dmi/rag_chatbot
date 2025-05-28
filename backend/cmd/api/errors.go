package main

import (
	"errors"
	"net/http"
)

var (
	ErrorMissingSessionIDHeader                 = errors.New("error missing session ID in the Header of the request")
	ErrorSessionIDHeaderDifferentFromJWTSubject = errors.New("error session ID in the Header is different from JWT Subject")
	ErrorUserNotAuthorized                      = errors.New("user not authorized")
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorf("internal server error: %s", err)
	writeJSONError(w, http.StatusInternalServerError, "server encountered a problem")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorf("bad request error: %s", err)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorf("not found error: %s", err)
	writeJSONError(w, http.StatusNotFound, err.Error())
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauthorized error - method=%s, path=%s, error=%s", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusBadRequest, err.Error())
}
