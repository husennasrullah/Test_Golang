package main

import "fmt"

type LivingThings interface {
	Walk(length int64)
	Run(length int64, speed float64)
}

type Human struct {
	Age           int64
	Name          string
	HandsomeLevel int64
}

func (h Human) Walk(length int64) {
	fmt.Println(h.Name+" walk for ", length, " meters")
}
func (h Human) Run(length int64, speed float64) {
	fmt.Println(h.Name+" run for ", length, " meters with speed ", speed)
}

type Dog struct {
	Age  int64
	Name string
}

func (h Dog) Bark() {
	fmt.Println("Guk guk guk auuumm. Guk guk guk")
}
func (h Dog) Walk(length int64) {
	fmt.Println(h.Name+" walk for ", length, " meters")
}
func (h Dog) Run(length int64, speed float64) {
	fmt.Println(h.Name+" run for ", length, " meters with speed ", speed)
}

func AcceptLivingThingsOnly(an LivingThings, length int64, speed float64) {
	an.Run(length, speed)
	an.Walk(length)
}

func main() {
	imanTamvan := Human{
		Age:           22,
		Name:          "Iman Syahputra Situmorang",
		HandsomeLevel: 999,
	}
	// Call AcceptLivingThingsOnly
	AcceptLivingThingsOnly(imanTamvan, 20, 5)

	anjingImut := Dog{
		Age:  1,
		Name: "Simba",
	}
	// Call AcceptLivingThingsOnly
	AcceptLivingThingsOnly(anjingImut, 10, 6)
}
