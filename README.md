# fivebin

Desktop vocabulary learner written in **Go** with a **Fyne** GUI.

It keeps a large dictionary out of RAM by streaming JSON, supports search and random practice words, and stores personal notes per word.

## Features

- Streaming dictionary load (memory-friendly for large word lists)
- Fast word search
- Random word practice
- Personal notes persisted to `user_meanings.json`
- Embedded default dictionary (`resources/data.json`)

## Requirements

- Go 1.21+
- OS packages required by [Fyne](https://docs.fyne.io/started/) on your platform (Linux: OpenGL / X11 or Wayland deps)

## Run

```bash
go mod download
go run .
```

Optional custom dictionary file:

```bash
go run . /path/to/dictionary.json
```

Build:

```bash
go build -o fivebin .
```

## Project layout

| File | Role |
| --- | --- |
| `main.go` | App bootstrap, dictionary path resolution |
| `dictionary.go` | Streaming dictionary access |
| `gui.go` | Fyne UI |
| `user_meanings.go` | Personal notes storage |
| `resources/data.json` | Embedded word list |

## Stack

- Language: Go
- UI: Fyne v2

## License

MIT
