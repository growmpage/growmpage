package growganizer

import (
	"fmt"

	"github.com/growmpage/growhelper"
)

var (
	growsensorActions = []string{measure, picture, alert} //TODO: convert to struct
	measure           = "Measure"
	picture           = "Picture"
	alert             = "Alert"

	growtableValues = []string{anyCondition, humidity, temperature, countPastDays, minTemperature, maxTemperature, minHumidity, maxHumidity}
	anyCondition    = "AnyCondition"
	humidity        = "Humidity"
	temperature     = "Temperature"
	countPastDays   = "CountPastDays"
	minTemperature  = "MinTemperature"
	maxTemperature  = "MaxTemperature"
	minHumidity     = "MinHumidity"
	maxHumidity     = "MaxHumidity"

	comparisonSigns = []string{bigger, smaller}
	bigger          = "bigger"
	smaller         = "smaller"
)

type Growcontrol struct {
	Action    string
	Time      Time
	Condition Condition
}

type Condition struct {
	Value          string
	SinceHours     int
	ComparisonSign string
	ComparedTo     int
}

type Time struct {
	Start        string
	End          string
	EveryMinutes int
}

func (g Growganizer) NewDailyGrowcontrol(week Week) *chan struct{} {
	growcontroller := &Growcontroller{
		week:   week,
		Cancel: make(chan struct{}),
	}
	growcontroller.startDailyGrowcontrol()
	return &growcontroller.Cancel
}

func (g *Growcontroller) startDailyGrowcontrol() {
	fmt.Printf("activate Growcontrol for week %v on %v\n", g.week.Name, growhelper.Today())
	g.Cancel = make(chan struct{})
	timeTableForToday := g.week.newTimeTableForToday()
	go timeTableForToday.execute(g)
}

func (t *TimeTableEntry) execute() {
	fmt.Printf("try to t.Action: %v\n", t.Action)
	if t.Condition.Value == anyCondition || t.Condition.Value == "" {
		execute(t)
		return
	}
	value := t.value()
	if value == growhelper.ErrorValue {
		fmt.Printf("%v value from growtable report ErrorValue\n", t.Action)
		return
	}
	switch t.Condition.ComparisonSign {
	case bigger:
		if value > t.Condition.ComparedTo {
			execute(t)
		}
	case smaller:
		if value < t.Condition.ComparedTo {
			execute(t)
		}
	default:
		if value == t.Condition.ComparedTo {
			execute(t)
		}
	}
}

func execute(t *TimeTableEntry) {
	switch t.Action {
	case measure:
		growhelper.Get("MEASURE")
	case picture:
		growhelper.Get("PICTURE")
	case alert:
		growhelper.Post("ALERT", t.createAlertMessage())
	default:
		growhelper.Post("PlugControl", t.Action)
	}
}

func (t *TimeTableEntry) createAlertMessage() string {
	message := growhelper.DateToGerman(growhelper.DateTime() + " -> ")
	switch t.Condition.Value {
	case "AnyCondition":
		return "" //not added to alert list
	case countPastDays:
		message += "Switch to next week!"
	default:
		message += t.Condition.Value + " " + t.Condition.ComparisonSign + " than " + growhelper.ToString(t.Condition.ComparedTo)
	}
	if t.Condition.SinceHours != 0 {
		message += " (since " + growhelper.ToString(t.Condition.SinceHours) + " hours)"
	}
	return message
}

func (t *TimeTableEntry) value() int {
	if t.Condition.Value == countPastDays {
		return growhelper.ToInt(growhelper.Get("DAYSSINCEACTIVESTARTDATE"))
	}
	sinceHours := growhelper.ToString(t.Condition.SinceHours)
	answer := growhelper.Post("GROWTABLEREPORT", sinceHours)
	report := growhelper.ReadReport(answer)
	var statistic growhelper.Statistic
	switch t.Condition.Value {
	case temperature, minTemperature, maxTemperature:
		statistic = report.Temperature
	case humidity, minHumidity, maxHumidity:
		statistic = report.Humidity
	}
	switch t.Condition.Value {
	case maxTemperature, maxHumidity:
		return statistic.Max
	case minTemperature, minHumidity:
		return statistic.Min
	case temperature, humidity:
		return statistic.Average
	default:
		return growhelper.ErrorValue
	}
}

func (g *Growganizer) GetActions() []string {
	return append(growsensorActions, g.GetPlugControls()...)
}

func (g *Growcontrol) GetComparisonSigns() []string {
	return comparisonSigns
}

func (g *Growcontrol) GetValues() []string {
	return growtableValues
}

func (g *Growganizer) GetAnyConditionString() string {
	return anyCondition
}
