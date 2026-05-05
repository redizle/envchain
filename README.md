# envchain

> Manage layered `.env` files across environments with secret interpolation and diff support.

---

## Installation

```bash
go install github.com/yourname/envchain@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/envchain.git && cd envchain && go build -o envchain .
```

---

## Usage

envchain merges `.env` files in layers (base → environment → local), interpolates secrets, and lets you diff configurations across environments.

```bash
# Load and print resolved env for staging
envchain resolve --env staging

# Diff staging vs production
envchain diff staging production

# Run a command with the resolved environment
envchain run --env production -- ./myapp serve
```

**Layer resolution order:**

```
.env              ← base defaults
.env.staging      ← environment overrides
.env.staging.local ← local secrets (gitignored)
```

Secret interpolation uses `${secret:my-secret-name}` syntax and supports backends like AWS Secrets Manager and HashiCorp Vault.

```ini
DB_PASSWORD=${secret:prod/db/password}
API_KEY=${secret:prod/api-key}
```

---

## Configuration

envchain looks for an optional `envchain.yaml` in the project root to configure secret backends and layer paths.

---

## License

MIT © 2024 yourname