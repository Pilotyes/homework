package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServerConfig(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Config is not nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewServerConfig()
			assert.NotNil(t, got)
		})
	}
}
