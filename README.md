# Configurator API

### By Skye Lane Goetz

Configurator API is a Gin-based Go service for:
- CURIE lookup against a DuckDB-backed Datassert database
- Downloading PMC TAR archives from a local filesystem mirror

## Quick Start

```bash
# from repo root
go run .
```

The API starts on `http://localhost:8550`.

## Requirements

- Go `1.25.7`
- Redis available at `localhost:6379` (used for API authentication)
- A Datassert DuckDB file
- A local directory containing PMC `.tar.xz` files

## Environment Setup (direnv recommended)

Use `direnv` so required paths and client credentials are loaded automatically in your shell.

```bash
# .envrc
export DATASSERT_PATH="/absolute/path/to/datassert.duckdb"
export PMC_TARS_PATH="/absolute/path/to/PMC/tars"

# optional: shell vars for curl examples
export API_USERNAME="your-username"
export API_KEY="your-api-key"
```

Then allow it:

```bash
direnv allow
```

## Build and Test Commands

```bash
# build binary
go build -o configurator-api .

# run service directly
go run .

# run all tests
go test ./...

# run tests with verbose output
go test -v ./...

# run a single test file
go test -v ./path/to/file_test.go

# run a specific test function in one file
go test -v ./path/to/file_test.go -run TestFunctionName

# run tests with coverage
go test -cover ./...
```

## Authentication Model

- Protected endpoints require query params: `username` and `api-key`.
- The server checks Redis key `username` and compares stored hash to the provided `api-key`.
- On auth failure, endpoints return `401`.

## Endpoints

| Method | Path | Auth | Purpose |
| --- | --- | --- | --- |
| `GET` | `/health` | No | Basic service health check |
| `GET` | `/search-for-curies` | Yes | Search CURIE candidates by term |
| `GET` | `/get-canonical-curie-info` | Yes | Resolve one canonical CURIE record |
| `GET` | `/download-from-pmc-tars` | Yes | Download a PMC TAR archive |

## Endpoint Details

### `GET /health`

Use this for basic service health checks (no authentication required).

Response:
- Status: `200 OK`
- Body:

```json
{"status": "ok"}
```

Example:

```bash
curl "http://localhost:8550/health"
```

### `GET /search-for-curies`

Required query params:
- `username` (string)
- `api-key` (string)
- `term` (string)

Behavior:
- Lowercases `term`
- Returns up to 50 matching CURIE rows

Example:

```bash
curl "http://localhost:8550/search-for-curies?username=$API_USERNAME&api-key=$API_KEY&term=brca1"
```

Success response shape:

```json
{
  "curies": [
    {
      "CURIE": "HGNC:1100",
      "PREFERRED_NAME": "BRCA1",
      "CATEGORY_NAME": "gene",
      "NCBI_TAXON_ID": 9606
    }
  ]
}
```

### `GET /get-canonical-curie-info`

Required query params:
- `username` (string)
- `api-key` (string)
- `curie` (string)

Behavior:
- Returns one resolved CURIE record
- Returns `404` if not found

Example:

```bash
curl "http://localhost:8550/get-canonical-curie-info?username=$API_USERNAME&api-key=$API_KEY&curie=HGNC:1100"
```

Success response shape:

```json
{
  "curie": {
    "CURIE": "HGNC:1100",
    "PREFERRED_NAME": "BRCA1",
    "CATEGORY_NAME": "gene",
    "NCBI_TAXON_ID": 9606
  }
}
```

### `GET /download-from-pmc-tars`

Required query params:
- `username` (string)
- `api-key` (string)
- `pmc-id` (string)

Accepted `pmc-id` input forms include:
- `PMC123456789`
- `123456789` (service adds `PMC`)
- `PMC:PMC123456789` (service strips prefix before `:`)

Behavior:
- Validates normalized PMC ID length is exactly 12 chars (`PMC` + 9 digits)
- Streams `application/octet-stream`
- Sets `Content-Disposition: attachment; filename=<PMC_ID>.tar.xz`

Example:

```bash
curl -L -o PMC123456789.tar.xz \
  "http://localhost:8550/download-from-pmc-tars?username=$API_USERNAME&api-key=$API_KEY&pmc-id=PMC123456789"
```

## Common Error Codes

- `400` missing or invalid required query parameter
- `401` missing/invalid auth credentials
- `404` record or PMC TAR not found
- `500` unexpected internal/database error
- `503` backend dependency unavailable (for example, DB open failure)

## Maintainer and Contributors

Maintainer and contributors are the same person.

- Skye Lane Goetz (`skye.lane.goetz@gmail.com`)
