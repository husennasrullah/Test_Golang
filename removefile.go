package main

import (
	"log"
	"os"
)

func main () {
	err := os.Remove("./avatar.jpg")
	if err != nil {
		log.Fatal("error bro :", err)
	}
}
