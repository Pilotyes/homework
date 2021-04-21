package mapstorage

import (
	"reflect"
	"shop-api/internal/config"
	"testing"
)

func TestStorage_Items(t *testing.T) {
	config := config.NewConfig()
	storage := New(config)
	storage.Items()
	storage.Items()
}

func TestNew(t *testing.T) {
	config := config.NewConfig()

	tests := []struct {
		name string
		want *Storage
	}{
		{
			name: "New storage",
			want: &Storage{
				config: config,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
