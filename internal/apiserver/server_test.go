package apiserver

import (
	"shop-api/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		config *config.Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid config",
			args: args{
				config: config.NewConfig(),
			},
			wantErr: false,
		},
		{
			name: "Invalid loglevel in config",
			args: args{
				config: &config.Config{
					Server: &config.ServerConfig{
						LogLevel: "invalid",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "MongoDB driver not implemented yet",
			args: args{
				config: &config.Config{
					Server: config.NewServerConfig(),
					Databases: &config.DatabasesConfig{
						Driver: config.MongoDBDriver,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid database driver in config",
			args: args{
				config: &config.Config{
					Server: config.NewServerConfig(),
					Databases: &config.DatabasesConfig{
						Driver: "invalid driver",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Nil SyncConfig in MainConfig for sync job error",
			args: args{
				config: &config.Config{
					Server:    config.NewServerConfig(),
					Databases: config.NewDatabasesConfig(),
					Sync:      nil,
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
