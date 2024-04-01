package growganizer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/growmpage/growhelper"
)

func TestWeek_createGrowcontrolPlan(t *testing.T) {
	type fields struct {
		Growcontrols []Growcontrol
	}
	tests := []struct {
		name   string
		fields fields
		want   TimeTableForToday
	}{
		{
			fields: fields{Growcontrols: []Growcontrol{}},
			want:   TimeTableForToday{},
		},
		{
			fields: fields{Growcontrols: []Growcontrol{Growcontrol{Time: Time{Start: "06:00"}, Action: growsensorActions[1]}, Growcontrol{Time: Time{Start: "06:15"}, Action: growsensorActions[0]}}},
			want:   TimeTableForToday{timeTableEntrys: []TimeTableEntry{TimeTableEntry{TimeInMinutes: growhelper.MinutesNoTime() + (6 * 60), Action: "Picture"}, TimeTableEntry{TimeInMinutes: growhelper.MinutesNoTime() + (6*60 + 15), Action: "Measure"}}},
		},
		{
			fields: fields{Growcontrols: []Growcontrol{Growcontrol{Time: Time{Start: "06:00"}, Action: growsensorActions[2]}}},
			want:   TimeTableForToday{timeTableEntrys: []TimeTableEntry{TimeTableEntry{TimeInMinutes: growhelper.MinutesNoTime() + (6 * 60), Action: "Alert"}}},
		},
		{
			fields: fields{Growcontrols: []Growcontrol{Growcontrol{Time: Time{Start: "06:00", EveryMinutes: 345, End: "18:00"}, Action: growsensorActions[2]}}},
			want: TimeTableForToday{timeTableEntrys: []TimeTableEntry{
				TimeTableEntry{TimeInMinutes: growhelper.MinutesNoTime() + (6 * 60), Action: "Alert"},
				TimeTableEntry{TimeInMinutes: growhelper.MinutesNoTime() + (6*60 + 345), Action: "Alert"},
				TimeTableEntry{TimeInMinutes: growhelper.MinutesNoTime() + (6*60 + 690), Action: "Alert"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := Week{
				Growcontrols: tt.fields.Growcontrols,
			}
			if got := w.newTimeTableForToday(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Week.createGrowcontrolPlan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeTableForToday_getIndexToStart(t *testing.T) {
	type fields struct {
		timeTableEntrys []TimeTableEntry
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			fields: fields{timeTableEntrys: []TimeTableEntry{
				TimeTableEntry{TimeInMinutes: growhelper.MinutesNoTime() + (6 * 60), Action: "Alert"},
				TimeTableEntry{TimeInMinutes: growhelper.MinutesNoTime() + (24 * 60), Action: "Alert"},
				TimeTableEntry{TimeInMinutes: growhelper.MinutesNoTime() + (24 * 61), Action: "Alert"},
			}},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TimeTableForToday{
				timeTableEntrys: tt.fields.timeTableEntrys,
			}
			if got := tr.getIndexToStart(); got != tt.want {
				t.Errorf("TimeTableForToday.getIndexToStart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeTableEntry_value(t *testing.T) {
	type fields struct {
		TimeInMinutes int
		Action        string
		Condition     Condition
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			fields: fields{Condition: Condition{Value: "Temperature", SinceHours: 5}},
			want:   21,
		},
		{
			fields: fields{Condition: Condition{Value: "MaxHumidity", SinceHours: 5}},
			want:   60,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fmt.Printf("growhelper.ReadString(r.Body): %v\n", growhelper.ReadString(r.Body))
		if r.URL.Path != "/GROWTABLEREPORT" {
			t.Errorf("Expected to request '/GROWTABLEREPORT', got: %s", r.URL.Path)
		}
		sinceHours := growhelper.ToInt(growhelper.ReadString(r.Body))
		if sinceHours != 5 {
			t.Error("sinceHours not 5")
		}
		report := growhelper.Report{
			Temperature: growhelper.Statistic{Min: 20, Max: 22, Average: 21},
			Humidity:    growhelper.Statistic{Min: 50, Max: 60, Average: 55},
		}
		responseBytes, _ := json.Marshal(report)
		fmt.Printf("string(responseBytes): %v\n", string(responseBytes))
		fmt.Fprint(w, string(responseBytes))
	}))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TimeTableEntry{
				TimeInMinutes: tt.fields.TimeInMinutes,
				Action:        tt.fields.Action,
				Condition:     tt.fields.Condition,
			}
			fmt.Printf("url: %v\n", server.URL)
			growhelper.Port = ":" + strings.Split(server.URL, ":")[2]
			if got := tr.value(); got != tt.want {
				t.Errorf("TimeTableEntry.value() = %v, want %v", got, tt.want)
			}

		})
	}
	defer server.Close()
}

func TestTimeTableEntry_execute(t *testing.T) {
	type fields struct {
		TimeInMinutes int
		Action        string
		Condition     Condition
	}
	tests := []struct {
		name   string
		fields fields
		server *httptest.Server
	}{
		{
			fields: fields{Action: "3", Condition: Condition{Value: temperature, SinceHours: 5, ComparisonSign: smaller, ComparedTo: 1}},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/GROWTABLEREPORT" {
					return
				}
				if r.URL.Path != "/PlugControl" {
					t.Errorf("Expected '/PlugControl', got: %s", r.URL.Path)
				}
				body := growhelper.ToInt(growhelper.ReadString(r.Body))
				if body != 3 {
					t.Errorf("Expected 3, got: %v", body)
				}
			})),
		},
		{
			fields: fields{Action: alert, Condition: Condition{Value: temperature, SinceHours: 5, ComparisonSign: smaller, ComparedTo: 1}},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/GROWTABLEREPORT" {
					return
				}
				if r.URL.Path != "/ALERT" {
					t.Errorf("Expected '/ALERT', got: %s", r.URL.Path)
				}
				body := growhelper.ReadString(r.Body)

				expect := growhelper.DateTime() + " -> Temperature smaller than 1 (since 5 hours)"
				expect = growhelper.DateToGerman(expect)
				if body != expect {
					t.Errorf("Expected %v, got: %v", expect, body)
				}
			})),
		},
		{
			fields: fields{Action: measure, Condition: Condition{Value: temperature, SinceHours: 5, ComparisonSign: smaller, ComparedTo: 1}},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/GROWTABLEREPORT" {
					return
				}
				if r.URL.Path != "/MEASURE" {
					t.Errorf("Expected '/MEASURE', got: %s", r.URL.Path)
				}
			})),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TimeTableEntry{
				TimeInMinutes: tt.fields.TimeInMinutes,
				Action:        tt.fields.Action,
				Condition:     tt.fields.Condition,
			}
			growhelper.Port = ":" + strings.Split(tt.server.URL, ":")[2]
			tr.execute()
		})
	}
}

func TestTimeTableEntry_createAlertMessage(t *testing.T) {
	type fields struct {
		TimeInMinutes int
		Action        string
		Condition     Condition
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{Action: measure, Condition: Condition{Value: temperature, SinceHours: 5, ComparisonSign: smaller, ComparedTo: 1}},
			want:   growhelper.DateToGerman(growhelper.DateTime()) + " -> Temperature smaller than 1 (since 5 hours)",
		},
		{
			fields: fields{Condition: Condition{Value: "AnyCondition", SinceHours: 5, ComparisonSign: smaller, ComparedTo: 0}},
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TimeTableEntry{
				TimeInMinutes: tt.fields.TimeInMinutes,
				Action:        tt.fields.Action,
				Condition:     tt.fields.Condition,
			}
			if got := tr.createAlertMessage(); got != tt.want {
				t.Errorf("TimeTableEntry.createAlertMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
