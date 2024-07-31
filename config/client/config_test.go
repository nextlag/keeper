package client_test

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"

	config "github.com/nextlag/keeper/config/client"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name           string
		expectedConfig *config.Config
	}{
		{
			expectedConfig: &config.Config{
				App: &config.App{
					Name: "keeper",
				},
				Server: &config.Server{
					ServerURL: "http://localhost:8080",
				},
				Log: &config.Log{
					Level: slog.Level(-4),
				},
				SQLite: &config.SQLite{
					DSN: "keeper.sqlite",
				},
				FilesStorage: &config.FilesStorage{
					ServerLocation: "data",
					ClientLocation: "tmp",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg := config.Load()

			if tt.expectedConfig == nil {
				require.Nil(t, cfg)
			} else {
				require.NotNil(t, cfg)
				require.Equal(t, tt.expectedConfig.App.Name, cfg.App.Name)
				require.Equal(t, tt.expectedConfig.Server.ServerURL, cfg.Server.ServerURL)
				require.Equal(t, tt.expectedConfig.Log.Level, cfg.Log.Level)
				require.Equal(t, tt.expectedConfig.SQLite.DSN, cfg.SQLite.DSN)
				require.Equal(t, tt.expectedConfig.FilesStorage.ServerLocation, cfg.FilesStorage.ServerLocation)
				require.Equal(t, tt.expectedConfig.FilesStorage.ClientLocation, cfg.FilesStorage.ClientLocation)
			}
		})
	}
}
