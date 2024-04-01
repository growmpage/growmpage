package growtable

import (
	"math"
	"time"

	"github.com/growmpage/growhelper"
)

func (g *Growtable) Report(sinceHours int) growhelper.Report {
	temperatureArray := []int{}
	humidityArray := []int{}
	for _, m := range g.getMeasurementsInTime(sinceHours) {
		humidityArray = append(humidityArray, m.Humidity)
		temperatureArray = append(temperatureArray, m.Temperature)
	}
	return growhelper.Report{
		Temperature: growhelper.ToStatistic(temperatureArray),
		Humidity:    growhelper.ToStatistic(humidityArray),
	}
}

func (g *Growtable) getMeasurementsInTime(sinceHours int) []Measurement {
	thenInMinutes := growhelper.Minutes() - (sinceHours * 60)
	var measurementsInTime []Measurement
	for _, m := range g.Measurements {
		if thenInMinutes <= m.TimeInMinutes {
			measurementsInTime = append(measurementsInTime, m)
		}
	}
	return measurementsInTime
}

func (g *Growtable) UpdateAllColorDeviations(weeks []Week) {
	if len(g.Measurements) == 0 || len(weeks) == 0 {
		return
	}
	g.setAllColorDeviations(weeks)
}

func (g *Growtable) setAllColorDeviations(weeks []Week) {
	mi := g.setColorDeviationsBefore(weeks[0].Start)
	weekIndex := 0
	for mi >= 0 { //Calculate Measurements before next week start
		nextWeekIndex := weekIndex + 1
		if nextWeekIndex == len(weeks) {
			g.setColorDeviationsAfter(mi, weeks[weekIndex]) //all weeks have valid week.Start
			return
		}
		measurementBeforeNextWeekStart := g.getMeasurementTime(mi).Before(weeks[nextWeekIndex].Start)
		if measurementBeforeNextWeekStart {
			g.setColorDeviations(mi, weeks[weekIndex])
			mi--
		} else {
			weekIndex++
		}
	}
}

func (g *Growtable) setColorDeviationsBefore(date time.Time) int {
	for mi := len(g.Measurements) - 1; mi != -1; mi-- { //Measurements before first start date should be marked blue for not calculated
		if g.getMeasurementTime(mi).Before(date) {
			g.Measurements[mi].HumidityDeviation = 200
			g.Measurements[mi].TemperatureDeviation = 200
		} else {
			return mi
		}
	}
	return -1
}

func (g *Growtable) setColorDeviationsAfter(mi int, week Week) {
	for mi >= 0 { //Measurements after last week.Start
		g.setColorDeviations(mi, week)
		mi--
	}
}

func (g *Growtable) setColorDeviations(mi int, week Week) {
	g.Measurements[mi].HumidityDeviation = getHslValue(g.Measurements[mi].Humidity, week.Humidity, 20)
	g.Measurements[mi].TemperatureDeviation = getHslValue(g.Measurements[mi].Temperature, week.Temperature, 10)
}

func (g *Growtable) getMeasurementTime(i int) time.Time {
	return time.Unix(int64(g.Measurements[i].TimeInMinutes*60), 0)
}

func getHslValue(value int, optimal int, maxDiff int) int { //0=red, 100=green
	result := -1
	diff := int(math.Abs(float64(value - optimal)))
	if diff > maxDiff {
		result = 0
	}
	result = 100 - (diff * (100 / maxDiff))
	return result
}
func (m Measurement) TimeInDays() int {
	return growhelper.Days(m.TimeInMinutes)
}
func (m Measurement) TimeInString() string {
	return growhelper.Time(m.TimeInMinutes)
}
func (m Measurement) DateInString() string {
	return growhelper.Date(m.TimeInMinutes)
}
