package growtable

import (
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/growmpage/growhelper"
)

// func TestEmptyPictureDirectory(t *testing.T) { TODO:
// 	files, _ := os.ReadDir(growhelper.Filename_pictures)
// 	demoPictureName := "toCopy.png"
// 	ok := len(files) == 1 && files[0].Name() == demoPictureName
// 	if !ok {
// 		name := ""
// 		if len(files) > 0 {
// 			name = files[0].Name()
// 		}
// 		t.Errorf("%v must contain one file called %v\nBut was: %v and %v", growhelper.Filename_pictures, demoPictureName, len(files), name)
// 	}
// }

func TestReadFromFile(t *testing.T) {
	// growhelper.Filename_growtable = "../" + growhelper.Filename_growtable //TODO: change vscode testing working directory
	tests := []struct {
		name string
		want *Growtable
	}{
		// {name: "growtable", want: &Growtable{Measurements: []Measurement{}}}, //TODO: re-enable so growtable.json is always empty
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fmt.Sprintf("%#v", ReadFromFile())
			want := fmt.Sprintf("%#v", tt.want)
			if got != want {
				t.Errorf("ReadFromFile() = %v, want %v", got, want)
			}
		})
	}
}

func TestGrowtable_SaveToFile(t *testing.T) {
	before := ReadFromFile()
	modified := append(before.Measurements, Measurement{Humidity: 10})
	tests := []struct {
		g Growtable
	}{
		{g: Growtable{Measurements: modified}},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			tt.g.SaveToFile()
			before.SaveToFile()
		})
	}
}

func TestGrowtable_RawUpdate(t *testing.T) {
	tests := []struct {
		growtableJson string
		want          int
	}{
		{
			growtableJson: "{\"Measurements\": [ {\"Temperature\": 22}]}",
			want:          22,
		},
		{
			growtableJson: "{\"Measurements\": [ {\"Humidity\": 50}]}",
			want:          50,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			growtable := &Growtable{}
			growtable.RawUpdateGrowtable(tt.growtableJson, false)
			tempBad := growtable.Measurements[0].Temperature != tt.want
			humBad := growtable.Measurements[0].Humidity != tt.want
			if tempBad && humBad {
				t.Fail()
			}
		})
	}
}

