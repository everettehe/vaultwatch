# vaultwatch

A CLI tool that monitors HashiCorp Vault secret expiration and sends configurable alerts before rotation deadlines.

---

## Installation

```bash
go install github.com/yourusername/vaultwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/vaultwatch.git
cd vaultwatch && go build -o vaultwatch .
```

---

## Usage

Set your Vault address and token, then run:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.yourtoken"

vaultwatch --path secret/myapp --warn-before 72h --alert slack
```

**Common flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--path` | Vault secret path to monitor | *(required)* |
| `--warn-before` | How far in advance to alert | `48h` |
| `--alert` | Alert channel (`slack`, `email`, `pagerduty`) | `stdout` |
| `--interval` | How often to check for expirations | `1h` |

**Example: run as a background watcher**

```bash
vaultwatch watch --config ./vaultwatch.yaml
```

A sample config file (`vaultwatch.yaml`) can be generated with:

```bash
vaultwatch init
```

---

## Requirements

- Go 1.21+
- HashiCorp Vault 1.10+

---

## License

MIT © 2024 [yourusername](https://github.com/yourusername)