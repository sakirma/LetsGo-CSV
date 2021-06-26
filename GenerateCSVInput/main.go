package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var columns = []string{"metering_point_id", "type", "reading", "created_at"}

const averageGasM3 float32 = 0.048
const averageElectricityWh float32 = 81.621

func main() {
	if len(os.Args) >= 2 {
		i, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("First argument is not a number")
			os.Exit(1)
		}

		t := 1
		if len(os.Args) == 3 {
			t, err = strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Printf("second argument is not a number")
				os.Exit(1)
			}
		}

		fmt.Println("Amount is : " + os.Args[1] + "\n")
		GenerateCSVFile(i, t)
		os.Exit(0)
	}

	GenerateCSVFile(10, 1)
}

func GenerateCSVFile(amount int, t int) {
	file, err := os.Create("input.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush() // I love defer

	// First field (Columns)
	writeInFile(writer, columns)

	// Could use parameter to use a seed
	currentTime := time.Now()

	fmt.Println(time.Now())

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var currentReading float32

	for i := 0; i < amount; i++ {

		reading := currentReading
		// Electricity
		if t == 1 {
			reading += averageElectricityWh + float32(r.Intn(10)-5)
		} else { // Gas
			reading += averageGasM3 + r.Float32()*0.01
		}

		currentReading = reading

		unixTimeStamp := currentTime.Add(time.Minute * time.Duration(i*15)).Unix()
		row := []string{strconv.Itoa(i), strconv.Itoa(t), fmt.Sprintf("%f", reading), strconv.FormatInt(unixTimeStamp, 10)}
		writeInFile(writer, row)
	}

	elapsed := time.Since(currentTime)
	log.Printf("It took: %s", elapsed)
}

func writeInFile(w *csv.Writer, record []string) {
	err := w.Write(record)
	checkError("Cannot write to file", err)
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
