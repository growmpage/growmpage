package growganizer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/growmpage/growhelper"
)

type Growganizer struct {
	Version         string
	ActiveWeekIndex int
	PlugControls    []PlugControl
	Weeks           []Week
	Alerts          []string
}

type PlugControl struct { //TODO: move to growswitcher
	Name             string
	Code             int
	PinNumberReceive int
	PinNumberSend    int
	ProtocolIndex    int
	Repeat           int
	PulseLength      int
	Length           int
}

type Week struct {
	Name         string
	Start        string
	Temperature  int
	Humidity     int
	EC           float64
	Dung         string
	LightMode    string
	LightHours   int
	Done         string
	Todo         string
	Growcontrols []Growcontrol
}

func ReadFromFile() Growganizer {
	file, err := os.ReadFile(growhelper.Filename_growganizer)
	if err != nil {
		fmt.Println(err)
	}
	growganizer := &Growganizer{}
	err = json.Unmarshal([]byte(file), &growganizer)
	if err != nil {
		fmt.Println(err)
	}
	return *growganizer
}

func (g *Growganizer) SaveToFile() {
	file, err := json.MarshalIndent(g, "", " ")
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile(growhelper.Filename_growganizer, file, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func (g *Growganizer) UpdateWeek(w http.ResponseWriter, r *http.Request) {
	control := r.FormValue("control")
	selectedWeekIndex := growhelper.ToInt(r.FormValue("selectedWeekIndex"))
	growcontrolIndex := growhelper.ToInt(r.FormValue("growcontrolIndex"))
	if selectedWeekIndex >= len(g.Weeks) {
		fmt.Println("ERROR: selectedWeekIndex out of range")
		return
	}
	switch control {
	case "AddGrowcontrol":
		if growcontrolIndex <= len(g.Weeks[selectedWeekIndex].Growcontrols) {
			g.Weeks[selectedWeekIndex].Growcontrols = append(g.Weeks[selectedWeekIndex].Growcontrols[:growcontrolIndex+1], g.Weeks[selectedWeekIndex].Growcontrols[growcontrolIndex:]...)
		} else {
			fmt.Println("ERROR: growcontrolIndex out of range")
		}
	case "DeleteGrowcontrol":
		if growcontrolIndex < len(g.Weeks[selectedWeekIndex].Growcontrols) {
			g.Weeks[selectedWeekIndex].Growcontrols = append(g.Weeks[selectedWeekIndex].Growcontrols[:growcontrolIndex], g.Weeks[selectedWeekIndex].Growcontrols[growcontrolIndex+1:]...)
		} else {
			fmt.Println("ERROR: growcontrolIndex out of range")
		}
		if len(g.Weeks[selectedWeekIndex].Growcontrols) == 0 {
			g.Weeks[selectedWeekIndex].Growcontrols = []Growcontrol{{Action: measure}}
		}
	case "AddWeek":
		g.addWeek(selectedWeekIndex)
	case "DeleteWeek":
		g.deleteWeek(selectedWeekIndex)
	case "ActivateWeek":
		g.ActiveWeekIndex = selectedWeekIndex
		g.Weeks[selectedWeekIndex].Start = growhelper.Today()
		for i := range g.Weeks {
			if i > selectedWeekIndex {
				g.Weeks[i].Start = ""
			}
		}
	default:
		fmt.Printf("ERROR: could not parse %v\n", control)
	}
	http.Redirect(w, r, "/week#"+growhelper.ToString(selectedWeekIndex), http.StatusSeeOther)
}

func (g *Growganizer) RawUpdateGrowganizer(growganizerJson string) {
	growganizer := &Growganizer{}
	err := json.Unmarshal([]byte(growganizerJson), &growganizer)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	*g = *growganizer
}

func (g *Growganizer) addWeek(selectedWeekIndex int) {
	g.Weeks = append(g.Weeks[:selectedWeekIndex+1], g.Weeks[selectedWeekIndex:]...)
	g.Weeks[selectedWeekIndex+1].Name = g.Weeks[selectedWeekIndex+1].Name + " (copy)"
	if selectedWeekIndex < g.ActiveWeekIndex {
		g.ActiveWeekIndex++
	}
}

func (g *Growganizer) deleteWeek(selectedWeekIndex int) {
	g.Weeks = append(g.Weeks[:selectedWeekIndex], g.Weeks[selectedWeekIndex+1:]...)
	if selectedWeekIndex < g.ActiveWeekIndex {
		g.ActiveWeekIndex--
	}
}

func (g *Growganizer) DeleteAlert(message string) {
	fmt.Printf("DeleteAlert: %v\n", message)
	for i, alert := range g.Alerts {
		if alert == message {
			g.Alerts = append(g.Alerts[:i], g.Alerts[i+1:]...)
		}
	}
	if len(g.Alerts) == 0 {
		g.Alerts = []string{}
	}
}

func formToWeek(r *http.Request, growcontrols []Growcontrol, weekIndex int) Week {
	weekIndexSuffix := "-" + growhelper.ToString(weekIndex)
	week := Week{
		r.FormValue("Name" + weekIndexSuffix),
		r.FormValue("Start" + weekIndexSuffix),
		growhelper.ToInt(r.FormValue("Temperature" + weekIndexSuffix)),
		growhelper.ToInt(r.FormValue("Humidity" + weekIndexSuffix)),
		growhelper.ToFloat64(r.FormValue("EC" + weekIndexSuffix)),
		r.FormValue("Dung" + weekIndexSuffix),
		r.FormValue("LightMode" + weekIndexSuffix),
		growhelper.ToInt(r.FormValue("LightHours" + weekIndexSuffix)),
		r.FormValue("Done" + weekIndexSuffix),
		r.FormValue("Todo" + weekIndexSuffix),
		growcontrols,
	}
	return week
}

func (g *Growganizer) UpdateWeeksByForm(r *http.Request) {
	for weekIndex := 0; weekIndex < len(g.Weeks); weekIndex++ {
		growcontrols := formToGrowcontrols(r, weekIndex)
		week := formToWeek(r, growcontrols, weekIndex)
		g.Weeks[weekIndex] = week
	}
}

func formToGrowcontrols(r *http.Request, weekIndex int) []Growcontrol {
	growcontrols := []Growcontrol{}
	i := 0
	growcontrolIndexSuffix := "-" + growhelper.ToString(i)
	weekIndexSuffix := "-" + growhelper.ToString(weekIndex)
	Action := r.FormValue("Action" + weekIndexSuffix + growcontrolIndexSuffix)
	for Action != "" {
		growcontrol := Growcontrol{
			Action,
			Time{
				r.FormValue("Start" + weekIndexSuffix + growcontrolIndexSuffix),
				r.FormValue("End" + weekIndexSuffix + growcontrolIndexSuffix),
				growhelper.ToInt(r.FormValue("EveryMinutes" + weekIndexSuffix + growcontrolIndexSuffix)),
			},
			Condition{
				r.FormValue("Value" + weekIndexSuffix + growcontrolIndexSuffix),
				growhelper.ToInt(r.FormValue("SinceHours" + weekIndexSuffix + growcontrolIndexSuffix)),
				r.FormValue("ComparisonSign" + weekIndexSuffix + growcontrolIndexSuffix),
				growhelper.ToInt(r.FormValue("ComparedTo" + weekIndexSuffix + growcontrolIndexSuffix)),
			},
		}
		growcontrols = append(growcontrols, growcontrol)
		i++
		growcontrolIndexSuffix = "-" + growhelper.ToString(i)
		Action = r.FormValue("Action" + weekIndexSuffix + growcontrolIndexSuffix)
	}
	return growcontrols
}

func (g *Growganizer) GetPlugControls() []string {
	names := []string{}
	for _, plug := range g.PlugControls {
		names = append(names, plug.Name)
	}
	return names
}

func (g *Growganizer) PlugIndex(name string) int {
	for i, plug := range g.PlugControls {
		if name == plug.Name {
			return i
		}
	}
	return -1 //TODO: google best practice
}

func (g *Growganizer) GetOptimalConditionsByWeekStarts() map[string]map[string]int {
	weekMap := make(map[string]map[string]int)
	for _, week := range g.Weeks {
		_, err := time.ParseInLocation("2006-01-02", week.Start, time.Local) //TODO: move to growtime
		if err == nil {
			weekMap[week.Start] = make(map[string]int)
			weekMap[week.Start]["Temperature"] = week.Temperature
			weekMap[week.Start]["Humidity"] = week.Humidity
		} else {
			fmt.Printf("err: %v\n", err)
		}
	}
	return weekMap
}
