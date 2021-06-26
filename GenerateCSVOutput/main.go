package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

import csvReader "Reader"

var columns = []string{"metering_point_id", "total_price"}

func main() {
	file, err := os.Create("result.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush() // I love defer

	err = writer.Write(columns)
	checkError("Couldn't write in new file.", err)


	ReadAndWriteInput(writer)
}

func ReadAndWriteInput(writer *csv.Writer) {
	file := "./input.csv"
	f, err := os.Open(file)
	if err != nil {
		log.Panic(err)
	}

	defer f.Close()

	reader := csv.NewReader(f)
	_, err = reader.Read()
	if err != nil {
		log.Panic(err)
	}

	err = csvReader.Init(reader)
	if err != nil {
		log.Panic(err)
	}

	for {
		cost, err := csvReader.GetNextCost(reader)

		if cost != nil {
			writeErr := writer.Write([]string{
				strconv.FormatUint(cost.Id, 10),
				fmt.Sprintf("%f", cost.Value),
			})
			checkError("Couldn't write to new file.", writeErr)
		}


		if err != nil {
			break
		}
	}
}

func checkError(errorMessage string, err error) {
	if err != nil {
		log.Fatal(errorMessage, err)
	}
}