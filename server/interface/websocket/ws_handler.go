package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WSHandler struct {
}

func NewWSHandler() *WSHandler {
	return &WSHandler{}
}
func (uh *WSHandler) SendMessage(writer http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("📩 Message: %s", msg)
		err = conn.WriteMessage(websocket.TextMessage, []byte("Hello from server"))
		if err != nil {
			return
		}
	}
}
