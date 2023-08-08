package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	tgl := "2023-01-31"

	for i := 1; i <= 3; i++ {
		fromDate, thruDate := getThreeLastMonth(tgl, i)

		fmt.Println(fromDate, " --> ", thruDate)
	}

}

func StrToDateFormatNCO(input string) (temp time.Time, errorS error) {
	if len(strings.Split(input, "/")) > 2 {
		temp, errorS = time.Parse("2/1/2006", input)
	} else {
		temp, errorS = time.Parse("2006-1-2", input)
	}
	return
}

func getThreeLastMonth(strDate string, index int) (from, last string) {
	const layout = "2006-01-02"
	t, _ := time.Parse(layout, strDate)

	year, month, _ := t.Date()

	newMonth := month - time.Month(index)
	fmt.Println(newMonth)
	startTime := time.Date(year, newMonth, 1, 0, 0, 0, 0, t.Location())

	endTime := startTime.AddDate(0, 1, 0).Add(-time.Nanosecond)

	return startTime.Format(layout), endTime.Format(layout)
}

func getLastMonth(strDate string) (from, last string) {
	const layout = "2006-01-02"
	t, _ := time.Parse(layout, strDate)

	year, month, _ := t.Date()

	startTime := time.Date(year, month-1, 1, 0, 0, 0, 0, t.Location())
	endTime := startTime.AddDate(0, 1, 0).Add(-time.Nanosecond)

	return startTime.Format(layout), endTime.Format(layout)
}
