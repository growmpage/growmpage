package growsensoric

import (
	"sync"
	"testing"

	"github.com/growmpage/growganizer"
)

func TestGrowswitcher_Switch(t *testing.T) {
	type fields struct {
		mu sync.Mutex
	}
	type args struct {
		plug growganizer.PlugControl
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			fields: fields{},
			args:   args{plug: growganizer.PlugControl{Name: "test"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Growswitcher{
				mu: tt.fields.mu,
			}
			g.Switch(tt.args.plug)
		})
	}
}

func TestGrowswitcher_SwitchSim(t *testing.T) {
	type fields struct {
		mu sync.Mutex
	}
	type args struct {
		plugOn  growganizer.PlugControl
		plugOff growganizer.PlugControl
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			fields: fields{},
			args:   args{plugOn: growganizer.PlugControl{Name: "on"}, plugOff: growganizer.PlugControl{Name: "off"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Growswitcher{
				mu: tt.fields.mu,
			}
			g.SwitchSim(tt.args.plugOn, tt.args.plugOff)
		})
	}
}
