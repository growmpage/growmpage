package main

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/growmpage/growganizer"
	"github.com/growmpage/growhelper"
	"github.com/growmpage/growsensoric"
	"github.com/growmpage/growtable"
)

func TestGrowmpage_GrowganizerToGrowtableWeeks(t *testing.T) {





	type fields struct {
		growcache    *bytes.Buffer
		growganizer  growganizer.Growganizer
		growtable    *growtable.Growtable
		growobserver *growsensoric.Growobserver
		growswitcher *growsensoric.Growswitcher
	}
	tests := []struct {
		name   string
		fields fields
		want   []growtable.Week
	}{
		{
			fields: fields{growganizer: growganizer.Growganizer{Weeks: []growganizer.Week{growganizer.Week{Start: "2014-03-20", Name: "ganizer", Temperature: 5}}}},
			want:   []growtable.Week{growtable.Week{Start: growhelper.HtmlDateAsTime("2014-03-20"), Temperature: 5}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Growmpage{
				growcache:    tt.fields.growcache,
				growganizer:  tt.fields.growganizer,
				growtable:    tt.fields.growtable,
				growobserver: tt.fields.growobserver,
				growswitcher: tt.fields.growswitcher,
			}
			if got := g.growganizerToGrowtableWeeks(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Growmpage.GrowganizerToGrowtableWeeks() = %v, want %v", got, tt.want)
			}
		})
	}
}
