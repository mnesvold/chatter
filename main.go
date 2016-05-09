package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
)

const (
	staticRoot = "./src/github.com/mnesvold/chatter/www"
)

var (
	port = flag.Int("port", 8000, "port to serve site over")
)

var (
	connections   []*websocket.Conn
	broadcastChan chan []byte
	closeChan     chan bool
)

func broadcast() {
	for {
		select {
		case payload := <-broadcastChan:
			for _, conn := range connections {
				conn.Write(payload)
			}
		case <-closeChan:
			break
		}
	}
}

func echoServer(ws *websocket.Conn) {
	connections = append(connections, ws)
	decoder := json.NewDecoder(ws)
	for {
		var payload map[string]interface{}
		if err := decoder.Decode(&payload); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		message := payload["message"]

		response := make(map[string]interface{})
		response["message"] = message
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Fatal(err)
		}
		broadcastChan <- responsePayload
	}
}

func main() {
	flag.Parse()

	broadcastChan = make(chan []byte)
	closeChan = make(chan bool)
	go broadcast()

	//http.HandleFunc("/mirror/", serveMirror)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(staticRoot))))
	http.Handle("/chat", websocket.Handler(echoServer))

	log.Printf("Listening on port %d\n", *port)
	bind := fmt.Sprintf(":%d", *port)
	log.Fatal(http.ListenAndServe(bind, nil))
}
