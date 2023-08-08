package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"time"
)

//data dari DB
type loginLog struct {
	Waktu  time.Time `json:"waktu"`
	Id     int       `json:"id"`
	Name   string    `json:"name"`
	Status string    `json:"status"`
	Alamat string    `json:"alamat"`
}

func DataExport(data []loginLog) (result [][]string) {
	//1. create header
	head := []string{
		"waktu", "id", "name", "status", "alamat",
	}

	result = append(result, head)

	//2. create content
	for i := 0; i < len(data); i++ {
		var tempContent []string
		tes := reflect.ValueOf(data[i])

		for i := 0; i < tes.NumField(); i++ {
			var temp string
			item := tes.Field(i).Interface()

			switch item.(type) {
			case float64, float32:
				temp = fmt.Sprintf("%0.2f", tes.Field(i))
			case time.Time:
				times := item.(time.Time)
				temp = times.Format("02 Sep 15 08:00 WIB")
			default:
				temp = fmt.Sprintf("%v", tes.Field(i))

			}
			tempContent = append(tempContent, temp)
		}

		result = append(result, tempContent)

	}

	return

}

func main() {
	var kirimData []loginLog
	data := loginLog{
		Waktu:  time.Now(),
		Id:     1,
		Name:   "husen",
		Status: "jdasjdas",
		Alamat: "tangerang",
	}

	kirimData = append(kirimData, data)

	exportData := DataExport(kirimData)


	berkas, errs := os.Create("tes.csv")
	if errs != nil {
		os.Exit(1)
	}

	csvWriter := csv.NewWriter(berkas)

	errs = csvWriter.WriteAll(exportData)
	if errs != nil {
		os.Exit(1)
	}
}
