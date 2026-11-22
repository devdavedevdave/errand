// Package launchd turns errand jobs into launchd property lists (plists).
//
// plist.go: the Plist type, launchd's view of a job (Label, ProgramArguments,
// RunAtLoad, StartInterval, StartCalendarInterval, WatchPaths, ...).

package launchd

type Plist struct {
	Label            string   `plist:"Label"`
	ProgramArguments []string `plist:"ProgramArguments"`

	RunAtLoad             bool               `plist:"RunAtLoad,omitempty"`
	StartInterval         int                `plist:"StartInterval,omitempty"`
	StartCalendarInterval []CalendarInterval `plist:"StartCalendarInterval,omitempty"`
	WatchPaths            []string           `plist:"WatchPaths,omitempty"`
	QueueDirectories      []string           `plist:"QueueDirectories,omitempty"`
	StartOnMount          bool               `plist:"StartOnMount,omitempty"`
	StandardOutPath       string             `plist:"StandardOutPath,omitempty"`
	StandardErrorPath     string             `plist:"StandardErrorPath,omitempty"`
	EnvironmentVariables  map[string]string  `plist:"EnvironmentVariables,omitempty"`
}

type CalendarInterval struct {
	Minute  *int `plist:"Minute,omitempty"`  // 0-59
	Hour    *int `plist:"Hour,omitempty"`    // 0-23
	Day     *int `plist:"Day,omitempty"`     // day of month, 1-31
	Weekday *int `plist:"Weekday,omitempty"` // 0 and 7 are both Sunday
	Month   *int `plist:"Month,omitempty"`   // 1-12
}
