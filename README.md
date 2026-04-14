# vaultpulse

> A lightweight CLI to monitor HashiCorp Vault secret expiration and rotation schedules with alerting hooks.

---

## Installation

```bash
go install github.com/yourusername/vaultpulse@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultpulse/releases).

---

## Usage

Set your Vault address and token, then run:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.yourtoken"

vaultpulse monitor --path secret/myapp --warn-before 72h
```

**Example output:**

```
[WARN]  secret/myapp/db-password   expires in 68h (threshold: 72h)
[OK]    secret/myapp/api-key       expires in 312h
[CRIT]  secret/myapp/tls-cert      expires in 2h  (threshold: 72h)
```

### Alerting Hooks

Send alerts to a webhook (e.g., Slack, PagerDuty) when secrets approach expiration:

```bash
vaultpulse monitor \
  --path secret/myapp \
  --warn-before 72h \
  --webhook https://hooks.slack.com/services/your/webhook/url
```

### Available Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--path` | *(required)* | Vault secret path to monitor |
| `--warn-before` | `48h` | Alert threshold before expiration |
| `--webhook` | `""` | Webhook URL for alerting |
| `--interval` | `5m` | Polling interval |
| `--output` | `text` | Output format: `text` or `json` |

---

## Requirements

- Go 1.21+
- HashiCorp Vault 1.10+

---

## License

[MIT](LICENSE)