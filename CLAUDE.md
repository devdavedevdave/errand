# errand

## What errand is
errand is a free, open-source tool for macOS, written in Go. It sits as a friendly layer on top of launchd.
The idea: you hand it a script and a trigger ("every 15 minutes", "weekdays at 9:00", "at login",
"when ~/Downloads changes"), and it takes care of the launchd plumbing. It's "your guy" that does the
little chores for you: cleaning the Desktop, emptying Downloads, running backups, pulling repos.

Writing launchd plists by hand is painful: verbose XML, `launchctl bootstrap/bootout` syntax, a minimal
PATH that makes scripts fail silently, and logs nobody can find. errand hides all of that.

## Vision (direction, not a fixed spec)
The exact scope is still open. Treat this as where it's heading, not as requirements:
- **Go library + CLI + TUI.** A core library that wraps launchd, a terminal CLI (`errand add`, `errand ls`,
  `errand run`, …) and a TUI built with Bubble Tea + Lip Gloss.
- **TOML is the source of truth.** Jobs live in a TOML config. Plists are generated build output
  (`errand apply`/`sync`), so a setup is portable across machines and diffable.
- **Human-readable triggers** that map to launchd keys (StartInterval, StartCalendarInterval, WatchPaths, RunAtLoad, …).
- **Ideas for later, not decided:** built-in logs (`errand logs -f`), run history with ✓/✗, run-now and
  `--dry-run`, failure notifications, `errand doctor` (PATH, permissions, shebangs), importing existing
  LaunchAgents, pause/resume, job templates, waking the Mac for jobs via `pmset` (needs root, unreliable on laptops).

## Architecture so far
- `cmd/errand/main.go` – entry point
- `internal/job` – core domain model (Job, Trigger). Everything depends on this, it depends on nothing.
- `internal/config` – Load/Save the TOML config (go-toml/v2). Paths are passed in, not hardcoded (testability, `--config` flag).
- `internal/launchd` – plist generation and launchctl calls
- `internal/apply` – sync TOML → plists
- `internal/tui` – Bubble Tea UI
- `Trigger` uses dedicated typed fields per trigger kind (Login, Calendar, Interval, Watch, …),
  because zero values are safe for some fields (Interval) but dangerous for others (Hour = 0 means midnight).

## Who I am and how to help me
I'm learning Go with this project. I have a Java/Kotlin background and I'm new to Go.
The goal is to understand Go, not just to have working code.
- Don't write or refactor my code. Give hints, ask questions, point me to docs (pkg.go.dev, `go doc`).
- Explain Go concepts by comparing them to Java where it helps (pointers vs. references, struct tags vs. annotations, …).
- When reviewing code: list the problems and why they matter, then let me fix them myself.
- Short snippets are fine to illustrate a concept, but not as a ready-made solution for errand itself.
- Prefer a Socratic, step-by-step style: check that I've understood before moving on.
- When there's a design decision, show me the trade-offs and let me decide.
- Point out idiomatic Go (naming, error handling, package layout) when I write Java-style Go.
