package main

import (
	"fmt"
	"strconv"
)

func main () {
	test, errS := strconv.ParseInt("", 10, 32)
	fmt.Println(test)
	if errS != nil {
		fmt.Println(errS)
	}
}
