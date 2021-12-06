// This example shows how to work with websockets in Go, based on the example available at:
// https://gowebexamples.com/websockets/

package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"net/http"
)

const (
	ContentTypeHeaderKey       = "Content-Type"
	ContentTypeHeaderValueJson = "application/json; charset=UTF-8"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		wsConn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			handleError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if wsConn == nil {
			handleError(w, "no connection", http.StatusInternalServerError)
			return
		}
		for {
			// Read message from websocket client (e.g., browser)
			msgType, msg, err := wsConn.ReadMessage()
			if err != nil {
				return
			}

			// Log the message locally
			glog.Infof("%s sent: %s\n", wsConn.RemoteAddr(), string(msg))

			// Write message back to websocket client (e.g., browser)
			if err = wsConn.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websockets.html")
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		glog.Errorf("%s", err)
	}
	glog.Flush()
}

func handleError(w http.ResponseWriter, errMsg string, statusCode int) {
	response := struct {
		Err string
	}{Err: fmt.Sprintf("%s\n", errMsg)}
	responseJSON, _ := json.Marshal(response)
	w.Header().Set(ContentTypeHeaderKey, ContentTypeHeaderValueJson)
	w.WriteHeader(statusCode)
	_, _ = w.Write(responseJSON)
}
