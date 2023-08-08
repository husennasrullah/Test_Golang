package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var done chan interface{}
var interrupt chan os.Signal

func receiveHandler(connection *websocket.Conn) {
	defer close(done)
	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			log.Println(`Error in receive:`, err)
			return
		}
		log.Printf(`Received: %s\n`, msg)
	}
}

func main() {
	done = make(chan interface{})    // Channel to indicate that the receiverHandler is done
	interrupt = make(chan os.Signal) // Channel to listen for interrupt signal to terminate gracefully

	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	socketUrl := `ws://localhost:8091/nexsocket-notif/235/q3mdm3yi/websocket`
	d := websocket.DefaultDialer
	header := http.Header{}
	header.Set("X-Token-Nexsoft", "bmV4Y2hpZWYuZGV2LWM1NzkwMTcwYjVmNzQ2ZGViNGYwYmZjMTUyOGU0MjQ0LTE2NDk5MjY2NjkwOTM")
	header.Set("Cookie", "SOCKETSERVER=")
	header.Set("heart-beat", "25000")
	header.Set("Origin", "https://socket.gromart.club:8443")

	conn, _, err := d.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal(`Error connecting to Websocket Server:`, err)
	}
	defer conn.Close()
	go receiveHandler(conn)

	// Our main loop for the client
	// We send our relevant packets here
	for {
		select {
		case <-time.After(time.Duration(1) * time.Millisecond * 1000):
			// Send an echo packet every second
			err = conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second*2))
			if err != nil {
				//conn.closeWs()
				return
			}

		case <-interrupt:
			// We received a SIGINT (Ctrl + C). Terminate gracefully…
			log.Println(`Received SIGINT interrupt signal. Closing all pending connections`)

			// Close our websocket connection
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println(`Error during closing websocket:`, err)
				return
			}

			select {
			case <-done:
				log.Println(`Receiver Channel Closed! Exiting….`)
			case <-time.After(time.Duration(1) * time.Second):
				log.Println(`Timeout in closing receiving channel. Exiting….`)
			}
			return
		}
	}
}
