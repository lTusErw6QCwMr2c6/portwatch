# portwatch

A lightweight CLI daemon that monitors and logs port activity changes on a host in real time.

---

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git
cd portwatch
go build -o portwatch .
```

---

## Usage

Start the daemon with default settings:

```bash
portwatch start
```

Watch specific ports and log output to a file:

```bash
portwatch start --ports 80,443,8080 --interval 5s --log /var/log/portwatch.log
```

Run a one-time snapshot of active ports:

```bash
portwatch scan
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--ports` | all | Comma-separated list of ports to monitor |
| `--interval` | `10s` | Polling interval |
| `--log` | stdout | Path to log file |
| `--verbose` | false | Enable verbose output |

---

## Example Output

```
2024/01/15 10:32:01 [OPEN]   port 8080 — PID 3821 (go run main.go)
2024/01/15 10:32:11 [CLOSED] port 8080
2024/01/15 10:33:05 [OPEN]   port 5432 — PID 1042 (postgres)
```

---

## Requirements

- Go 1.21+
- Linux or macOS (Windows support experimental)

---

## License

MIT © 2024 [yourusername](https://github.com/yourusername)