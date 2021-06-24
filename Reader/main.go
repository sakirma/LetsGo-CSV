package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"
)

type Cost struct {
	Id    uint64
	Value float64
}

type row struct {
	Id      uint64
	Reading float64
	Type    int
	Date    time.Time
}

var bus [2]row // TODO: Assumption Explain why the size is two in readme

func main() {
	file := "./result.csv"
	f, err := os.Open(file)
	if err != nil {
		log.Panic("Couldn't open file " + file)
	}

	defer f.Close()

	reader := csv.NewReader(f)
	reader.Read() // Ignore headers

	rows, err := Init(reader)
	if err != nil {
		log.Panic("Could not iterate through the file")
	}

	for _, _ = range rows {

	}
}

func Init(reader *csv.Reader) ([2]row, error) {
	// Find the first valid value
	for {
		read, err := reader.Read()
		if err != nil {
			return [2]row{}, err
		}

		parsedRow := parseRow(read)
		if isRowValid(parsedRow) {
			bus[0] = parsedRow
			break
		}
	}

	// Get the next value. Replace the non-valid value with the previous value (check assumption).
	read, err := reader.Read()
	if err != nil {
		return [2]row{}, err
	}

	parsedRow := parseRow(read)
	if isRowValid(parsedRow) {
		bus[1] = parsedRow
	} else {
		bus[1] = parsedRow
	}

	return bus, nil
}

// GetNextUsage returns the next usage of the readings
func GetNextUsage(reader *csv.Reader) (Cost, error) {
	cost := readingToCost(bus[0], bus[1])

	r, err := reader.Read()
	if err != nil {
		return cost, err
	}

	nextRow := parseRow(r)
	if !isRowValid(nextRow) {
		diff := bus[1].Reading - bus[0].Reading
		nextRow.Reading = clampReading(bus[1].Reading+diff, 0, 100) // TODO: Assumption: Linear Readings should be kept min maxed or else we will get problems with a continues invalid data.
	}
	addToBus(parseRow(r))

	return cost, nil
}

func clampReading(reading float64, min float64, max float64) float64 {
	if reading > max {
		return max
	}
	if reading < min {
		return min
	}

	return reading
}

const idIndex = 0
const typeIndex = 1
const readingIndex = 2
const timeIndex = 3

func parseRow(r []string) row {
	id, err := strconv.ParseUint(r[idIndex], 10, 64)
	if err != nil {
		panicNumberConversion(r[idIndex])
	}

	reading, err := strconv.ParseFloat(r[readingIndex], 64)
	if err != nil {
		panicNumberConversion(r[readingIndex])
	}

	typeNumber, err := strconv.Atoi(r[typeIndex])
	if err != nil {
		panicNumberConversion(r[typeIndex])
	}

	i, err := strconv.ParseInt(r[timeIndex], 10, 64)
	if err != nil {
		panicNumberConversion(r[timeIndex])
	}

	return row{id, reading, typeNumber, time.Unix(i, 0)}
}

func readingToCost(row1 row, row2 row) Cost {
	usage := getReadingKWh(row2) - getReadingKWh(row1)

	var cost float64
	if row1.Type == 2 {
		cost = usage / 0.06
	} else {
		day := row1.Date.Weekday()
		if day == time.Saturday || day == time.Sunday {
			cost = usage / 0.18
		} else {
			hours := row1.Date.Hour()
			if hours >= 7 && hours <= 23 {
				cost = usage / 0.2
			} else {
				cost = usage / 0.18
			}
		}
	}

	return Cost{ Id: row1.Id, Value: cost}
}

func getReadingKWh(r row) float64 {
	if r.Type == 1 {
		return r.Reading / 1000
	} else {
		return r.Reading * 9.769
	}
}

func isRowValid(r row) bool {
	return r.Reading >= 0 && r.Reading <= 100
}

func addToBus(r row) {
	temp := bus[1]
	bus[0] = temp
	bus[1] = r
}

func panicNumberConversion(name string) {
	log.Panic("Could not convert " + name + " reading into a number")
}
