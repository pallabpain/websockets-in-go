package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var (
	addr = flag.String("addr", ":8080", "http address")
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{
		Scheme: "ws",
		Host:   *addr,
		Path:   "/echo",
	}
	log.Printf("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	done := make(chan struct{})

	// Start a go routine to read the message from server
	go func() {
		defer close(done)
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			log.Printf("received message: %s", string(msg))
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// Send a message every second
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := conn.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println(err)
				return
			}
		case <-interrupt:
			log.Println("interrupted")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseGoingAway, ""))
			if err != nil {
				log.Println(err)
				return
			}

			select {
			case <-done:
				log.Println("connection closed")
			case <-time.After(time.Second):
			}

			return
		}
	}
}
