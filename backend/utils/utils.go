package utils

import (
	"net/http"
)

type userKey string

const UserCtx userKey = "user"

func GetUserFromContext(r *http.Request) string {
	user, _ := r.Context().Value(UserCtx).(string)
	return user
}
