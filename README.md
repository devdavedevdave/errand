# errand

**Your guy for the little chores on your Mac.**

You give errand a script and tell it when to run. errand handles the rest.

- Clean up the Desktop every evening
- Empty Downloads once a week
- Run a backup at 4:00 every night
- Sort screenshots the moment they land

> ⚠️ **Early days.** errand is being built right now. The config format and the
> launchd plumbing exist; the `errand` command itself doesn't work yet.

## Why?

macOS already has a built-in way to run things on a schedule: **launchd**.
But using it by hand means:

- writing long, fiddly XML files (`.plist`)
- remembering `launchctl` commands
- scripts failing silently because launchd's `PATH` barely has anything in it
- logs that end up who-knows-where

errand hides all of that. You write a short, readable config file and errand
turns it into what launchd needs.

## How it works

All your jobs live in one file, `config.toml`:

```toml
version = 1

[[jobs]]
id = "nightly-backup"
name = "Nightly backup"
script = "~/scripts/backup.sh"

    [jobs.triggers]
    calendar = [{ hour = 4, minute = 0 }]   # every day at 4:00

[[jobs]]
id = "sort-screenshots"
name = "Sort screenshots"
script = "~/scripts/sort-screenshots.sh"

    [jobs.triggers]
    watch = ["~/Desktop"]                   # whenever the Desktop changes
```

That's it. One file, easy to read, easy to copy to a new Mac.

## Triggers

A trigger is *when* a job runs. A job can have more than one.

| Trigger    | Runs the job…                       | Example                          |
|------------|-------------------------------------|----------------------------------|
| `login`    | when you log in                     | `login = true`                   |
| `calendar` | at a time of day                    | `calendar = [{ hour = 9, minute = 0 }]` |
| `interval` | every N seconds                     | `interval = 900` (15 minutes)    |
| `watch`    | when a file or folder changes       | `watch = ["~/Downloads"]`        |
| `queue`    | when files appear in a folder       | `queue = ["~/Inbox"]`            |
| `on_mount` | when a disk is plugged in           | `on_mount = true`                |

## Environment variables

Scripts run by launchd get almost no `PATH`, so tools like `git` or anything
from Homebrew often "aren't found." Set what you need, for all jobs or per job:

```toml
[env]                      # for every job
LANG = "en_US.UTF-8"

[[jobs]]
id = "pull-repos"
script = "~/scripts/pull.sh"

    [jobs.env]             # for this job only
    PATH = "/opt/homebrew/bin:/usr/bin:/bin"
```

## Planned

- `errand apply`: turn `config.toml` into launchd jobs
- `errand ls`, `errand run`, `errand logs`
- A friendly terminal UI

## License

[Apache 2.0](LICENSE)
