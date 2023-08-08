package main

import (
	"fmt"
	"time"
)

const retryperiode = time.Second * 5

func main() {
	var count, index int
	intervalTime := []time.Duration{
		time.Second, time.Second * 3, time.Second * 5, time.Second * 8, time.Second * 10,
	}

	ticker := time.NewTicker(intervalTime[index])

	for ; ; <-ticker.C {
		fmt.Println(time.Now().Format("15:04:05"), " ...halo")

		if count == 5 {
			fmt.Println("=========================")
			if index > len(intervalTime)-1 {
				index = 0
			}
			ticker.Reset(intervalTime[index])
			count = 0
			index++
		}

		count++
	}

}
