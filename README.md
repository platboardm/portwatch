# portwatch

A lightweight CLI daemon that monitors port availability and sends alerts when services go down or come back up.

---

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

---

## Usage

Start monitoring a port with a simple command:

```bash
portwatch --host localhost --port 8080 --interval 10s
```

Monitor multiple ports using a config file:

```yaml
# portwatch.yaml
targets:
  - host: localhost
    port: 8080
    name: "API Server"
  - host: db.internal
    port: 5432
    name: "Postgres"
interval: 15s
alert: slack
```

```bash
portwatch --config portwatch.yaml
```

When a service goes down or recovers, `portwatch` logs the event and optionally sends an alert via webhook, Slack, or email.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--host` | `localhost` | Target host to monitor |
| `--port` | — | Target port number |
| `--interval` | `30s` | Check interval |
| `--config` | — | Path to config file |
| `--timeout` | `5s` | Connection timeout |

---

## License

MIT © 2024 [yourusername](https://github.com/yourusername)