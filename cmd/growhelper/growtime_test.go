package growhelper

import (
	"math"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestHtmlDate(t *testing.T) {
	type args struct {
		date string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name:    "foo",
			args:    args{date: "2006-01-02"},
			want:    time.Date(2006, time.January, 2, 0, 0, 0, 0, time.Local),
			wantErr: false,
		},
		{
			name:    "bla",
			args:    args{date: "2005/01/02"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HtmlDate(tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("HtmlDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HtmlDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDay(t *testing.T) {
	type args struct {
		minutes int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{minutes: 28611240+(24*60)},
			want: 27,
		},
		{
			name: "1-1",
			args: args{minutes: 28611240+(23*60)},
			want: 26,
		},
		{
			name: "1-2",
			args: args{minutes: 28611240},
			want: 26,
		},
		{
			name: "1-3",
			args: args{minutes: 28611240-1},
			want: 25,
		},
		{
			name: "2",
			args: args{minutes: 48 * 60},
			want: 3,
		},
		{
			name: "0",
			args: args{minutes: 0},
			want: 1,
		},
		{
			name: "5-0",
			args: args{minutes: 5},
			want: 1,
		},
		{
			name: "25-1",
			args: args{minutes: 25 * 60},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Day(tt.args.minutes); got != tt.want {
				t.Errorf("Day() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTime(t *testing.T) {
	type args struct {
		minutes int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{minutes: 5},
			want: "01:05",
		},
		{
			name: "2",
			args: args{minutes: 0},
			want: "01:00",
		},
		{
			name: "3",
			args: args{minutes: 65},
			want: "02:05",
		},
		{
			name: "4",
			args: args{minutes: 24 * 60},
			want: "01:00",
		},
		{
			name: "5",
			args: args{minutes: 25 * 60},
			want: "02:00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Time(tt.args.minutes); got != tt.want {
				t.Errorf("Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate(t *testing.T) {
	type args struct {
		minutes int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{minutes: 0},
			want: "Thu 1 Jan 1970",
		},
		{
			args: args{minutes: 5},
			want: "Thu 1 Jan 1970",
		},
		{
			args: args{minutes: 25 * 60},
			want: "Fri 2 Jan 1970",
		},
		{
			args: args{minutes: 24 * 60 * 2},
			want: "Sat 3 Jan 1970",
		},
		{
			args: args{minutes: 24 * 60 * 367},
			want: "Sun 3 Jan 1971",
		},
		{
			args: args{minutes: 24 * 60 * 368},
			want: "Mon 4 Jan 1971",
		},
		{
			args: args{minutes: 24 * 60 * 369},
			want: "Tue 5 Jan 1971",
		},
		{
			args: args{minutes: 24 * 60 * 370},
			want: "Wed 6 Jan 1971",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Date(tt.args.minutes); got != tt.want {
				t.Errorf("Date() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDaysToMinutes(t *testing.T) {
	type args struct {
		days int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{days: 1},
			want: 24 * 60,
		},
		{
			args: args{days: 0},
			want: 0,
		},
		{
			args: args{days: 10},
			want: 24 * 60 * 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DaysToMinutes(tt.args.days); got != tt.want {
				t.Errorf("DaysToMinutes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinutesNoTime(t *testing.T) {
	now := time.Now().Local()
	tests := []struct {
		name string
		want int
	}{
		{
			want: (int(now.Unix()) / 60) - (now.Hour() * 60) - (now.Minute()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MinutesNoTime(); got != tt.want {
				t.Errorf("MinutesNoTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinutes(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{want: ((time.Now().Year() - 1970) * 365 * 24 * 60) + (int(time.Now().Month() * 30 * 24 * 60)) + ((time.Now().Day() - 18) * 24 * 60) + (time.Now().Hour() * 60) + (time.Now().Minute())},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Minutes(); math.Abs(float64(got-tt.want)) > (3 * 24 * 60) {
				t.Errorf("Minutes() = %v, want %v, diff %v", got, tt.want, math.Abs(float64(got-tt.want)))
			}
		})
	}
}

func TestMinutesMidnight(t *testing.T) {
	if got := MinutesMidnight(); got < (1700154595 / 60) {
		t.Errorf("MinutesMidnight()")
	}
}

func TestAddTodayParsedTime(t *testing.T) {
	type args struct {
		timeFormat24 string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{timeFormat24: "00:05"},
			want: MinutesNoTime() + 5,
		},
		{
			args: args{timeFormat24: "02:05"},
			want: MinutesNoTime() + 125,
		},
		{
			args: args{timeFormat24: "15:00"},
			want: MinutesNoTime() + 15*60,
		},
		{
			args: args{timeFormat24: "00:00"},
			want: MinutesNoTime(),
		},
		{
			args: args{timeFormat24: "23:59"},
			want: MinutesNoTime() - 1 + 24*60,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddTodayParsedTime(tt.args.timeFormat24); got != tt.want {
				t.Errorf("AddTodayParsedTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToday(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{want: ToString(time.Now().Year()) + "-" + ToString(int(time.Now().Month())) + "-" + ToString(time.Now().Day())},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Today(); strings.Split(got, "-")[0] != strings.Split(tt.want, "-")[0] {
				t.Errorf("Today() = %v, want %v", got, tt.want)
			}
		})
	}
}
