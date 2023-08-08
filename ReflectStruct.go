package main

import (
	"fmt"
	"reflect"
)

type tes struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func generateFromTypeOfStruct (x interface{}){
	rflct := reflect.TypeOf(x)
	fmt.Println(rflct)

	results :=  reflect.New(rflct)
	fmt.Println(results)
}

func main() {
	data := tes{
		Name:    "husen",
		Address: "Nasrullah",
	}

	generateFromTypeOfStruct(data)
}
