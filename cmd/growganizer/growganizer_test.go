package growganizer

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestReadFromFile(t *testing.T) { // TODO: add ../ for all paths used here, or better: change vscode test working directory
	_, filename, _, _ := runtime.Caller(0)
	fmt.Printf("main: %s", filename)
	tests := []struct {
		name     string
		contains string
	}{
		{name: "growganizer", contains: "Succesfull tested my growbox with all controls for a week"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadFromFile(); !strings.Contains(fmt.Sprintf("%v", got), tt.contains) {
				t.Errorf("ReadFromFile() = %v, want containing in %v", got, tt.contains)
			}
		})
	}
}

func TestGrowganizer_SaveToFile(t *testing.T) {
	before := ReadFromFile()
	modified := append(before.Weeks, Week{Name: "foo"})
	tests := []struct {
		g Growganizer
	}{
		{g: Growganizer{Weeks: modified}},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			tt.g.SaveToFile()
			before.SaveToFile()
		})
	}
}

func TestGrowganizer_UpdateWeek(t *testing.T) {
	type fields struct {
		ActiveWeekIndex int
		Weeks           []Week
	}
	tests := []struct {
		name   string
		fields fields
		args   string
	}{
		{
			fields: fields{Weeks: []Week{Week{Name: "test", EC: 1}}, ActiveWeekIndex: 0},
			args:   "control=bla&selectedWeekIndex=0&growcontrolIndex=0",
		},
		{
			fields: fields{Weeks: []Week{Week{Name: "test2", Growcontrols: []Growcontrol{Growcontrol{Action: "foo"}}}}, ActiveWeekIndex: 0},
			args:   "control=AddGrowcontrol&selectedWeekIndex=0&growcontrolIndex=0&",
		},
		{
			fields: fields{Weeks: []Week{Week{Name: "test", EC: 1}}, ActiveWeekIndex: 0},
			args:   "control=DeleteGrowcontrol&selectedWeekIndex=0&growcontrolIndex=0",
		},
		{
			fields: fields{Weeks: []Week{Week{Name: "test", EC: 1}}},
			args:   "control=AddWeek&selectedWeekIndex=0&growcontrolIndex=0",
		},
		{
			fields: fields{Weeks: []Week{Week{Name: "test", EC: 1}}, ActiveWeekIndex: 0},
			args:   "control=DeleteWeek&selectedWeekIndex=0&growcontrolIndex=0",
		},
		{
			fields: fields{Weeks: []Week{Week{Name: "test", EC: 1}, Week{Name: "test2", EC: 1}}, ActiveWeekIndex: 0},
			args:   "control=ActivateWeek&selectedWeekIndex=0&growcontrolIndex=0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Growganizer{
				ActiveWeekIndex: tt.fields.ActiveWeekIndex,
				Weeks:           tt.fields.Weeks,
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.args))
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			g.UpdateWeek(w, r)
			if tt.fields.ActiveWeekIndex != g.ActiveWeekIndex {
				t.Errorf("ActiveWeekIndex = %v, want %v", g.ActiveWeekIndex, tt.fields.ActiveWeekIndex)
			}
			result := w.Result().Header
			if len(result) != 0 && result.Values("Location")[0] != "/week#0" {
				t.Errorf("redirect = %v, want %v", w.Result().Header.Values("Location")[0], "/week#0")
			}

		})
	}
}

func TestGrowganizer_RawUpdateGrowganizer(t *testing.T) {
	type fields struct {
		Version         string
		ActiveWeekIndex int
		PlugControls    []PlugControl
		Weeks           []Week
		Alerts          []string
	}
	type args struct {
		growganizerJson string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			fields: fields{Version: "new", ActiveWeekIndex: 0, Weeks: []Week{Week{Name: "blu"}}},
			args:   args{growganizerJson: "{\"Version\": \"8a2643e\"}"},
		},
		{
			fields: fields{Version: "new", ActiveWeekIndex: 0, Weeks: []Week{Week{Name: "blu"}}},
			args:   args{growganizerJson: "{\"Version\": \"8a2643e\",\"PlugControls\": [{\"Name\": \"SimOnC\",\"Code\": 5201,\"PinNumberReceive\": 27,\"PinNumberSend\": 17,\"ProtocolIndex\": 1,\"Repeat\": 10,\"PulseLength\": 307,\"Length\": 24}],\"Weeks\": [{\"Name\": \"Testweek\"}]}"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Growganizer{
				Version:         tt.fields.Version,
				ActiveWeekIndex: tt.fields.ActiveWeekIndex,
				PlugControls:    tt.fields.PlugControls,
				Weeks:           tt.fields.Weeks,
				Alerts:          tt.fields.Alerts,
			}
			g.RawUpdateGrowganizer(tt.args.growganizerJson)
			if g.Version != "8a2643e" {
				t.Errorf("Version %v, want %v", g.Version, "8a2643e")
			}
		})
	}
}

