package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var messageChan chan string

var mapStatus = make(map[string]checkSyncStatus)

type content struct {
	data interface{} `json:"data"`
}

type checkSyncStatus struct {
	Counter int    `json:"counter"`
	Total   int    `json:"total"`
	Status  string `json:"status"`
}

type hasil struct {
	Status string `json:"status"`
	Catatan string `json:"catatan"`
}

func handleSSE() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Printf("Get handshake from client")

		// prepare the header
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// instantiate the channel
		messageChan = make(chan string)

		// close the channel after exit the function
		defer func() {
			close(messageChan)
			messageChan = nil
			log.Printf("client connection is closed")
		}()

		// prepare the flusher
		flusher, _ := w.(http.Flusher)

		// trap the request under loop forever
		for {

			select {

			// message will received here and printed
			case message := <-messageChan:
				fmt.Fprintf(w, "%s\n", message)
				flusher.Flush()
				if message == "done" {
					return
				}

			// connection is closed then defer will be executed
			case <-r.Context().Done():
				return

			}
		}

	}
}

func sendMessage(message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if messageChan != nil {
			log.Printf("print message to client")

			// send the message through the available channel
			messageChan <- message
		}

	}
}

func updateStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		var count int
		var total = 100
		UUID, _ := uuid.NewRandom()
		output := UUID.String()

		defer func() {
			delete(mapStatus, output)
		}()

		for {
			count++

			a := checkSyncStatus{
				Counter: count,
				Total:   total,
				Status:  "ON-PROGRESS",
			}

			if count <= total {
				if count == total {
					a.Status = "DONE"
				}

				time.Sleep(time.Millisecond *100)
				mapStatus[output] = a
				// bagian untuk write event respon secara buffered
				tes, _ := json.Marshal(mapStatus[output])
				flusher, _ := w.(http.Flusher)
				fmt.Fprintf(w, "data: %v\n\n", string(tes))
				flusher.Flush()
			} else {
				break
			}

		}
		//testt := hasil{
		//	Status:  "Berhasilllllllll",
		//	Catatan: "testttingggggggg",
		//}
		//b, _ := json.Marshal(testt)

		file, errs := os.Open("aa.xls")
		if errs != nil {
			return
		}

		_, _ = io.Copy(w, file)

		w.WriteHeader(200)
		//_, _ = w.Write(file)
	}
}

func cekStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var output string
		uuids := mux.Vars(r)["ID"]
		respon := mapStatus[uuids]

		if respon.Status == "" {
			output = "Data Not Found"
		} else {
			b, _ := json.Marshal(respon)
			output = string(b)
		}

		w.WriteHeader(200)
		_, _ = w.Write([]byte(output))
	}
}

func main() {

	handler := mux.NewRouter()

	handler.HandleFunc("/counter", updateStatus())

	handler.HandleFunc("/cekstatus/{ID}", cekStatus())

	handler.HandleFunc("/handshake", handleSSE())

	handler.HandleFunc("/sendmessage", sendMessage("hello client"))

	handler.HandleFunc("/sendmessage2", sendMessage("done"))

	log.Fatal("HTTP server error: ", http.ListenAndServe("localhost:3000", handler))
}
