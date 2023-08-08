package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func main () {
	filePath := "C:\\Test_Golang\\tmp\\foo\\obd.csv"
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file " + filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for " + filePath, err)
	}

	fmt.Println(records)
	for i := 0; i < len(records); i++ {
		fmt.Println("jumlah data ke :", i, " = ", len(records[i]))
	}

}
