# Workshot

**Never lose your development context again.**

Workshot is a lightweight, cross-platform CLI tool that captures and restores your development environment in seconds. Built for developers who context-switch frequently and need a fast, reliable way to resume work exactly where they left off.

---

## The Problem

You're deep in flow state when suddenly:

* 💬 Urgent Slack message
* 🔥 Production incident
* 📞 Unexpected meeting

When you return 30 minutes later:

* *"Which directory was I in?"*
* *"What branch was I on?"*
* *"What commands was I running?"*

**Workshot eliminates this mental overhead with a single command.**

---

## How It Works

### Save your context

```bash
workshot freeze my-feature
```

### Resume instantly

```bash
workshot restore my-feature
```

That's it. You're back in context.

---

## Key Features

### 🎯 **Instant Context Switching**

Freeze and restore your entire development state in under a second.

### 🔍 **Git-Aware**

Captures branch, commit, remote, stash count, and dirty state automatically.

### 📝 **Command History**

Saves your recent terminal commands (with automatic secret filtering).

### 🔌 **Plugin Architecture**

Extensible design — add your own context capturers.

### 🔒 **Privacy First**

* No telemetry
* No network calls
* Local-only storage
* Automatic credential filtering

---

## 🌍 Cross-Platform Support

Workshot is designed to work on **Linux, macOS, and Windows**.

### 🧪 Tested On

* **Windows** (PowerShell)

### 🧩 Expected to Work On (not yet manually tested)

* **Linux** (Bash, Zsh)
* **macOS** (Zsh)

> ⚠️ If you encounter any platform-specific issues, please open an issue with your OS, shell, and Workshot version.

---

## Installation

### Option 1: Go Install (Recommended)

```bash
go install github.com/ansoncodes/workshot/cmd/workshot@latest
```

### Option 2: Pre-built Binaries

Download from [GitHub Releases](https://github.com/ansoncodes/workshot/releases)

Available for:

* **Linux** (amd64, arm64)
* **macOS** (Intel & Apple Silicon)
* **Windows** (amd64)

### Verify Installation

```bash
workshot --version
```

---

## Quick Start Guide

### 1. Save Your Current Context

```bash
cd ~/projects/my-app
git checkout feature/api-redesign
workshot freeze api-work
```

**Output:**

```
✓ Snapshot saved: api-work
  Location: ~/.workshot/shots/api-work.json
```

### 2. Switch to Another Task

```bash
cd ~/projects/another-project
git checkout hotfix/critical-bug
# Work on something else...
```

### 3. Return to Your Original Context

```bash
workshot restore api-work
```

**Output:**

```
Snapshot: api-work
Created: 2025-01-15 14:23:05

Working Directory:
 /home/user/projects/my-app

Git State:
 Branch:  feature/api-redesign
 Remote:  git@github.com:user/my-app.git
 Status:  Modified (2 uncommitted changes)
 Commit:  a3f2c1b
 Stashes: 1

Recent Commands:
 git checkout -b feature/api-redesign
 npm install
 npm run dev
 git status

Commands to restore:
 cd "/home/user/projects/my-app"
 git checkout feature/api-redesign
```

### 4. Execute Restore Commands

#### Linux / macOS (Bash, Zsh)

```bash
eval $(workshot restore api-work -c)
```

#### Windows (PowerShell)

```powershell
Invoke-Expression (workshot restore api-work -c | Out-String)
```

> On Windows, use `Invoke-Expression` instead of `eval`.

---

## What Gets Captured

### 📂 **Directory Context**

* Absolute path of your working directory

### 🌿 **Git Information**

* Current branch
* Remote repository URL
* Dirty state (uncommitted changes)
* Latest commit SHA
* Number of stashes

### 💻 **Terminal History**

* Last ~20 commands
* Supports Bash, Zsh, and PowerShell
* Automatically filters sensitive data (API keys, passwords, tokens)

### 📋 **Metadata**

* Snapshot creation timestamp
* Schema version (for future compatibility)

---

## All Commands

| Command                      | Description                                                                                          |
| ---------------------------- | ---------------------------------------------------------------------------------------------------- |
| `workshot freeze <name>`     | Capture the current **working directory, git context, and recent terminal commands** into a snapshot |
| `workshot restore <name>`    | Show the saved snapshot details and **print the steps required to restore the context**              |
| `workshot restore <name> -c` | **Emit shell commands** that restore the **working directory and git branch** (for `eval` / `iex`)   |
| `workshot list`              | List all saved workshot snapshots                                                                    |
| `workshot show <name>`       | Display detailed information about a snapshot (directory, git info, commands)                        |
| `workshot show <name> -j`    | Output the snapshot data as **raw JSON**                                                             |
| `workshot delete <name>`     | Permanently delete a saved snapshot                                                                  |
| `workshot --version`         | Display the installed Workshot version                                                               |


---

## Example Snapshot

```bash
workshot show api-work
```

**Output:**

```bash
Snapshot: api-work
Created: 2025-01-15 14:23:05

Working Directory:
 /home/user/projects/my-app

Git State:
 Branch:  feature/api-redesign
 Remote:  git@github.com:user/my-app.git
 Status:  Dirty (uncommitted changes)
 Commit:  a3f2c1b
 Stashes: 1

Recent Commands:
 git checkout -b feature/api-redesign
 npm install
 npm run dev

Metadata:
 Schema Version: 1
 Age:            just now
 Plugins:        2 active
```

**Storage Location:** `~/.workshot/shots/`

---

## Understanding Limitations

Workshot captures **context**, not live execution state.

### ✅ What Workshot CAN Do

* Remember which directory you were in
* Track your Git branch and commit
* Recall recent commands
* Generate restore commands

### ❌ What Workshot CANNOT Do

* Change your current shell's directory
* Restore terminal output or scrollback
* Resume running processes or servers
* Restore editor tabs or cursor positions

This is a **fundamental OS and shell limitation**, not a flaw in Workshot.

### The Solution: `eval`

```bash
eval $(workshot restore <name> -c)
```

This is the **only correct way** for a CLI tool to modify shell context.

---

## Advanced Usage

### Scripting & Automation

```bash
workshot freeze before-branch-switch && git checkout main
workshot freeze build-$BUILD_NUMBER
workshot list
```

### Alias Shortcuts

```bash
alias wsf='workshot freeze'
alias wsr='eval $(workshot restore $1 -c)'
alias wsl='workshot list'
```

---

## Technical Architecture

### Storage

* JSON snapshots in `~/.workshot/shots/`
* Indexed via `index.json`
* Atomic writes to prevent corruption

### Plugin System

```go
type Capturer interface {
    Name() string
    Priority() int
    Capture() (map[string]interface{}, error)
    Restore(data map[string]interface{}) error
    CanRestore(data map[string]interface{}) bool
}
```

**Built-in plugins:**

* `git`
* `terminal`

---

## Privacy & Security

* No telemetry
* No network calls
* No cloud sync
* Local-only storage
* Automatic secret filtering

---

## Roadmap

### v0.2.0

* VS Code open files
* Docker container capture
* Improved PowerShell support

### v0.3.0

* Browser tabs
* Tmux sessions
* Snapshot diffing

---

## License

MIT License — see [LICENSE](./LICENSE)

---

Built with ❤️ by [Anson](https://github.com/ansoncodes)

**Workshot: Save context. Switch tasks. Come back instantly.**

[⭐ Star on GitHub](https://github.com/ansoncodes/workshot)
