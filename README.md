# Arkeo LCD Proxy

Small Go reverse proxy that rewrites legacy LCD tx queries:
`/cosmos/tx/v1beta1/txs?events=...` -> `/cosmos/tx/v1beta1/txs?query=...`

## Requirements

- Docker + Docker Compose

## Configure

Copy the example env file and adjust as needed:

```bash
cp .env.example .env
```

Environment variables:

- `BACKEND_LCD_URL` (default: `http://127.0.0.1:1317`)
- `BACKEND_RPC_URL` (not used by the proxy; kept for reference)
- `LISTEN` (default: `:1318`)
- `LOG_FILE` (default: `~/.lcd-proxy/lcd-proxy.log`)

Note: `LOG_FILE` expands `~` at runtime. In a container, `~` resolves to the
container user's home (by default `/root`).

## Run locally (Docker Compose)

```bash
docker compose up --build -d
```

Test:

```bash
ADDR=arkeo1w2ln0prejgrztmf9w23e0rsnlks7djneh5te7p
curl -s "http://127.0.0.1:1318/cosmos/tx/v1beta1/txs?events=transfer.sender%3D%27${ADDR}%27&limit=1"
```

Logs:

- Docker: `docker compose logs -f`
- File: `LOG_FILE` (defaults to `~/.lcd-proxy/lcd-proxy.log`)

## Docker image (GitHub Container Registry)

The GitHub Actions workflow publishes to GHCR on every push to `main` and on
`v*` tags.

Image:

- `ghcr.io/arkeonetwork/arkeo-lcd-proxy:latest`

Example:

```bash
docker run -e BACKEND_LCD_URL=https://rest-seed.arkeo.network \
  -p 1318:1318 ghcr.io/arkeonetwork/arkeo-lcd-proxy:latest
```

## Server install (public image)

```bash
docker pull ghcr.io/arkeonetwork/arkeo-lcd-proxy:latest
docker run -d --name lcd-proxy --restart unless-stopped \
  -e BACKEND_LCD_URL=https://rest-seed.arkeo.network \
  -e LISTEN=:1318 \
  -e LOG_FILE=/root/.lcd-proxy/lcd-proxy.log \
  -p 1318:1318 ghcr.io/arkeonetwork/arkeo-lcd-proxy:latest
```
