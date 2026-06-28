# Changelog

## v0.2.0

### Highlights

- Migrate from Gin to Fiber v3 with a single-server architecture (replaces dual Gin + reverse-proxy design)
- Upgrade to Go 1.26
- Add graceful shutdown on SIGINT/SIGTERM
- Add `GET /health` endpoint for health checks
- New environment variables: `HOST`, `DOCS`, `INDEX`, `CONFIG`
- Config file accepts a single JSON object (auto-wrapped as array); invalid rules with empty `from`/`to` are skipped
- Expanded HTML analytics muting (Matomo, CNZZ, etc.); HTML noise removal runs even without a config file
- Add Windows amd64 and arm64 release builds

### Breaking Changes

- **Default static serving**: `INDEX=off` by default — root `/` shows a directory listing instead of auto-serving `index.html`. Single-site users should set `INDEX=on`.
- **Single port**: v0.1.7 used `PORT` for the public server and `PORT-1` for the internal static server. v0.2.0 uses only one port.
- **Bind address**: `HOST` env var replaces the fixed public `0.0.0.0` + internal `127.0.0.1` split.

### Upgrade Guide

- Single-site deployments: add `-e INDEX=on` (Docker) or set `INDEX=on` in docker-compose.
- Update image tag to `v0.2.0`:
  ```bash
  docker pull soulteary/docker-quick-docs:v0.2.0
  ```
