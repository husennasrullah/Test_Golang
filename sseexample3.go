package main

import (
	"fmt"
	"github.com/bernerdschaefer/eventsource"
	"net/http"
	"strconv"
	"time"
)

func main() {
	es := eventsource.Handler(func(lastID string, e *eventsource.Encoder, stop <-chan bool) {
		var id int64
		for {
			select {
			case <-time.After(3 * time.Second):
				fmt.Println("sending event...")
				id += 1
				e.Encode(eventsource.Event{ID: strconv.FormatInt(id, 10),
					Type: "add",
					Data: []byte("some data")})
			case <-stop:
				return
			}
		}
	})
	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		es.ServeHTTP(w, r)
	})
	if e := http.ListenAndServe(":9090", nil); e != nil {
		fmt.Println(e)
	}
	//}
}