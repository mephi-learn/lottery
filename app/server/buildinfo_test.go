package main

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseBuildInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		debugbi  debug.BuildInfo
		expected BuildInfo
	}{
		{
			name: "go run",
			debugbi: debug.BuildInfo{
				Main: debug.Module{
					Path:    "whatever",
					Version: "(devel)",
				},
				Settings: nil,
			},
			expected: BuildInfo{"whatever", "(devel)"},
		},
		{
			name: "tagged",
			debugbi: debug.BuildInfo{
				Main: debug.Module{
					Path:    "whatever",
					Version: "v1.2.3",
				},
				Settings: nil,
			},
			expected: BuildInfo{"whatever", "v1.2.3"},
		},
		{
			name: "untagged",
			debugbi: debug.BuildInfo{
				Main: debug.Module{
					Path:    "whatever",
					Version: "(devel)",
				},
				Settings: []debug.BuildSetting{
					{Key: "vcs.revision", Value: "12345678901234567890"},
					{Key: "vcs.time", Value: "2024-01-01T11:12:13Z"},
					{Key: "vcs.modified", Value: "false"},
				},
			},
			expected: BuildInfo{"whatever", "v0.0.0-20240101111213-123456789012"},
		},
		{
			name: "dirty",
			debugbi: debug.BuildInfo{
				Main: debug.Module{
					Path:    "whatever",
					Version: "(devel)",
				},
				Settings: []debug.BuildSetting{
					{Key: "vcs.revision", Value: "12345678901234567890"},
					{Key: "vcs.time", Value: "2024-01-01T11:12:13Z"},
					{Key: "vcs.modified", Value: "true"},
				},
			},
			expected: BuildInfo{"whatever", "v0.0.0-20240101111213-123456789012+dirty"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			bi := parseBuildInfo(&test.debugbi)
			require.Equal(t, test.expected, bi)
		})
	}
}
