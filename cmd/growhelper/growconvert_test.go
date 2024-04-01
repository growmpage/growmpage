package growhelper

import (
	"encoding/json"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestReadReport(t *testing.T) {
	report := Report{Temperature: Statistic{Min: 3, Max: 5, Average: 2}, Humidity: Statistic{Min: 0}}
	requestBytes, _ := json.Marshal(report)
	reportAsString := string(requestBytes)
	type args struct {
		r string
	}
	tests := []struct {
		name string
		args args
		want Report
	}{
		{
			args: args{
				r: reportAsString,
			},
			want: report,
		},
		{
			args: args{
				r: "{\"Temperature\":{\"Min\":3,\"Max\":5,\"Average\":2},\"Humidity\":{\"Min\":0,\"Max\":0,\"Average\":0}}",
			},
			want: Report{
				Temperature: Statistic{Min: 3, Max: 5, Average: 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadReport(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadReport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadString(t *testing.T) {
	s := "Hi Ho World!. "
	reader := strings.NewReader(s)
	readCloser := io.NopCloser(reader)

	type args struct {
		r io.ReadCloser
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				r: readCloser,
			},
			want: s,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadString(tt.args.r); got != tt.want {
				t.Errorf("ReadString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt(t *testing.T) {
	type args struct {
		string string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{string: "5"},
			want: 5,
		},
		{
			args: args{string: " -14 "},
			want: -14,
		},
		{
			args: args{string: ""},
			want: 0,
		},
		{
			args: args{string: "+1"},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToInt(tt.args.string); got != tt.want {
				t.Errorf("ToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToString(t *testing.T) {
	type args struct {
		int int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{int: 5},
			want: "5",
		},
		{
			args: args{int: -5},
			want: "-5",
		},
		{
			args: args{int: 0},
			want: "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToString(tt.args.int); got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToFloat64(t *testing.T) {
	type args struct {
		string string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			args: args{string: "5"},
			want: 5,
		},
		{
			args: args{string: " -14 "},
			want: -14,
		},
		{
			args: args{string: ""},
			want: 0,
		},
		{
			args: args{string: "+1"},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToFloat64(tt.args.string); got != tt.want {
				t.Errorf("ToFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToStatistic(t *testing.T) {
	min, max, average := -2, 50, 26

	type args struct {
		measurements []int
	}
	tests := []struct {
		name string
		args args
		want Statistic
	}{
		{
			args: args{measurements: []int{10, 15, 23, 24, -2, 0, 0, 0, 1, 2, 2, 1, 50}},
			want: Statistic{Min: min, Max: max, Average: average},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToStatistic(tt.args.measurements); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToStatistic() = %v, want %v", got, tt.want)
			}
		})
	}
}
