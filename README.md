# axidev-corrige

A small **auto-correct helper** for all platforms, built with [Wails](https://wails.io/).

It listens to your **global keyboard input**, builds the current word as you type, and when you press **Space** it checks the word against a **French dictionary**. If the word looks wrong, it can:

- show a suggestion in a small overlay window
- optionally **auto-replace** the last typed word by simulating key presses (requires Accessibility permissions on macOS)

## What it does

- Tracks the word you are currently typing.
- On **Space**:
  - if the word is correct: nothing happens (you just typed a valid word)
  - if the word is incorrect: it requests spelling suggestions and uses the top suggestion
- Displays status in an overlay:
  - `Waiting...` (no current word)
  - `bonjour ✓` (word is correct)
  - `bonjor → bonjour` (suggested correction)

## Prerequisites

- Go 1.21+
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

## Development

Run in development mode with hot reload:

```bash
wails dev
```

You should see console logs like:

- `Listening for keyboard events... (Press Space to clear word)`
- `Loaded N French words into dictionary`

If needed, macOS will prompt for Accessibility permissions. Granting them is required for auto-replacement.

## Build

Build for production:

```bash
wails build
```

The binary will be in `build/bin/`.

## Notes / limitations

- The "end of word" trigger is currently **Space** only.
- Backspace, punctuation, and cursor navigation aren't handled as word-edit operations (so the tracked word can get out of sync if you edit mid-word).
- The replacement shortcut assumes the OS/app uses the standard **Option+Shift+Left** word-selection behavior.
