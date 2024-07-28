package growtable

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/growmpage/growhelper"
)

type Growtable struct {
	mu           sync.Mutex
	Measurements []Measurement
}

type Measurement struct {
	HumidityDeviation    int
	TemperatureDeviation int
	TimeInMinutes        int
	Temperature          int
	Humidity             int
	Picture              bool
}

func ReadFromFile() *Growtable {
	file, _ := os.ReadFile(growhelper.Filename_growtable)
	growtable := &Growtable{}
	_ = json.Unmarshal([]byte(file), &growtable)
	sort.Sort(growtable)
	return growtable
}

func (g *Growtable) SaveToFile() {
	// g.withLockContext(func() {
	file, err := json.MarshalIndent(g, "", " ")
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile(growhelper.Filename_growtable, file, 0644)
	if err != nil {
		fmt.Println(err)
	}
	// })
}

func (g *Growtable) RawUpdateGrowtable(growtableJson string, removePictures bool) {
	err := json.Unmarshal([]byte(growtableJson), &g)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		growhelper.Post("ALERT", err.Error())
	}
	if removePictures {
		g.removeUndocumentedPictures() //TODO: new function -> cleanup button?
	}
}

func (g *Growtable) removeUndocumentedPictures() {
	base := growhelper.Filename_pictures
	files, _ := os.ReadDir(base)
	for _, file := range files {
		// fmt.Printf("file.Name(): %v\n", file.Name())
		if file.Name() == "toCopy.png" {
			continue
		}
		f := base + file.Name()
		minutes := growhelper.ToInt(strings.Split(file.Name(), ".")[0])
		fmt.Printf("think of deleting %v with minutes: %v and path %v\n", file.Name(), minutes, f)
		if !g.hasPicture(minutes) {
			fmt.Println("remove " + f)
			os.Remove(f)
		}
	}
}

func (g *Growtable) AddPicture(minute int, week Week) {
	if len(g.Measurements) == 0 {
		g.Measurements = []Measurement{{
			TimeInMinutes: minute,
		}}
	}

	sameMinute := g.Measurements[0].TimeInMinutes == minute
	if sameMinute {
		g.Measurements[0].Picture = true
	} else {
		g.Measurements = append(
			[]Measurement{{
				TimeInMinutes: minute,
				Picture:       true,
			}}, g.Measurements...)
	}

	growhelper.Get("UPDATE")
}

func (g *Growtable) AddMeasurement(temperature int, humidity int, week Week) {
	var minute int = growhelper.Minutes()
	if len(g.Measurements) == 0 {
		g.Measurements = []Measurement{{
			TimeInMinutes: minute,
		}}
	}

	diffMinute := g.Measurements[0].TimeInMinutes - minute
	if math.Abs(float64(diffMinute)) < 2 {
		g.Measurements[0].Temperature = temperature
		g.Measurements[0].Humidity = humidity
	} else {
		g.Measurements = append(
			[]Measurement{{
				TimeInMinutes: minute,
				Temperature:   temperature,
				Humidity:      humidity,
			}}, g.Measurements...)
	}

	g.setColorDeviations(0, week)
	growhelper.Get("UPDATE")
}

type Week struct {
	Start       time.Time
	Temperature int
	Humidity    int
}

// func (g *Growtable) withLockContext(fn func()) { //TODO: really?
// 	g.mu.Lock()
// 	defer g.mu.Unlock()
// 	fn()
// }

func (g *Growtable) Len() int {
	return len(g.Measurements)
}

func (g *Growtable) Less(i, j int) bool {
	return g.Measurements[i].TimeInMinutes > g.Measurements[j].TimeInMinutes
}

func (g *Growtable) Swap(i, j int) {
	g.Measurements[i], g.Measurements[j] = g.Measurements[j], g.Measurements[i]
}

func (g *Growtable) hasPicture(minutes int) bool {
	for _, m := range g.Measurements {
		if m.TimeInMinutes == minutes {
			return m.Picture
		}
	}
	return false
}
