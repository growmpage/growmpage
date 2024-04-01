package growhelper

import (
	"testing"
)

func TestUrl(t *testing.T) {
	type args struct {
		suffix string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{suffix: "/test"},
			want: "http://127.0.0.1:99999/test",
		},
		{
			args: args{suffix: "bla?fuu"},
			want: "http://127.0.0.1:99999/bla?fuu",
		},
	}
	for _, tt := range tests {
		Port = ":99999"
		t.Run(tt.name, func(t *testing.T) {
			if got := Url(tt.args.suffix); got != tt.want {
				t.Errorf("Url() = %v, want %v", got, tt.want)
			}
		})
	}
}
