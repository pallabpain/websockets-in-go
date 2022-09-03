package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
	addr     = flag.String("addr", ":8080", "http address")
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	http.HandleFunc("/echo", echoServiceHandler)

	log.Println("Started listening on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		panic(err)
	}
}

func echoServiceHandler(w http.ResponseWriter, r *http.Request) {
	// Accept all origin headers
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	// Upgrade the HTTP connection to a Websocket
	conn, _ := upgrader.Upgrade(w, r, nil)

	defer conn.Close()

	// Since this is an echo service, keep reading messages
	// and writing them back
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("client %s sent: %s\n", conn.RemoteAddr(), string(msg))

		if err = conn.WriteMessage(msgType, msg); err != nil {
			log.Println(err)
		}
	}
}
