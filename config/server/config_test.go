package server_test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	config "github.com/nextlag/keeper/config/server"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name           string
		expectedConfig *config.Config
	}{
		{
			name: "Load default config",
			expectedConfig: &config.Config{
				Network: &config.Network{
					Host:  ":8080",
					HTTPS: false,
				},
				Log: &config.Log{
					Level:       slog.Level(-4),
					ProjectPath: "/app/",
					LogToFile:   false,
					LogPath:     "/app/logs/out.log",
				},
				Security: &config.Security{
					Domain:                "localhost",
					AccessTokenExpiresIn:  600 * time.Minute,
					RefreshTokenExpiresIn: 6000 * time.Minute,
					AccessTokenMaxAge:     600,
					RefreshTokenMaxAge:    6000,
				},
				PG: &config.PG{
					PoolMax: 2,
				},
				Cache: &config.Cache{
					DefaultExpiration: 5 * time.Minute,
					CleanupInterval:   10 * time.Minute,
				},
				FilesStorage: &config.FilesStorage{
					Location: "data",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg, err := config.Load()
			if tt.expectedConfig == nil {
				require.Nil(t, cfg)
			}
			if err != nil {
				require.Error(t, err)
			}

			require.NotNil(t, cfg)
			require.Equal(t, tt.expectedConfig.Network.Host, cfg.Network.Host)
			require.Equal(t, tt.expectedConfig.Network.HTTPS, cfg.Network.HTTPS)
			require.Equal(t, tt.expectedConfig.Log.Level, cfg.Log.Level)
			require.Equal(t, tt.expectedConfig.Log.ProjectPath, cfg.Log.ProjectPath)
			require.Equal(t, tt.expectedConfig.Log.LogToFile, cfg.Log.LogToFile)
			require.Equal(t, tt.expectedConfig.Log.LogPath, cfg.Log.LogPath)
			require.Equal(t, tt.expectedConfig.Security.Domain, cfg.Security.Domain)
			require.Equal(t, tt.expectedConfig.Security.AccessTokenExpiresIn, cfg.Security.AccessTokenExpiresIn)
			require.Equal(t, tt.expectedConfig.Security.RefreshTokenExpiresIn, cfg.Security.RefreshTokenExpiresIn)
			require.Equal(t, tt.expectedConfig.Security.AccessTokenMaxAge, cfg.Security.AccessTokenMaxAge)
			require.Equal(t, tt.expectedConfig.Security.RefreshTokenMaxAge, cfg.Security.RefreshTokenMaxAge)
			require.Equal(t, tt.expectedConfig.PG.PoolMax, cfg.PG.PoolMax)
			require.Equal(t, tt.expectedConfig.Cache.DefaultExpiration, cfg.Cache.DefaultExpiration)
			require.Equal(t, tt.expectedConfig.Cache.CleanupInterval, cfg.Cache.CleanupInterval)
			require.Equal(t, tt.expectedConfig.FilesStorage.Location, cfg.FilesStorage.Location)
		})
	}
}
