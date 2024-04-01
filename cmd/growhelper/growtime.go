package growhelper

import (
	"fmt"
	"strings"
	"time"
)

func getNow() time.Time{
	return time.Now().Local()
}

func getTime(minutes int) time.Time{
	return time.Unix(int64(minutes*60), 0).Local()
}


func HtmlDate(date string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", date, time.Local)
}

func HtmlDateAsTime(date string) time.Time {
	t, err := HtmlDate(date)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	return t
}

func DaysSinceHtmlDate(date string) int {
	days := int(getNow().Sub(HtmlDateAsTime(date)).Abs().Hours() / 24)
	return days
}

func Days(minutes int) int {
	return minutes/60/24
}

func Time(minutes int) string {
	return getTime(minutes).Format("15:04")
}

func Date(minutes int) string {
	return getTime(minutes).Format("Mon 2 Jan 2006")
}

func DaysToMinutes(days int) int {
	return days * 24 * 60
}

func MinutesNoTime() int {
	minutesOver := (getNow().Hour() * 60) + getNow().Minute()
	return Minutes() - minutesOver
}
func Minutes() int {
	return int(getNow().Unix()) / 60
}
func MinutesFromTime(t time.Time) int {
	return int(t.Local().Unix()) / 60
}
func MinutesMidnight() int {
	return MinutesNoTime() + (24 * 60)
}

func DateTime() string {
	return getNow().Format("Mon 2.1, 15:04")
}

func AddTodayParsedTime(timeFormat24 string) int {
	minutesNoTime := MinutesNoTime()
	hours := ToInt(strings.Split(timeFormat24, ":")[0])
	minutes := ToInt(strings.Split(timeFormat24, ":")[1])
	minutesToAdd := (hours * 60) + minutes
	result := minutesNoTime + minutesToAdd
	return result
}

func Today() string {
	return getNow().Format("2006-01-02")
}

func DateToGerman(message string) string {
	message = strings.Replace(message, "Mon", "Montag", 1)
	message = strings.Replace(message, "Tue", "Dienstag", 1)
	message = strings.Replace(message, "Wed", "Mittwoch", 1)
	message = strings.Replace(message, "Thu", "Donnerstag", 1)
	message = strings.Replace(message, "Fri", "Freitag", 1)
	message = strings.Replace(message, "Sat", "Samstag", 1)
	message = strings.Replace(message, "Son", "Sonntag", 1)
	return message
}
