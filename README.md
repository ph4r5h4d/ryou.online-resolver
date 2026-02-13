# ryou.online — Dual-Stack IP Detection Service

A minimal Go service that returns the connecting client's IP address as plain text. Used by [ryou.onl](https://ryou.onl) to detect a user's IPv4 and IPv6 addresses separately.

## How It Works

The service relies on a DNS trick to force clients onto a specific protocol:

- `v4.ryou.online` has **only an A record** → clients must connect over IPv4
- `v6.ryou.online` has **only a AAAA record** → clients must connect over IPv6

The frontend fetches both endpoints in parallel and gets both addresses.

> DNS records must be **DNS only** (gray cloud in Cloudflare) — proxying would break the single-stack forcing.

## Project Structure

```
├── main.go              # Go HTTP service
├── Dockerfile           # Multi-stage build
├── docker-compose.yml   # Traefik + ip-service
├── traefik.yml          # Traefik static configuration
└── README.md
```

## Endpoints

| Method | Path | Response |
|--------|------|----------|
| `GET` | `/` | Client IP as plain text |
| `OPTIONS` | `/` | `204 No Content` (CORS preflight) |
| `GET` | `/health` | `ok` |
| * | * | `404` |

CORS is restricted to `https://ryou.onl`.

## DNS Setup (Cloudflare)

| Type | Name | Content | Proxy |
|------|------|---------|-------|
| A | `v4` | `<server-ipv4>` | DNS only |
| AAAA | `v6` | `<server-ipv6>` | DNS only |
| A | `ryou.online` | `<server-ipv4>` | Proxied (optional) |
| AAAA | `ryou.online` | `<server-ipv6>` | Proxied (optional) |

The `ryou.online` records are only needed if handling the redirect to `ryou.onl` via Traefik. Alternatively, set up the redirect as a Cloudflare Redirect Rule.

## Stack

- **Go** (stdlib only) — HTTP service returning client IP
- **Traefik v3.3** — Reverse proxy, automatic Let's Encrypt TLS, rate limiting
- **Docker Compose** — Orchestration

## Deployment

### Prerequisites

- A dual-stack server (IPv4 + IPv6), e.g., Hetzner
- Docker and Docker Compose installed
- DNS records configured as above

### Deploy

```bash
git clone <repo-url> && cd ryou.online-resolver
docker compose up -d
```

### Verify

```bash
curl https://v4.ryou.online
# → 185.123.45.67

curl -6 https://v6.ryou.online
# → 2a01:4f8:c17:abcd::1

curl https://v4.ryou.online/health
# → ok
```

## Configuration

### Rate Limiting

Configured in `docker-compose.yml` via Traefik labels:

- **Average:** 30 requests/second
- **Burst:** 50 requests

Adjust the `ratelimit` middleware labels to change these values.

### TLS Certificates

Traefik automatically obtains and renews Let's Encrypt certificates. Update the ACME email in `traefik.yml` if needed:

```yaml
certificatesResolvers:
  letsencrypt:
    acme:
      email: admin@ryou.online
```

### Redirect

`ryou.online` → `ryou.onl` is handled as a permanent (301) redirect by Traefik.

## Frontend Integration

```js
const [v4, v6] = await Promise.allSettled([
  fetch('https://v4.ryou.online').then(r => r.text()),
  fetch('https://v6.ryou.online').then(r => r.text()),
]);

const ipv4 = v4.status === 'fulfilled' ? v4.value.trim() : null;
const ipv6 = v6.status === 'fulfilled' ? v6.value.trim() : null;
```

If the user lacks one protocol, that fetch will fail gracefully and return `null`.

## License

MIT
