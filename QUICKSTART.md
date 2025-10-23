# Quick Start Guide

This guide will help you get the Zeedzad app up and running quickly.

## Prerequisites

- Go 1.25+
- Bun (or npm/pnpm)
- YouTube Data API Key ([Get one here](https://console.cloud.google.com/apis/credentials))
- SQLite3 command-line tool

## Step 1: Database Setup

```bash
# Create database directory
mkdir -p data

# Set environment variable
export SQLITE_PATH="$PWD/data/zeedzad.db"

# Initialize database
sqlite3 $SQLITE_PATH < pkg/db/schema.sql

# Verify tables were created
sqlite3 $SQLITE_PATH ".tables"
# Should show: games  videos
```

## Step 2: Start Backend

```bash
cd pkg

# Install Go dependencies
go mod download

# Start the server
go run main.go

# Server will start on http://localhost:8088
```

Keep this terminal running.

## Step 3: Sync YouTube Videos

In a new terminal:

```bash
# Set your YouTube API key
export YOUTUBE_API_KEY="your_api_key_here"

# Run sync script
./scripts/sync-youtube.sh

# Or manually:
curl -X POST "http://localhost:8088/api/videos/sync?api_key=$YOUTUBE_API_KEY&max_results=50"
```

This will fetch the latest 50 videos from the OPZTV YouTube channel.

## Step 4: Start Frontend (Development)

In another new terminal:

```bash
cd web

# Install dependencies
bun install

# Start dev server
bun dev

# Frontend will start on http://localhost:3000
```

## Step 5: Use the App

1. Open your browser to `http://localhost:3000`
2. You should see a grid of video cards from OPZTV
3. Click "Match Game" on any video
4. Search for the game name
5. Select the correct game from Steam results
6. The video is now matched!

## Troubleshooting

### Backend won't start
- Check that `SQLITE_PATH` is set: `echo $SQLITE_PATH`
- Verify database exists: `ls -l $SQLITE_PATH`
- Check if port 8088 is in use: `lsof -i :8088`

### Frontend can't connect to API
- Verify backend is running: `curl http://localhost:8088/`
- Check browser console for errors
- Ensure CORS is working (should be enabled by default)

### No videos showing
- Run the sync script to fetch videos
- Check API response: `curl http://localhost:8088/api/videos`
- Verify database has data: `sqlite3 $SQLITE_PATH "SELECT COUNT(*) FROM videos;"`

### YouTube sync fails
- Verify API key is valid
- Check YouTube API quota limits
- Ensure channel username "OPZTV" is correct

## Next Steps

- Browse the API documentation at `http://localhost:8088/api/swagger/`
- Check the full README.md for production deployment
- Customize the frontend styling in `web/app/assets/css/main.css`

## Quick Commands Reference

```bash
# Start backend
cd pkg && go run main.go

# Start frontend dev
cd web && bun dev

# Sync videos
YOUTUBE_API_KEY=xxx ./scripts/sync-youtube.sh

# Build for production
cd web && bun run build
cd pkg && go build -o zeedzad

# Run production binary
export SQLITE_PATH=/path/to/db
./pkg/zeedzad
```

## API Quick Test

```bash
# Health check
curl http://localhost:8088/

# Get videos
curl http://localhost:8088/api/videos

# Search games on Steam
curl "http://localhost:8088/api/games/steam/search?q=team%20fortress%202"

# Get database health
curl http://localhost:8088/api/databasez
```

Enjoy using Zeedzad! 🎮📺
