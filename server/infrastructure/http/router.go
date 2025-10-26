package http

import (
	"net/http"
	"server/infrastructure/http/handlers"
)

func NewRouter(uh *handlers.UserHandler) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/api/users/check-username", uh.CheckUserName)
	router.HandleFunc("/api/users/check-email", uh.CheckEmail)
	return router
}
