package v6

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOperatingSystemAlias_VersionMutualExclusivity(t *testing.T) {
	db := setupTestStore(t).db

	msg := "cannot have both version and version_pattern set"

	tests := []struct {
		name   string
		input  *OperatingSystemSpecifierOverride
		errMsg string
	}{
		{
			name: "version and version_pattern are mutually exclusive",
			input: &OperatingSystemSpecifierOverride{
				Alias:          "ubuntu",
				Version:        "20.04",
				VersionPattern: "20.*",
			},
			errMsg: msg,
		},
		{
			name: "only version is set",
			input: &OperatingSystemSpecifierOverride{
				Alias:   "ubuntu",
				Version: "20.04",
			},
			errMsg: "",
		},
		{
			name: "only version_pattern is set",
			input: &OperatingSystemSpecifierOverride{
				Alias:          "ubuntu",
				VersionPattern: "20.*",
			},
			errMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.Create(tt.input).Error
			if tt.errMsg == "" {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			}
		})
	}
}

func TestOperatingSystem_VersionNumber(t *testing.T) {
	tests := []struct {
		name           string
		os             *OperatingSystem
		expectedResult string
	}{
		{
			name:           "nil OS",
			os:             nil,
			expectedResult: "",
		},
		{
			name:           "major and minor versions",
			os:             &OperatingSystem{MajorVersion: "10", MinorVersion: "1"},
			expectedResult: "10.1",
		},
		{
			name:           "major version only",
			os:             &OperatingSystem{MajorVersion: "10"},
			expectedResult: "10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedResult, tt.os.VersionNumber())
		})
	}
}

func TestOperatingSystem_Version(t *testing.T) {
	tests := []struct {
		name           string
		os             *OperatingSystem
		expectedResult string
	}{
		{
			name:           "nil OS",
			os:             nil,
			expectedResult: "",
		},
		{
			name:           "label version",
			os:             &OperatingSystem{LabelVersion: "unstable"},
			expectedResult: "unstable",
		},
		{
			name:           "major and minor versions",
			os:             &OperatingSystem{MajorVersion: "10", MinorVersion: "1"},
			expectedResult: "10.1",
		},
		{
			name:           "major version only",
			os:             &OperatingSystem{MajorVersion: "10"},
			expectedResult: "10",
		},
		{
			name:           "codename",
			os:             &OperatingSystem{Codename: "buster"},
			expectedResult: "buster",
		},
		{
			name:           "with channel",
			os:             &OperatingSystem{MajorVersion: "10", MinorVersion: "1", Channel: "stable"},
			expectedResult: "10.1+stable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedResult, tt.os.Version())
		})
	}
}

func TestOperatingSystem_clean(t *testing.T) {

	tests := []struct {
		name  string
		input OperatingSystem
		want  OperatingSystem
	}{
		{
			name: "trim 0s",
			input: OperatingSystem{
				Name:         "Ubuntu",
				MajorVersion: "20",
				MinorVersion: "04",
			},
			want: OperatingSystem{
				Name:         "Ubuntu",
				MajorVersion: "20",
				MinorVersion: "4",
			},
		},
		{
			name: "preserve 0 value",
			input: OperatingSystem{
				Name:         "Redhat",
				MajorVersion: "9",
				MinorVersion: "0",
			},
			want: OperatingSystem{
				Name:         "Redhat",
				MajorVersion: "9",
				MinorVersion: "0", // important! ...9 != 9.0 since 9 includes multiple minor versions
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := tt.input
			o.clean()
			if d := cmp.Diff(tt.want, o); d != "" {
				t.Errorf("OperatingSystem.clean() mismatch (-want +got):\n%s", d)
			}
		})
	}
}
