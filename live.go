package LanGenius

import (
	"golang.org/x/net/websocket"
	"net/http"
)

func chat(ws *websocket.Conn) {
	defer ws.Close()
}
func live(w http.ResponseWriter, r *http.Request) {

}
func camera(w http.ResponseWriter, r *http.Request) {

}
