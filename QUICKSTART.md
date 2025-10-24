# Quick Start (Cloudflare D1)

This quick start assumes the repository is using Cloudflare D1 as the primary database. If you still need a local SQLite workflow for model generation, see `pkg/db/README.md`.

## Prerequisites

- Go 1.25+
- Bun (or npm/pnpm)
- Cloudflare account with a D1 database
- Cloudflare API token with D1 permissions
- YouTube Data API Key

## 1. Configure environment variables

Set the required env vars in your shell (or a local `.env`):

```bash
export D1_ACCOUNT_ID="<your-account-id>"
export D1_DATABASE_ID="<your-d1-database-id>"
export CLOUDFLARE_API_TOKEN="<your-api-token>"
export YOUTUBE_API_KEY="<your-youtube-key>"
export IGDB_CLIENT_ID="<your-igdb-client-id>"
export IGDB_CLIENT_SECRET="<your-igdb-client-secret>"
```

## 2. Generate Go-Jet models (if needed)

If you modify the SQL schema and need generated Go-Jet models, generate them from the schema file:

```bash
cd pkg
jet -dsn=file://db/schema.sql -path=./repository/table
```

This uses the schema file directly and does not require a local SQLite DB.

## 3. Start the backend

```bash
cd pkg
go mod download
go run main.go
```

The backend will listen on `:8088` by default.

## 4. Sync YouTube videos

With the backend running, sync videos from the OPZTV channel:

```bash
curl -X POST "http://localhost:8088/api/videos/sync?api_key=$YOUTUBE_API_KEY&max_results=50"
```

Or use the provided script (ensure `YOUTUBE_API_KEY` is exported):

```bash
./scripts/sync-youtube.sh
```

## 5. Start the frontend (development)

```bash
cd web
bun install
bun dev
```

Open `http://localhost:3000` in your browser.

## Quick commands

```bash
# Backend
cd pkg && go run main.go

# Frontend dev
cd web && bun dev

# Generate models
cd pkg && jet -dsn=file://db/schema.sql -path=./repository/table

# Build for production
cd web && bun run build
cd pkg && go build -o zeedzad
```

## Health & API checks

```bash
curl http://localhost:8088/           # health
curl http://localhost:8088/api/videos  # list videos
curl http://localhost:8088/api/databasez # db health
```

If you require a local SQLite-only workflow (for offline model generation), see `pkg/db/README.md` for instructions.