func TestGrowtable_Report(t *testing.T) {
	type fields struct {
		mu           sync.Mutex
		Measurements []Measurement
	}
	type args struct {
		sinceHours int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   growhelper.Report
	}{
		{
			fields: fields{
				Measurements: []Measurement{
					{Temperature: 20, Humidity: 50, TimeInMinutes: growhelper.Minutes()},
					{Temperature: 22, Humidity: 60, TimeInMinutes: growhelper.Minutes()}},
			},
			args: args{sinceHours: 5},
			want: growhelper.Report{Temperature: growhelper.Statistic{Min: 20, Max: 22, Average: 21},
				Humidity: growhelper.Statistic{Min: 50, Max: 60, Average: 55}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Growtable{
				mu:           tt.fields.mu,
				Measurements: tt.fields.Measurements,
			}
			if got := g.Report(tt.args.sinceHours); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Growtable.Report() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrowtable_getMeasurementsInTime(t *testing.T) {
	type fields struct {
		mu           sync.Mutex
		Measurements []Measurement
	}
	type args struct {
		sinceHours int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Measurement
	}{
		{
			fields: fields{
				Measurements: []Measurement{
					{Temperature: 20, TimeInMinutes: growhelper.Minutes()},
					{Temperature: 24, TimeInMinutes: growhelper.Minutes() - 4*60},
					{Temperature: 15, TimeInMinutes: growhelper.Minutes() - (5 * 60) - 1},
					{Temperature: 25, TimeInMinutes: growhelper.Minutes() - (5 * 60)},
					{Temperature: 35, TimeInMinutes: growhelper.Minutes() - (5 * 60) + 1},
					{Temperature: 26, TimeInMinutes: growhelper.Minutes() - 6*60},
				},
			},
			args: args{sinceHours: 5},
			want: []Measurement{
				{Temperature: 20, TimeInMinutes: growhelper.Minutes()},
				{Temperature: 24, TimeInMinutes: growhelper.Minutes() - 4*60},
				{Temperature: 25, TimeInMinutes: growhelper.Minutes() - (5 * 60)},
				{Temperature: 35, TimeInMinutes: growhelper.Minutes() - (5 * 60) + 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Growtable{
				mu:           tt.fields.mu,
				Measurements: tt.fields.Measurements,
			}
			if got := g.getMeasurementsInTime(tt.args.sinceHours); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Growtable.getMeasurementsInTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrowtable_UpdateColorDeviations(t *testing.T) {
	type fields struct {
		mu           sync.Mutex
		Measurements []Measurement
	}
	type args struct {
		weeks []Week
	}
	tests := []struct {
		fields fields
		args   args
		want   []Measurement
	}{
		{
			fields: fields{
				Measurements: []Measurement{
					{Temperature: -20, Humidity: -20, TimeInMinutes: growhelper.MinutesFromTime(growhelper.HtmlDateAsTime("2023-12-12")) + 1},
				},
			},
			args: args{weeks: []Week{}},
			want: []Measurement{
				{Temperature: -20, Humidity: -20, TimeInMinutes: growhelper.MinutesFromTime(growhelper.HtmlDateAsTime("2023-12-12")) + 1},
			},
		},
		{
			fields: fields{Measurements: []Measurement{}},
			args: args{weeks: []Week{
				{Temperature: 20, Humidity: 50, Start: growhelper.HtmlDateAsTime("2023-11-25")},
			},
			},
			want: []Measurement{},
		},
		{
			want: []Measurement{
				{TemperatureDeviation: 0, HumidityDeviation: 50, Temperature: 30, Humidity: 60, TimeInMinutes: growhelper.MinutesFromTime(growhelper.HtmlDateAsTime("2023-12-12")) + 1},  //40-70 ->temp0, hum50
				{TemperatureDeviation: 80, HumidityDeviation: 90, Temperature: 28, Humidity: 58, TimeInMinutes: growhelper.MinutesFromTime(growhelper.HtmlDateAsTime("2023-11-27")) + 1}, //30-60 ->temp80, hum90 (aktuell: -20/40)
				{TemperatureDeviation: 40, HumidityDeviation: 70, Temperature: 26, Humidity: 56, TimeInMinutes: growhelper.MinutesFromTime(growhelper.HtmlDateAsTime("2023-11-25")) + 1}, //20-50 ->temp40, hum70 (aktuell: 60/80)
				{TemperatureDeviation: 80, HumidityDeviation: 90, Temperature: 22, Humidity: 52, TimeInMinutes: growhelper.MinutesFromTime(growhelper.HtmlDateAsTime("2023-11-25"))},     //20-50 ->temp80, hum90 (aktuell: 20/60)

				{TemperatureDeviation: 200, HumidityDeviation: 200, Temperature: 24, Humidity: 54, TimeInMinutes: growhelper.MinutesFromTime(growhelper.HtmlDateAsTime("2023-11-25")) - 1}, //nicht dabei, passt
				{TemperatureDeviation: 200, HumidityDeviation: 200, Temperature: 20, Humidity: 50, TimeInMinutes: growhelper.MinutesFromTime(growhelper.HtmlDateAsTime("2023-11-23"))},
			},
			fields: fields{
				Measurements: []Measurement{
					{Temperature: 30, Humidity: 60, TimeInMinutes: growhelper.MinutesFromTime(growhelper.HtmlDateAsTime("2023-12-12")) + 1}, //40-70 ->temp0, hum50
					{Temperature: 28, Humidity: 58, TimeInMinutes: growhelper.MinutesFromTime(growhelper.HtmlDateAsTime("2023-11-27")) + 1}, //30-60 ->temp80, hum90
					{Temperature: 26, Humidity: 56, TimeInMinutes: growhelper.MinutesFromTime(growhelper.HtmlDateAsTime("2023-11-25")) + 1}, //20-50 ->temp40, hum70
					{Temperature: 22, Humidity: 52, TimeInMinutes: growhelper.MinutesFromTime(growhelper.HtmlDateAsTime("2023-11-25"))},     //20-50 ->temp80, hum90
					{Temperature: 24, Humidity: 54, TimeInMinutes: growhelper.MinutesFromTime(growhelper.HtmlDateAsTime("2023-11-25")) - 1}, //nicht dabei, passt
					{Temperature: 20, Humidity: 50, TimeInMinutes: growhelper.MinutesFromTime(growhelper.HtmlDateAsTime("2023-11-23"))},     //nicht dabei, passt
				},
			},
			args: args{weeks: []Week{
				{Temperature: 20, Humidity: 50, Start: growhelper.HtmlDateAsTime("2023-11-25")},
				{Temperature: 30, Humidity: 60, Start: growhelper.HtmlDateAsTime("2023-11-26")},
				{Temperature: 40, Humidity: 70, Start: growhelper.HtmlDateAsTime("2023-12-05")},
			},
			},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			g := &Growtable{
				mu:           tt.fields.mu,
				Measurements: tt.fields.Measurements,
			}
			g.UpdateAllColorDeviations(tt.args.weeks)
			if !reflect.DeepEqual(tt.want, g.Measurements) {
				t.Errorf("g.UpdateColorDeviations() = %v, want %v", g.Measurements, tt.want)
			}
		})
	}
}

func Test_getHslValue(t *testing.T) {
	type args struct {
		value   int
		optimal int
		maxDiff int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{value: 40, optimal: 65, maxDiff: 25},
			want: 0,
		},
		{
			args: args{value: 26, optimal: 20, maxDiff: 10},
			want: 40,
		},
		{
			args: args{value: 50, optimal: 50, maxDiff: 20},
			want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getHslValue(tt.args.value, tt.args.optimal, tt.args.maxDiff); got != tt.want {
				t.Errorf("getHslValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeInDays(t *testing.T) {
	type args struct {
		m Measurement
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{Measurement{TimeInMinutes: 1*24*60 - 1}},
			want: 2,
		},
		{
			args: args{Measurement{TimeInMinutes: 10 * 24 * 60}},
			want: 11,
		},
		{
			args: args{Measurement{TimeInMinutes: 12*24*60 + 300}},
			want: 13,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.m.TimeInDays(); got != tt.want {
				t.Errorf("TimeInDays() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeInString(t *testing.T) {
	type args struct {
		m Measurement
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{Measurement{TimeInMinutes: 10*24*60 + 300}},
			want: "06:00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.m.TimeInString(); got != tt.want {
				t.Errorf("TimeInString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateInString(t *testing.T) {
	type args struct {
		m Measurement
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{Measurement{TimeInMinutes: 10*24*60 + 300}},
			want: "Sun 11 Jan 1970",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.m.DateInString(); got != tt.want {
				t.Errorf("DateInString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrowtable_hasPicture(t *testing.T) {
	type fields struct {
		mu           sync.Mutex
		Measurements []Measurement
	}
	type args struct {
		minutes int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			args:   args{minutes: 10*24*60 + 300},
			fields: fields{Measurements: []Measurement{Measurement{TimeInMinutes: 10*24*60 + 300}, Measurement{TimeInMinutes: 10*24*60 + 320, Picture: true}}},
			want:   false,
		},
		{
			args:   args{minutes: 10*24*60 + 320},
			fields: fields{Measurements: []Measurement{Measurement{TimeInMinutes: 10*24*60 + 300}, Measurement{TimeInMinutes: 10*24*60 + 320, Picture: true}}},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Growtable{
				mu:           tt.fields.mu,
				Measurements: tt.fields.Measurements,
			}
			if got := g.hasPicture(tt.args.minutes); got != tt.want {
				t.Errorf("Growtable.hasPicture() = %v, want %v", got, tt.want)
			}
		})
	}
}
