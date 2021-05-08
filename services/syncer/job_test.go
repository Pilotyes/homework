package syncer

import (
	"shop-api/internal/config"
	"shop-api/internal/storage"
	"shop-api/internal/storage/mapstorage"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewJobService(t *testing.T) {
	testStorage := mapstorage.New(config.NewConfig())
	testLogger := logrus.New()
	testSyncConfig := config.NewSyncConfig()
	type args struct {
		name       string
		storage    storage.Storage
		logger     *logrus.Logger
		syncConfig *config.SyncConfig
	}
	tests := []struct {
		name    string
		args    args
		want    *JobService
		wantErr bool
	}{
		{
			name: "empty name",
			args: args{
				name:       "",
				storage:    testStorage,
				logger:     testLogger,
				syncConfig: testSyncConfig,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "storage is nil",
			args: args{
				name:       "job name",
				storage:    nil,
				logger:     testLogger,
				syncConfig: testSyncConfig,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "logger is nil",
			args: args{
				name:       "job name",
				storage:    testStorage,
				logger:     nil,
				syncConfig: testSyncConfig,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "syncConfig is nil",
			args: args{
				name:       "job name",
				storage:    testStorage,
				logger:     testLogger,
				syncConfig: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid args",
			args: args{
				name:       "job name",
				storage:    testStorage,
				logger:     testLogger,
				syncConfig: testSyncConfig,
			},
			want: &JobService{
				name:    "job name",
				storage: testStorage,
				config:  testSyncConfig,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewJobService(tt.args.name, tt.args.storage, tt.args.logger, tt.args.syncConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJobService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil {
				assert.Equal(t, tt.want.name, got.name)
				assert.Equal(t, tt.want.storage, got.storage)
				assert.Equal(t, tt.want.config, got.config)
				assert.NotNil(t, got.logger)
			}
		})
	}
}
