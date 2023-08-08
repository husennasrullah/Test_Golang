package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

//variabel channel buat nampung
var dataSend = make(chan string, 2)

func main() {
	//receive handler socket
	go sendGoroutines()
	go sendGoroutines2()

	handler := mux.NewRouter()
	handler.HandleFunc("/testapi", handleApi())
	log.Fatal("HTTP server error: ", http.ListenAndServe("localhost:3000", handler))
}

func handleApi() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//todo send to socket (skipped)
		// waiting for channel datasend
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		for {
			select {
			case data := <-dataSend:
				var output string
				b, _ := json.Marshal(data)
				output = string(b)
				w.WriteHeader(200)
				_, _ = w.Write([]byte(output))
				return
			case <-ctx.Done():
				output := "NO DATA FOUND"
				w.WriteHeader(200)
				_, _ = w.Write([]byte(output))
				return
			}

		}

	}
}

func sendGoroutines() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for ; ; <-ticker.C {
		select {
		default:
			test := dataName{
				Nama:   "husen",
				Alamat: "tangerang",
			}
			kirim, _ := json.Marshal(test)

			dataSend <- string(kirim)
		}
	}
}

func sendGoroutines2() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for ; ; <-ticker.C {
		select {
		default:
			test := dataName{
				Nama:   "lilik",
				Alamat: "tangerang",
			}
			kirim, _ := json.Marshal(test)

			dataSend <- string(kirim)
		}
	}
}

type dataName struct {
	Nama   string
	Alamat string
}
