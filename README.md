# URL-StatusChecker-CLI

A simple CLI tool in Go that checks whether a URL is live.

## Requirements

- Go 1.22+

## Run

```bash
go run . <url>
```

Example:

```bash
go run . https://example.com
go run . example.com
```

## Build

```bash
go build -o url-status-checker .
./url-status-checker <url>
```

## Output

- `LIVE: <url> (HTTP <status>)` for successful responses (`2xx` and `3xx`)
- `DOWN: <url> (HTTP <status>)` for unsuccessful HTTP responses
- `DOWN: <url> (<error>)` for connection/timeout/DNS errors

## Exit Codes

- `0` if live
- `1` for invalid usage/input
- `2` if URL is down or unreachable