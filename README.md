# portwatch

Lightweight CLI daemon that monitors open ports and alerts on unexpected changes.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

Start the daemon with a default scan interval:

```bash
portwatch start
```

Specify a custom interval and alert on any new or closed ports:

```bash
portwatch start --interval 30s --notify
```

Define a baseline of expected ports to suppress known services:

```bash
portwatch start --baseline 22,80,443
```

On any unexpected change, `portwatch` prints an alert to stdout (and optionally sends a system notification):

```
[ALERT] New port detected: 4444 (TCP)
[ALERT] Port closed unexpectedly: 8080 (TCP)
```

### Commands

| Command | Description |
|---|---|
| `start` | Begin monitoring open ports |
| `snapshot` | Print current open ports and exit |
| `diff` | Compare current state against a saved baseline |

## Configuration

`portwatch` can be configured via a YAML file at `~/.portwatch.yaml`:

```yaml
interval: 30s
baseline:
  - 22
  - 80
  - 443
notify: true
```

## License

MIT © 2024 yourusername