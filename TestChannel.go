package main

import (
	"fmt"
	"time"
)

func fast(num int, out chan<- int) {
	result := num * 2
	time.Sleep(15 * time.Millisecond)
	out <- result
}

func slow(num int, out chan<- int) {
	result := num * 2
	time.Sleep(15 * time.Millisecond)
	out <- result
}

func main() {
	out1 := make(chan int)
	out2 := make(chan int)

	// we start both fast and slow in different
	// goroutines with different channels
	go fast(8, out1)
	go slow(3, out2)

	// perform some action depending on which channel
	// receives information first
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case res := <-out1:
			fmt.Println("fast finished first, result:", res)
		case <-ticker.C:
			fmt.Println("slow finished first, result:", ticker.C)
		}
	}
}
