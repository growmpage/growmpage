package growganizer

import (
	"fmt"
	"sort"
	"time"

	"github.com/growmpage/growhelper"
)

type TimeTableForToday struct {
	timeTableEntrys []TimeTableEntry
}

type TimeTableEntry struct {
	TimeInMinutes int
	Action        string
	Condition     Condition
}

type Growcontroller struct {
	week   Week
	Cancel chan struct{}
}

func (w Week) newTimeTableForToday() TimeTableForToday {
	timeTableForToday := TimeTableForToday{}
	minutesNoTime := growhelper.MinutesNoTime()
	for _, growcontrol := range w.Growcontrols {

		start := minutesNoTime
		if growcontrol.Time.Start != "" {
			start = growhelper.AddTodayParsedTime(growcontrol.Time.Start)
		}
		timeTableForToday.add(start, growcontrol.Action, growcontrol.Condition)

		end := growhelper.MinutesMidnight()
		if growcontrol.Time.End != "" {
			end = growhelper.AddTodayParsedTime(growcontrol.Time.End)
		}

		everyMinutes := growcontrol.Time.EveryMinutes
		if everyMinutes > 0 {
			if end < start {
				end = growhelper.MinutesMidnight()
			}
			for everyNow := start + everyMinutes; everyNow < end; everyNow = everyNow + everyMinutes {
				timeTableForToday.add(everyNow, growcontrol.Action, growcontrol.Condition)
			}
		}
	}
	sort.Sort(&timeTableForToday)
	return timeTableForToday
}

func (t *TimeTableForToday) add(timeInMinutes int, action string, condition Condition) {
	timeTableEntry := TimeTableEntry{
		TimeInMinutes: timeInMinutes,
		Action:        action,
		Condition:     condition,
	}
	t.timeTableEntrys = append(t.timeTableEntrys, timeTableEntry)
}

func NewTimer(minutes int) *time.Timer{
	untilExecuteDuration := (time.Minute * time.Duration(minutes)) - (time.Second * time.Duration(growhelper.SecondsNow()))
	return time.NewTimer(untilExecuteDuration)
}

func (t *TimeTableForToday) execute(g *Growcontroller) {
	timeTableIndex := t.getIndexToStart()
	if timeTableIndex == -1 {
		return
	}
	untilExecuteMinutes := t.timeTableEntrys[timeTableIndex].TimeInMinutes - growhelper.Minutes()
	actionTimer := NewTimer(untilExecuteMinutes)

	untilMidnightMinutes := growhelper.MinutesMidnight() - growhelper.Minutes()
	newDayTimer := NewTimer(untilMidnightMinutes)

	for {
		select {
		case <-actionTimer.C:
			t.timeTableEntrys[timeTableIndex].execute()
			if timeTableIndex == len(t.timeTableEntrys)-1 {
				untilExecuteMinutes = untilMidnightMinutes + 1 //TODO: just growhelper.Get("RESTARTGROWCONTROLLER")
			} else {
				timeTableIndex++
				untilExecuteMinutes = t.timeTableEntrys[timeTableIndex].TimeInMinutes - growhelper.Minutes()
			}
			actionTimer = NewTimer(untilExecuteMinutes)
		case <-newDayTimer.C:
			growhelper.Get("RESTARTGROWCONTROLLER")
			return
		case <-g.Cancel:
			actionTimer.Stop()
			newDayTimer.Stop()
			fmt.Println("end daily growcontroller")
			return
		}
	}
}

func (t *TimeTableForToday) getIndexToStart() int {
	sort.Sort(t)
	nowInMinutes := growhelper.Minutes()
	for i := range t.timeTableEntrys {
		if t.timeTableEntrys[i].TimeInMinutes >= nowInMinutes {
			return i
		}
	}
	return -1
}

func (t *TimeTableForToday) Len() int {
	return len(t.timeTableEntrys)
}

func (t *TimeTableForToday) Less(i, j int) bool {
	return t.timeTableEntrys[i].TimeInMinutes < t.timeTableEntrys[j].TimeInMinutes
}

func (t *TimeTableForToday) Swap(i, j int) {
	t.timeTableEntrys[i], t.timeTableEntrys[j] = t.timeTableEntrys[j], t.timeTableEntrys[i]
}
