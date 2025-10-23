# Zeedzad - YouTube Video to Game Matcher

A web application for matching YouTube videos from [OPZTV](https://www.youtube.com/@OPZTV) with Steam games.

## Features

- 📺 **Video Management**: Display YouTube videos in a card-based layout with thumbnails and metadata
- 🎮 **Game Matching**: Match videos with Steam games using the Steam Community API
- 🔍 **Search**: Search videos by title or game name
- 📄 **Pagination**: Browse videos with 24 items per page
- 🎨 **Modern UI**: Built with Nuxt 4, Tailwind CSS, and DaisyUI components

## Tech Stack

### Backend
- **Go**: Main programming language
- **Fiber v3**: HTTP server framework
- **SQLite**: Database (pure Go via modernc.org/sqlite)
- **Go-Jet**: Type-safe SQL query builder
- **YouTube Data API v3**: Fetch videos from YouTube
- **Swagger**: API documentation

### Frontend
- **Nuxt 4**: Vue.js framework
- **Tailwind CSS**: Utility-first CSS framework
- **DaisyUI**: Component library
- **FontAwesome**: Icon library

## Prerequisites

- Go 1.25+
- Bun 1.3+ (or npm/pnpm)
- YouTube Data API Key (for syncing videos)
- mise (optional, for tool management)

## Setup

### 1. Database Setup

```bash
# Set database path
export SQLITE_PATH="/path/to/database.db"

# Initialize database with schema
sqlite3 $SQLITE_PATH < pkg/db/schema.sql
```

### 2. Backend Setup

```bash
cd pkg

# Install dependencies
go mod download

# Run development server (default port :8088)
go run main.go

# Or specify custom port
go run main.go -port :3000
```

### 3. Frontend Setup

```bash
cd web

# Install dependencies
bun install

# Run development server
bun dev
```

The frontend will be available at `http://localhost:3000`.

## API Endpoints

### Videos
- `GET /api/videos` - Get all videos (paginated)
  - Query params: `offset`, `limit`, `search`
- `GET /api/videos/:id` - Get video by ID
- `PUT /api/videos/:id/game` - Match video with game
- `POST /api/videos/sync` - Sync videos from YouTube
  - Query params: `api_key` (required), `max_results` (optional, default: 50)

### Games
- `GET /api/games` - Get all games (paginated)
  - Query params: `offset`, `limit`, `search`
- `GET /api/games/:id` - Get game by ID
- `POST /api/games` - Create new game
- `GET /api/games/steam/search` - Search Steam games
  - Query param: `q` (search query)

### Health
- `GET /` - Health check
- `GET /api/databasez` - Database health check

### Documentation
- `GET /api/swagger/` - Swagger API documentation

## Usage

### 1. Sync YouTube Videos

First, sync videos from the OPZTV YouTube channel:

```bash
curl -X POST "http://localhost:8088/api/videos/sync?api_key=YOUR_YOUTUBE_API_KEY&max_results=50"
```

Or use the Swagger UI at `http://localhost:8088/api/swagger/`

### 2. Browse and Match Videos

1. Open the web app at `http://localhost:3000`
2. Browse videos in the card layout
3. For unmatched videos, click "Match Game"
4. Search for the game name
5. Select the correct game from Steam search results
6. The video will be automatically matched with the game

### 3. Search Videos

Use the search bar to find videos by:
- Video title
- Game name (for matched videos)

## Production Build

### 1. Build Frontend

```bash
cd web
bun run build
```

This generates static files in `web/.output/public/`

### 2. Copy Assets

```bash
# Copy built assets to embed location
cp -r web/.output/public/* pkg/web/public/
```

### 3. Build Backend

```bash
cd pkg
go build -o zeedzad
```

### 4. Run Production Binary

```bash
export SQLITE_PATH="/path/to/database.db"
./zeedzad
```

The binary includes embedded frontend assets and serves both API and web app.

## Development

### Generate Go-Jet Models

After modifying `pkg/db/schema.sql`:

```bash
cd pkg
jet -dsn=/tmp/zeedzad.db -schema=file://db/schema.sql -path=./repository/table
```

### Generate Swagger Docs

After modifying handler godoc comments:

```bash
cd pkg
swag init -d . -g server/server.go -o docs --ot go
```

Or use mise:

```bash
mise run swag
```

## Project Structure

```
zeedzad/
├── pkg/                    # Backend Go code
│   ├── config/            # Configuration
│   ├── db/                # Database connection
│   ├── docs/              # Swagger documentation
│   ├── handler/           # HTTP handlers
│   ├── model/             # API models
│   ├── repository/        # Database layer
│   │   ├── model/         # Generated Go-Jet models
│   │   └── table/         # Generated Go-Jet tables
│   ├── server/            # Fiber server setup
│   ├── web/               # Embedded frontend assets
│   └── main.go            # Entry point
└── web/                   # Frontend Nuxt app
    ├── app/
    │   ├── assets/        # CSS and styles
    │   ├── components/    # Vue components
    │   ├── composables/   # Composables (API, etc.)
    │   ├── layouts/       # Layouts
    │   ├── pages/         # Pages
    │   └── plugins/       # Plugins
    └── nuxt.config.ts     # Nuxt configuration
```

## Configuration

### Backend
- `SQLITE_PATH`: Path to SQLite database (required)
- `-port`: Server port (default: :8088)

### Frontend (Development)
- `NUXT_PUBLIC_API_BASE`: API base URL (default: /api)

## API Response Format

All API responses follow this structure:

```json
{
  "data": <T>,
  "meta": {
    "total": 100,
    "limit": 24,
    "offset": 0
  }
}
```

## Contributing

1. Follow the project conventions in `.github/copilot-instructions.md`
2. Use tabs for indentation (width 2)
3. Add Swagger documentation to all endpoints
4. Use Go-Jet for database queries
5. Run `mise run swag` after handler changes

## License

This project is private and not licensed for public use.
