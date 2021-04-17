package mapstorage

import (
	"reflect"
	"testing"
)

func TestStorage_Items(t *testing.T) {
	storage := New()
	storage.Items()
	storage.Items()
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Storage
	}{
		{
			name: "New storage",
			want: &Storage{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
