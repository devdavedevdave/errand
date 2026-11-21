package config

import (
	"errors"
	"io/fs"
	"path/filepath"
	"reflect"
	"testing"
)

// Load
func TestLoadConfigValidToml(t *testing.T) {
	testPath := "testdata/valid.toml"

	got, err := Load(testPath)
	if err != nil {
		t.Fatalf("Load(%q) returned error: %v", testPath, err)
	}

	want := Config{
		Version: 1,
		Jobs: []Job{
			{
				ID:     "sort-screenshots",
				Name:   "Sort screenshots",
				Script: "~/scripts/sort-screenshots.sh",
				Triggers: Triggers{
					Login:    true,
					Interval: 3600,
					Watch:    []string{"~/Desktop"},
					Calendar: []CalendarTime{{Hour: 4, Minute: 30}},
				},
			},
		},
	}

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Load(%q)\n got: %+v\nwant: %+v", testPath, *got, want)
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	// A fresh temp dir is guaranteed to be empty, so this file can't exist.
	testPath := filepath.Join(t.TempDir(), "missing.toml")

	got, err := Load(testPath)

	if err == nil {
		t.Fatalf("Load(%q) returned no error, want one", testPath)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("Load(%q) error = %v, want it to wrap fs.ErrNotExist", testPath, err)
	}
	if got != nil {
		t.Errorf("Load(%q) config = %+v, want nil", testPath, got)
	}
}

func TestLoadConfigInvalidToml(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{name: "broken syntax", path: "testdata/broken_syntax.toml"},
		{name: "wrong type", path: "testdata/wrong_type.toml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.path)

			if err == nil {
				t.Fatalf("Load(%q) returned no error, want one", tt.path)
			}
			if got != nil {
				t.Errorf("Load(%q) config = %+v, want nil", tt.path, got)
			}
		})
	}
}

// Save

// SaveInitial

// ConfigExists

// CreateNewConfig
