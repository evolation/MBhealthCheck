//go:build !js || !wasm

package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleModbusCommunication(conn *websocket.Conn) {
	runModbusHealthcheckTarget("p_filename string", "COM14")
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error during WebSocket upgrade:", err)
		return
	}
	defer conn.Close()

	// Handle WebSocket communication (to be implemented).
	handleModbusCommunication(conn)
}

func main() {

	// WebSocket endpoint
	http.HandleFunc("/ws", wsHandler)
	if err := http.ListenAndServe(":7998", nil); err != nil {
		log.Fatal(err)
	}

}
