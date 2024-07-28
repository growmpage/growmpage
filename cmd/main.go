package main

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"

	"github.com/growmpage/growganizer"
	"github.com/growmpage/growhelper"
	"github.com/growmpage/growsensoric"
	"github.com/growmpage/growtable"
)

type Growmpage struct {
	growcache *bytes.Buffer

	growganizer    growganizer.Growganizer
	growcontroller chan struct{}
	growtable      *growtable.Growtable

	growobserver *growsensoric.Growobserver
	growswitcher *growsensoric.Growswitcher
}

func main() {
	grow := &Growmpage{
		growganizer:  growganizer.ReadFromFile(),
		growtable:    growtable.ReadFromFile(),
		growcache:    new(bytes.Buffer),
		growobserver: &growsensoric.Growobserver{},
		growswitcher: &growsensoric.Growswitcher{},
	}

	grow.handleGrowmpage()
	grow.handleExpertpage()
	grow.handleInternal()
	grow.handleFiles()

	grow.restartGrowcontroller()
	grow.updateGrowmpage()
	grow.startServer()
}

func (g *Growmpage) startServer() {
	fmt.Printf("open: %v\n", growhelper.Url(""))
	err := http.ListenAndServe(growhelper.Port, nil)
	if err != nil {
		fmt.Printf("ListenAndServe error: %v\n", err)
	}
}

func (g *Growmpage) SaveToDatabase() {
	g.growganizer.SaveToFile()
	g.growtable.SaveToFile()
}

func (g *Growmpage) updateGrowmpage() {
	tmpl := template.Must(template.ParseFiles(growhelper.Filename_growmpage))
	g.growcache = new(bytes.Buffer) //TODO: must be in struct?
	err := tmpl.Execute(g.growcache, g)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}

func (g *Growmpage) restartGrowcontroller() {
	if g.growcontroller != nil {
		close(g.growcontroller)
	}
	g.growcontroller = *g.growganizer.NewDailyGrowcontrol(g.Weeks()[g.growganizer.ActiveWeekIndex])
}

func (g *Growmpage) updateGrowtableColors() {
	g.growtable.UpdateAllColorDeviations(g.growganizerToGrowtableWeeks())
}

func (g *Growmpage) currentWeek() growtable.Week {
	weeks := g.growganizerToGrowtableWeeks()
	if len(weeks) == 0 {
		return growtable.Week{}
	}
	return weeks[len(weeks)-1] //TODO: really?
}

func (g *Growmpage) growganizerToGrowtableWeeks() []growtable.Week {
	weeks := []growtable.Week{}
	for _, week := range g.growganizer.Weeks {
		parsedStart, err := growhelper.HtmlDate(week.Start)
		if err != nil {
			continue
		}
		weeks = append(weeks,
			growtable.Week{
				Start:       parsedStart,
				Temperature: week.Temperature,
				Humidity:    week.Humidity,
			},
		)
	}
	return weeks
}

// these fields are used by growmpage.html (must be public, extending *Growmpage, TODO: add all here eg. TimeInDays):
func (g *Growmpage) GetDaysSinceActiveStartDate() int {
	if g.growganizer.Weeks[g.growganizer.ActiveWeekIndex].Start == "" {
		return 0
	}
	days := growhelper.DaysSinceHtmlDate(g.growganizer.Weeks[g.growganizer.ActiveWeekIndex].Start)
	return days
}
func (g *Growmpage) ActiveWeekIndex() int {
	return g.growganizer.ActiveWeekIndex
}
func (g *Growmpage) Weeks() []growganizer.Week {
	return g.growganizer.Weeks
}
func (g *Growmpage) Measurements() []growtable.Measurement {
	return g.growtable.Measurements
}
func (g *Growmpage) GetDummyGrowtableEntrys() []struct{} {
	return make([]struct{}, 60)
}
func (g *Growmpage) Alerts() []string {
	return g.growganizer.Alerts
}
func (g *Growmpage) GetActions() []string {
	return g.growganizer.GetActions()
}
func (g *Growmpage) GetAnyConditionString() string {
	return g.growganizer.GetAnyConditionString()
}
