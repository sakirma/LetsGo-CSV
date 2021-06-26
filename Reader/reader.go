package reader

import (
	"log"
	"strconv"
	"time"
)

type Cost struct {
	Id    uint64
	Value float32
}

type Row struct {
	Id      uint64
	Reading float32
	Type    int
	Date    time.Time
}

var bus [2]Row // TODO: Assumption Explain why the size is two in readme

type Reader interface {
	Read() (record []string, err error)
}

// Init Returns error if something went bad while reading the file
// Returns the array of the first two costs
func Init(reader Reader) error {

	getNextRow := func() (Row, error) {
		read, err := reader.Read()
		if err != nil {
			return Row{}, err
		}

		return parseRow(read), nil
	}

	currentReading, err := getNextRow()
	if err != nil {
		return err
	}
	for {
		row2, err := getNextRow()
		if err != nil {
			return err
		}

		usage := getUsage(currentReading, row2)
		if IsUsageValid(usage) {
			bus[0] = currentReading
			bus[1] = row2
			break
		} else {
			currentReading = row2
		}
	}

	return nil
}

// GetNextCost returns the next usage of the readings
func GetNextCost(reader Reader) (*Cost, error) {
	cost := ReadingToCost(bus[0], bus[1])

	r, err := reader.Read()
	if err != nil {
		return &cost, err
	}

	nextRow := parseRow(r)
	if !IsUsageValid(getUsage(bus[1], nextRow)) {
		diff := bus[1].Reading - bus[0].Reading
		nextRow.Reading = clampReading(bus[1].Reading+diff, 0, 100)
	}
	addToBus(nextRow)

	return &cost, nil
}

func clampReading(reading float32, min float32, max float32) float32 {
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

func parseRow(r []string) Row {
	id, err := strconv.ParseUint(r[idIndex], 10, 64)
	if err != nil {
		panicNumberConversion(r[idIndex])
	}

	reading, err := strconv.ParseFloat(r[readingIndex], 32)
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

	return Row{id, float32(reading), typeNumber, time.Unix(i, 0)}
}

func getUsage(row1 Row, row2 Row) float32 {
	return getReadingKWh(row2) - getReadingKWh(row1)
}

func ReadingToCost(row1 Row, row2 Row) Cost {
	usage := getUsage(row1, row2)

	var cost float32
	if row1.Type == 2 {
		cost = usage * 0.06
	} else {
		day := row1.Date.Weekday()
		if day == time.Saturday || day == time.Sunday {
			cost = usage * 0.18
		} else {
			hours := row1.Date.Hour()
			if hours >= 7 && hours <= 23 {
				cost = usage * 0.2
			} else {
				cost = usage * 0.18
			}
		}
	}

	return Cost{Id: row1.Id, Value: cost}
}

func getReadingKWh(r Row) float32 {
	if r.Type == 1 {
		return r.Reading / 1000
	} else {
		return r.Reading * 9.769
	}
}

func IsUsageValid(usage float32) bool {
	return usage >= 0 && usage <= 100
}

func addToBus(r Row) {
	temp := bus[1]
	bus[0] = temp
	bus[1] = r
}

func panicNumberConversion(name string) {
	log.Panic("Could not convert " + name + " reading into a number")
}

func GetBus() [2]Row {
	return bus
}