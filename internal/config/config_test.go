package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func expected() *Config {
	var config Config

	config.Server.HTTP.Addr = "127.0.0.1:389"
	config.Logger.Level = "debug"

	return &config
}

func TestNewConfig(t *testing.T) {
	t.Parallel()

	config := NewConfig()
	require.NotNil(t, config)
	require.True(t, reflect.DeepEqual(&Config{}, config))
}

func TestNewConfigFromFile(t *testing.T) {
	t.Parallel()

	type args struct {
		filename string
	}

	tests := []struct {
		args    args
		want    *Config
		name    string
		wantErr bool
	}{
		{
			name: "yaml",
			args: args{
				filename: "testdata/config.yml",
			},
			want:    expected(),
			wantErr: false,
		},
		{
			name: "json",
			args: args{
				filename: "testdata/config.json",
			},
			want:    expected(),
			wantErr: false,
		},
		{
			name: "no file",
			args: args{
				filename: "testdata/dummy.json",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "unknown format",
			args: args{
				filename: "testdata/config.dummy",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "wrong file",
			args: args{
				filename: "testdata/wrong.yaml",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewConfigFromFile(tt.args.filename)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.want, got)
		})
	}
}
