package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mik-dmi/rag_chatbot/backend/utils"
)

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
			return
		}
		parts := strings.Split(authHeader, " ") // authorization:
		if len(parts) != 2 {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		token := parts[1]
		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)

		sessionIDFromJWT := fmt.Sprintf("%s", claims["sub"])

		sessionIDFromRequest := utils.GetUserFromContext(r)

		if sessionIDFromJWT != sessionIDFromRequest {
			app.logger.Infof("sessionIDFromJWT - %s ;sessionIDFromRequest - %s ", sessionIDFromJWT, sessionIDFromRequest)
			app.unauthorizedErrorResponse(w, r, ErrorSessionIDHeaderDifferentFromJWTSubject)
			return
		}

		next.ServeHTTP(w, r)
	})
}
