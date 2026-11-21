package config

type Config struct {
	Version int   `toml:"version"`
	Jobs    []Job `toml:"jobs"`
}

type Job struct {
	ID       string   `toml:"id"`
	Name     string   `toml:"name"`
	Script   string   `toml:"script"`
	Triggers Triggers `toml:"triggers"`
}

type Triggers struct {
	Login    bool           `toml:"login"`    // RunAtLoad
	Calendar []CalendarTime `toml:"calendar"` // StartCalendarInterval
	Interval int            `toml:"interval"` // StartInterval, seconds
	Watch    []string       `toml:"watch"`    // WatchPaths
	Queue    []string       `toml:"queue"`    // QueueDirectories
	OnMount  bool           `toml:"on_mount"` // StartOnMount
}

type CalendarTime struct {
	Hour   int `toml:"hour"`
	Minute int `toml:"minute"`
}

func validate() {
	// todo: write validation for structs
}
