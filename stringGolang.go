package main

import "fmt"

func main () {
		multiline := "line " +
			"by line \n" +
			"and line \n" +
			"after line"

		fmt.Print(multiline) // New lines as interpreted \n
}
