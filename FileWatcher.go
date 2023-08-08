package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok:= <-watcher.Events:
				if !ok {
					fmt.Println("file done")
					return
				}
				log.Println("event:", event)
				fmt.Println(event.Name)
				if event.Op&fsnotify.Write == fsnotify.Write {
					//log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}

		}
	}()

	err = watcher.Add("C:\\Test_Golang\\tmp\\foo")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}


