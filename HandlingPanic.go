package main

import (
	"fmt"
	"time"
)

func main() {
	go goroutine1()
	time.Sleep(time.Hour)
	fmt.Println("program selesai")
}

func goroutine1() {
	var count int
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered")
			fmt.Println("start to re-run goroutines.....")
			count = 0
			goroutine1()
		}
	}()

	for {
		if count == 5 {
			panic("ada panik nih.....")
		}
		time.Sleep(5 * time.Second)
		fmt.Println("data from gouroutine1")
		count++
	}
}
