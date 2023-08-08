package main

import (
	"flag"
	"fmt"
)

func main () {
	var tes string
	flag.StringVar(&tes, "test", "default", "conversion type")

	flag.Parse()
	fmt.Println(tes)
}
