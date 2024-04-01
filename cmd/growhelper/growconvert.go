package growhelper

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

const ErrorValue int = -999

func ReadReport(r string) Report {
	report := Report{}
	err := json.Unmarshal([]byte(r), &report)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	return report
}

func ReadString(r io.ReadCloser) string {
	body, err := io.ReadAll(r)
	defer r.Close()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	return string(body)
}

func ToInt(string string) int {
	value, err := strconv.Atoi(strings.TrimSpace(string))
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	return value
}

func ToString(int int) string {
	return strconv.Itoa(int)
}

func ToFloat64(string string) float64 {
	value, err := strconv.ParseFloat(strings.TrimSpace(string), 64)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	return value
}

//TODO: move interchanged structs here? or create own if diff very high. google shared structs -> makes binding very strong, bad habit?

type Report struct {
	Temperature Statistic
	Humidity    Statistic
}

type Statistic struct {
	Min     int
	Max     int
	Average int
}

func ToStatistic(measurements []int) Statistic {
	if len(measurements) == 0 {
		return Statistic{}
	}
	maxValue := measurements[0]
	minValue := measurements[0]
	movingAverageValue := float64(measurements[0])
	for _, m := range measurements {
		if m > maxValue {
			maxValue = m
		}
		if m < minValue {
			minValue = m
		}
		movingAverageValue = (movingAverageValue + float64(m)) / 2
	}
	return Statistic{Min: minValue, Max: maxValue, Average: int(math.Round(movingAverageValue))}
}
