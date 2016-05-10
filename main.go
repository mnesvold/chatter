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

type chatServer struct {
	connections   []*websocket.Conn
	broadcastChan chan []byte
}

func NewServer() *chatServer {
	server = &chatServer{broadcastChan: make(chan []byte)}
	go server.broadcast()
	return server
}

func (s *chatServer) broadcast() {
	for {
		payload := <-s.broadcastChan
		for _, conn := range s.connections {
			conn.Write(payload)
		}
	}
}

func (s *chatServer) HandleConnection(ws *websocket.Conn) {
	s.connections = append(s.connections, ws)
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
		s.broadcastChan <- responsePayload
	}
}

var (
	port = flag.Int("port", 8000, "port to serve site over")
)

func main() {
	flag.Parse()

	server := NewServer()

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(staticRoot))))
	http.Handle("/chat", websocket.Handler(server.HandleConnection))

	log.Printf("Listening on port %d\n", *port)
	bind := fmt.Sprintf(":%d", *port)
	log.Fatal(http.ListenAndServe(bind, nil))
}
