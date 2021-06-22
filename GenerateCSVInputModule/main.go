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

type InputReadings struct {
	meteringPointId int
	readingType     string
	reading         int
	createdAt       time.Time
}

var columns = []string{"metering_point_id", "type", "reading", "created_at"}

const averageGasM3 float32 = 0.052
const averageElectricityKWh float32 = 81.621

func main() {
	if len(os.Args) > 2 {
		fmt.Printf("Please give amount as argument")
		os.Exit(1)
	} else if len(os.Args) == 2 {
		i, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("Given argument is not a number")
			os.Exit(1)
		}

		fmt.Println("Amount is : " + os.Args[1] + "\n")
		GenerateCSVFile(i)
		os.Exit(0)
	}

	GenerateCSVFile(70000000)
}

func GenerateCSVFile(amount int) {
	file, err := os.Create("result.csv")
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

	for i := 0; i < amount; i++ {

		readingTypeNumber := r.Intn(2) + 1
		var reading float32
		// Electricity
		if readingTypeNumber == 1 {
			reading = averageElectricityKWh + float32(r.Intn(10)-5)
		} else { // Gas
			reading = averageGasM3 + r.Float32()*0.01
		}

		readingType := strconv.Itoa(readingTypeNumber)

		unixTimeStamp := currentTime.Add(time.Minute * time.Duration(i*15)).Unix()
		row := []string{strconv.Itoa(i), readingType, fmt.Sprintf("%f", reading), strconv.FormatInt(unixTimeStamp, 10)}
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