func TestGrowganizer_UpdateWeeksByForm(t *testing.T) {
	type fields struct {
		Version         string
		ActiveWeekIndex int
		PlugControls    []PlugControl
		Weeks           []Week
		Alerts          []string
	}
	tests := []struct {
		name   string
		fields fields
		args   string
	}{
		{
			fields: fields{Weeks: []Week{Week{Name: "test", EC: 1}}, ActiveWeekIndex: 0},
			args:   "Name-0=crazy&control=bla&selectedWeekIndex=0&growcontrolIndex=0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Growganizer{
				Version:         tt.fields.Version,
				ActiveWeekIndex: tt.fields.ActiveWeekIndex,
				PlugControls:    tt.fields.PlugControls,
				Weeks:           tt.fields.Weeks,
				Alerts:          tt.fields.Alerts,
			}
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.args))
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			g.UpdateWeeksByForm(r)
			if g.Weeks[0].Name != "crazy" {
				t.Fail()
			}
			fmt.Printf("g.Weeks: %v\n", g.Weeks)
		})
	}
}

func TestGrowganizer_GetPlugControls(t *testing.T) {
	type fields struct {
		Version         string
		ActiveWeekIndex int
		PlugControls    []PlugControl
		Weeks           []Week
		Alerts          []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			fields: fields{Weeks: []Week{Week{Name: "test", EC: 1}}, ActiveWeekIndex: 0, PlugControls: []PlugControl{PlugControl{Name: "bla", Code: 111}, PlugControl{Name: "blu", Code: 111}}},
			want:   []string{"bla", "blu"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Growganizer{
				Version:         tt.fields.Version,
				ActiveWeekIndex: tt.fields.ActiveWeekIndex,
				PlugControls:    tt.fields.PlugControls,
				Weeks:           tt.fields.Weeks,
				Alerts:          tt.fields.Alerts,
			}
			if got := g.GetPlugControls(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Growganizer.GetPlugControls() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrowganizer_PlugIndex(t *testing.T) {
	type fields struct {
		Version         string
		ActiveWeekIndex int
		PlugControls    []PlugControl
		Weeks           []Week
		Alerts          []string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			fields: fields{Weeks: []Week{Week{Name: "test", EC: 1}}, ActiveWeekIndex: 0, PlugControls: []PlugControl{PlugControl{Name: "bla", Code: 111}, PlugControl{Name: "blu", Code: 111}}},
			want:   1,
			args:   args{name: "blu"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Growganizer{
				Version:         tt.fields.Version,
				ActiveWeekIndex: tt.fields.ActiveWeekIndex,
				PlugControls:    tt.fields.PlugControls,
				Weeks:           tt.fields.Weeks,
				Alerts:          tt.fields.Alerts,
			}
			if got := g.PlugIndex(tt.args.name); got != tt.want {
				t.Errorf("Growganizer.PlugIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrowganizer_GetOptimalConditionsByWeekStarts(t *testing.T) {
	weekMap := make(map[string]map[string]int)
	weekMap["2006-01-02"] = make(map[string]int)
	weekMap["2006-01-02"]["Temperature"] = 10
	weekMap["2006-01-02"]["Humidity"] = 40

	type fields struct {
		Version         string
		ActiveWeekIndex int
		PlugControls    []PlugControl
		Weeks           []Week
		Alerts          []string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			fields: fields{Weeks: []Week{Week{Start: "2006-01-02", Name: "test", EC: 1, Temperature: 10, Humidity: 40}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Growganizer{
				Version:         tt.fields.Version,
				ActiveWeekIndex: tt.fields.ActiveWeekIndex,
				PlugControls:    tt.fields.PlugControls,
				Weeks:           tt.fields.Weeks,
				Alerts:          tt.fields.Alerts,
			}
			if got := g.GetOptimalConditionsByWeekStarts(); !reflect.DeepEqual(got, weekMap) {
				t.Errorf("Growganizer.GetOptimalConditionsByWeekStarts() = %v, want %v", got, weekMap)
			}
		})
	}
}
