package growsensoric

import (
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/growmpage/growhelper"
)

func TestGrowobserver_Measure(t *testing.T) {
	type fields struct {
		mu        sync.Mutex
		observing Observing
	}
	tests := []struct {
		name   string
		fields fields
		want   Observing
	}{
		{
			fields: fields{observing: Observing{}},
			want:   Observing{Minutes: growhelper.Minutes(), Temperature: 25, Humidity: 60},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Growobserver{
				mu:        tt.fields.mu,
				observing: tt.fields.observing,
			}
			if got := g.Measure(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Growobserver.Measure() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrowobserver_Picture(t *testing.T) {
	type fields struct {
		mu        sync.Mutex
		observing Observing
	}
	tests := []struct {
		name   string
		fields fields
		want   Observing
	}{
		{
			fields: fields{observing: Observing{}},
			want:   Observing{Minutes: growhelper.Minutes(), Temperature: 0, Humidity: 0, Picture: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Growobserver{
				mu:        tt.fields.mu,
				observing: tt.fields.observing,
			}
			got := g.Picture()
			newFile := growhelper.Filename_pictures + growhelper.ToString(tt.want.Minutes) + ".png"
			os.Remove(newFile)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Growobserver.Picture() = %v, want %v", got, tt.want)
			}
			
		})
	}
}
