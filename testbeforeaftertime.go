package main

import (
	"fmt"
	"time"
)

func main() {
	var output time.Time
	validFrom := "2023-05-20"

	output, _ = time.Parse("2006-01-02", validFrom)
	currentDate := time.Now().Truncate(time.Hour * 24)
	fromDate := (output).Truncate(time.Hour * 24)

	fmt.Println(currentDate.Add(-time.Hour * 24 * 31))

	if fromDate.Before(currentDate.Add(-31)) {
		fmt.Println("salah")
	} else {
		fmt.Println("benar")
	}
}
