package main

import (
	"fmt"
	"log"
	"time"
)

func main () {
	layout := "01/02/2006"
	date :=  "07/26/2022"

	dt, err := time.Parse(layout, date)
	if err != nil {
		log.Fatal("gagaglllll :", err)
	} else {
		fmt.Println(dt)
	}


}
