package config

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// -------------------------------
// Load
// -------------------------------

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

// -------------------------------
// Save
// -------------------------------

func TestSaveConfigValidToml(t *testing.T) {
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

	path := filepath.Join(os.TempDir(), "validToml.xml")

	err := Save(path, &want)

	if err != nil {
		t.Fatalf("Load(%q) returned no error, want one", path)
	}
}

func TestSaveConfigInvalidToml(t *testing.T) {
	// One temp dir for all cases. Each case builds its own path inside it.
	dir := t.TempDir()
	validConfig := &Config{Version: 1}

	tests := []struct {
		name   string
		path   string
		config *Config
	}{
		// toml.Marshal can't encode a nil pointer.
		{name: "nil config", path: filepath.Join(dir, "config.toml"), config: nil},
		// os.WriteFile can't write a file where a directory already is.
		{name: "path is a directory", path: dir, config: validConfig},
		// os.WriteFile creates the file, but not missing parent directories.
		{name: "parent directory missing", path: filepath.Join(dir, "missing", "config.toml"), config: validConfig},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Save(tt.path, tt.config)

			if err == nil {
				t.Errorf("Save(%q) returned no error, want one", tt.path)
			}
		})
	}
}

// -------------------------------
// InitConfigIfNeeded
// -------------------------------

// -------------------------------
// SaveInitial
// -------------------------------

// -------------------------------
// Config Exists
// -------------------------------

func TestConfigExistsReturnsTrue(t *testing.T) {
	// One temp dir for all cases. Each case builds its own path inside it.
	dir := t.TempDir()
	path := filepath.Join(dir, configPath())

	err := os.MkdirAll(filepath.Dir(path), dirPerm)

	if err != nil {
		t.Fatalf("MkdirAll(%q): %v", filepath.Dir(path), err)
	}

	err = SaveInitial(path)

	if err != nil {
		t.Fatalf("SaveInitial(%q): %v", path, err)
	}

	exists, err := configExists(dir)

	if err != nil {
		t.Errorf("configExists(%q) returned error: %v", dir, err)
	}

	if !exists {
		t.Errorf("configExists(%q) returned false, want true", dir)
	}
}

func TestConfigDoesntExistsReturnsFalse(t *testing.T) {
	dir := t.TempDir()

	exists, err := configExists(dir)

	if err != nil {
		t.Errorf("configExists(%q) returned error: %v, want nil", dir, err)
	}

	if exists {
		t.Errorf("configExists(%q) returned true, want false", dir)
	}
}

func TestConfigExistsStatFails(t *testing.T) {
	dir := t.TempDir()

	// Put a regular file where the ".config" directory should be. os.Stat on
	// ".config/errand/config.toml" then fails with "not a directory" (ENOTDIR),
	// which is not fs.ErrNotExist, so configExists takes its error branch.
	blocker := filepath.Join(dir, ".config")
	err := os.WriteFile(blocker, nil, filePerm)

	if err != nil {
		t.Fatalf("WriteFile(%q): %v", blocker, err)
	}

	exists, err := configExists(dir)

	if err == nil {
		t.Errorf("configExists(%q) returned no error, want one", dir)
	}

	if exists {
		t.Errorf("configExists(%q) returned true, want false", dir)
	}
}

// -------------------------------
// Create New Config
// -------------------------------

func TestCreateNewConfigSuccess(t *testing.T) {
	dir := t.TempDir()

	err := createNewConfig(dir)

	if err != nil {
		t.Fatalf("createNewConfig(%q) returned error: %v", dir, err)
	}

	// Checking that the file exists isn't enough: it must also contain the
	// initial config, so load it back and compare.
	path := filepath.Join(dir, configPath())
	got, err := Load(path)

	if err != nil {
		t.Fatalf("Load(%q) after createNewConfig: %v", path, err)
	}

	// Compare fields instead of reflect.DeepEqual: the round trip turns a nil
	// Jobs slice into an empty, non-nil one, and DeepEqual treats those as different.
	if got.Version != 1 {
		t.Errorf("Load(%q) Version = %d, want 1", path, got.Version)
	}

	if len(got.Jobs) != 0 {
		t.Errorf("Load(%q) Jobs = %+v, want none", path, got.Jobs)
	}
}

func TestCreateNewConfigFailure(t *testing.T) {
	tests := []struct {
		name string
		// setup puts something in dir that makes createNewConfig fail.
		setup func(t *testing.T, dir string)
	}{
		{
			// A regular file where the ".config" directory should be,
			// so os.MkdirAll can't create ".config/errand".
			name: "mkdir fails",
			setup: func(t *testing.T, dir string) {
				blocker := filepath.Join(dir, ".config")
				if err := os.WriteFile(blocker, nil, filePerm); err != nil {
					t.Fatalf("WriteFile(%q): %v", blocker, err)
				}
			},
		},
		{
			// A directory where config.toml should be, so MkdirAll
			// succeeds but SaveInitial can't write the file.
			name: "save fails",
			setup: func(t *testing.T, dir string) {
				blocker := filepath.Join(dir, configPath())
				if err := os.MkdirAll(blocker, dirPerm); err != nil {
					t.Fatalf("MkdirAll(%q): %v", blocker, err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			tt.setup(t, dir)

			err := createNewConfig(dir)

			if err == nil {
				t.Errorf("createNewConfig(%q) returned no error, want one", dir)
			}
		})
	}
}
