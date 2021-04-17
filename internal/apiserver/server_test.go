package apiserver

import (
	"shop-api/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		config *config.ServerConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid config",
			args: args{
				config: config.NewServerConfig(),
			},
			wantErr: false,
		},
		{
			name: "Invalid loglevel in config",
			args: args{
				config: &config.ServerConfig{
					LogLevel: "invalid",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.NotNil(t, got)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}
