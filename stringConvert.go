package main

import (
	"fmt"
	"time"
)

func main() {
	tes := time.Now().Format("06") + "01"

	fmt.Println(tes)
}