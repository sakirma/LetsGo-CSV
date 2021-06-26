package UnitTesting

import (
	reader "Reader"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"testing"
)

type readerMock struct {
	mock.Mock
}

func (r *readerMock) Read() (record []string, err error) {
	fmt.Println("Mocked Read function of Reader")

	args := r.Called()
	return args.Get(0).([]string), args.Error(1)
}

func Test_read_init(t *testing.T) {
	// Assign
	csvReader := new(readerMock)

	csvReader.On("Read").Return([]string{"0", "1", "69", "1624625795"}, nil).Once()
	csvReader.On("Read").Return([]string{"1", "1", "100", "1624625795"}, nil)

	// Act
	err := reader.Init(csvReader)
	if err != nil {
		return
	}

	// Assert
	assert.Equal(t, reader.GetBus()[0].Reading, float32(69))
	assert.Equal(t, reader.GetBus()[1].Reading, float32(100))
}

func Test_read_init_ignores_invalid_usage(t *testing.T) {
	// Assign
	csvReader := new(readerMock)

	csvReader.On("Read").Return([]string{"0", "1", "69", "1624625795"}, nil).Once()
	csvReader.On("Read").Return([]string{"0", "1", "50", "1624625795"}, nil).Once()
	csvReader.On("Read").Return([]string{"1", "1", "100", "1624625795"}, nil)

	// Act
	err := reader.Init(csvReader)
	if err != nil {
		return
	}

	// Assert
	assert.Equal(t, reader.GetBus()[0].Reading, float32(50))
	assert.Equal(t, reader.GetBus()[1].Reading, float32(100))
}

func Test_read_get_next_cost(t *testing.T) {
	// Assign
	result0 := reader.Cost{Id: 0, Value: 0.0062000006}
	result1 := reader.Cost{Id: 1, Value: 0.0039999993}
	result2 := reader.Cost{Id: 2, Value: 0.006000002}

	csvReader := new(readerMock)

	csvReader.On("Read").Return([]string{"0", "1", "69", "1624626695"}, nil).Once()
	csvReader.On("Read").Return([]string{"1", "1", "100", "1624627595"}, nil).Once()
	csvReader.On("Read").Return([]string{"2", "1", "120", "1624628495"}, nil).Once()
	csvReader.On("Read").Return([]string{"3", "1", "150", "1624629395"}, nil).Once()
	csvReader.On("Read").Return([]string{"3", "1", "250", "1624630295"}, nil).Once()
	csvReader.On("Read").Return([]string{"3", "1", "300", "1624631195"}, io.EOF).Once()

	// Act
	err := reader.Init(csvReader)
	if err != nil {
		return
	}

	var costs []*reader.Cost
	for  {
		cost, err := reader.GetNextCost(csvReader)
		if err != nil {
			break
		}

		costs = append(costs, cost)
	}


	// Assert
	assert.True(t, len(costs) == 3)
	assert.Equal(t, &result0, costs[0])
	assert.Equal(t, &result1, costs[1])
	assert.Equal(t, &result2, costs[2])
}

func Test_usage_can_check_if_valid(t *testing.T) {
	// Assign
	wrongUsage := float32(-100)
	wrongUsage1 := float32(110)
	wrongUsage2 := float32(-0.0001)

	validUsage := float32(69)

	// Assert
	assert.False(t, reader.IsUsageValid(wrongUsage))
	assert.False(t, reader.IsUsageValid(wrongUsage1))
	assert.False(t, reader.IsUsageValid(wrongUsage2))

	assert.True(t, reader.IsUsageValid(validUsage))
}