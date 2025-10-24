# Zeedzad - YouTube Video to Game Matcher

A full-stack web application that matches YouTube videos from the OPZTV channel with Steam games.

This repository has been migrated to use Cloudflare D1 as the primary database (serverless SQLite via Cloudflare's API).

## What this doc covers

- How to run the app locally against Cloudflare D1
- How to generate Go-Jet models from the schema file
- Backend and frontend run/build steps
- Important environment variables (D1 / API keys)

## Features

- Video listing and pagination (24 items per page)
- Match YouTube videos to Steam games using a modal search
- Sync videos from YouTube (YouTube Data API v3)
- Swagger API documentation for the backend

## Tech Stack

### Backend
- Go (Fiber v3)
- Cloudflare D1 (via github.com/cloudflare/cloudflare-go)
- Go-Jet for type-safe SQL models and queries

### Frontend
- Nuxt 4 + TypeScript
- Tailwind CSS + DaisyUI

## Prerequisites

- Go 1.25+
- Bun (or npm/pnpm) for the Nuxt app
- Cloudflare account with a D1 database
- A Cloudflare API token with D1 permissions
- YouTube Data API key (for syncing videos)
- IGDB Client ID / Secret (for IGDB lookups)

## Environment variables

Create a `.env` (or set env vars in your shell) with the following keys:

- `D1_ACCOUNT_ID` (required)
- `D1_DATABASE_ID` (required)
- `CLOUDFLARE_API_TOKEN` (required)
- `YOUTUBE_API_KEY` (required for sync)
- `IGDB_CLIENT_ID` (optional — required for IGDB searches)
- `IGDB_CLIENT_SECRET` (optional — required for IGDB searches)

There is an `.env.example` file with the template values.

## Local development (D1)

1) Ensure your D1 database exists (create it with Wrangler or in the dashboard) and you have the `D1_DATABASE_ID` and `D1_ACCOUNT_ID`.

2) Export the required environment variables (example):

```bash
export D1_ACCOUNT_ID="<your-account-id>"
export D1_DATABASE_ID="<your-d1-database-id>"
export CLOUDFLARE_API_TOKEN="<your-api-token>"
export YOUTUBE_API_KEY="<your-youtube-key>"
export IGDB_CLIENT_ID="<your-igdb-client-id>"
export IGDB_CLIENT_SECRET="<your-igdb-client-secret>"
```

3) Generate Go-Jet models (uses schema file directly):

```bash
cd pkg
# Generate types/tables from the schema file
jet -dsn=file://db/schema.sql -path=./repository/table
```

4) Start the backend:

```bash
cd pkg
go mod download
go run main.go
```

The backend listens on `:8088` by default and exposes the API under `/api/`.

5) Start the frontend (development):

```bash
cd web
bun install
bun dev
```

The frontend dev server runs at `http://localhost:3000` and will proxy API requests in development.

## API Endpoints (summary)

- `GET /` — Health check
- `GET /api/databasez` — Database health
- `GET /api/videos` — List videos (query: offset, limit, search)
- `GET /api/videos/:id` — Video by ID
- `PUT /api/videos/:id/game` — Attach a game to a video
- `POST /api/videos/sync` — Sync videos from YouTube (query: api_key, max_results)
- `GET /api/games` — List games
- `GET /api/games/steam/search` — Search Steam (query: q)

Swagger UI is available at `/api/swagger/` when running the backend in dev with docs generated.

## Production build

1. Build frontend:

```bash
cd web
bun run build
```

2. Copy built assets into the backend embed directory:

```bash
cp -r web/.output/public/* pkg/web/public/
```

3. Build backend binary:

```bash
cd pkg
go build -o zeedzad
```

4. Run production binary (env vars must be set):

```bash
export D1_ACCOUNT_ID="..."
export D1_DATABASE_ID="..."
export CLOUDFLARE_API_TOKEN="..."
./zeedzad
```

## Generating API docs

After updating handler comments, generate Swagger docs:

```bash
cd pkg
swag init -d . -g server/server.go -o docs --ot go
```

## Project layout

See the `pkg/` and `web/` folders for source code. Backend entrypoint is `pkg/main.go`.

## Security note

Do not commit files containing credentials (for example, `zerver.env`). Revoke and rotate any credentials that were exposed.

## Contributing

- Use tabs for indentation (project convention)
- Add Swagger comments for all new handlers
- Use Go-Jet for queries

## License

This project is private.
