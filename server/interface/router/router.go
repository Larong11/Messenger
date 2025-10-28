package router

import (
	"net/http"
	"server/interface/http/handlers"
	"server/interface/websocket"
)

func NewRouter(uh *handlers.UserHandler, ws *websocket.WSHandler) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/api/users/check-username", uh.CheckUserName)
	router.HandleFunc("/api/users/check-email", uh.CheckEmail)
	router.HandleFunc("/api/users/create-user", uh.CreateUser)
	router.HandleFunc("/api/users/request-verification-code", uh.RequestVerificationCode)
	router.HandleFunc("/ws", ws.SendMessage)
	return router
}
