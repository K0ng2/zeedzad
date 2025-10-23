# Zeedzad Project Guide

## Architecture Overview

Zeedzad is a full-stack web application with:
- **Backend**: Go (Fiber v3) REST API with SQLite (pure Go via modernc.org/sqlite) in `pkg/`
- **Frontend**: Nuxt 4 application in `web/`
- **Deployment**: Backend embeds frontend static assets via Go embed

### Service Boundaries

```
┌─────────────────────────────────────────────────┐
│  Fiber HTTP Server (pkg/server/server.go)      │
│  ├── / - Healthcheck (Fiber's built-in)         │
│  ├── /api/* - REST API endpoints                │
│  ├── /api/swagger/* - Auto-generated API docs   │
│  ├── /api/databasez - Database health check     │
│  └── /* - Static files from embedded web build  │
└─────────────────────────────────────────────────┘
```

The backend serves both API routes and the compiled Nuxt app. When frontend builds complete, assets are placed in `pkg/web/public/` where `web/embed.go` picks them up via `//go:embed **/* all:public`.

### Layer Architecture (Backend)

1. **Handler** (`pkg/handler/`) - HTTP request/response, validation, Swagger docs
2. **Repository** (`pkg/repository/`) - Data access using go-jet/jet SQL builder
3. **Database** (`pkg/db/`) - SQLite connection management and executor interface
4. **Model** (`pkg/model/`) - Shared data structures and API response formats

**Data Flow**: Handler → Repository → DB → SQLite

### Configuration

- **Environment**: `SQLITE_PATH` environment variable required (path to SQLite database file)
- **Server Port**: Default `:8088`, override with `-port` flag: `go run main.go -port :3000`
- **Database Driver**: Uses pure Go SQLite (modernc.org/sqlite) - no CGo required

## Critical Workflows

### Development Build & Run

```bash
# Set required environment variable
export SQLITE_PATH="/path/to/database.db"

# Start Go backend (from pkg/)
cd pkg
go run main.go                # Default port :8088
go run main.go -port :3000    # Custom port

# Start Nuxt dev server (from web/)
cd web
bun dev  # or npm/pnpm dev
```

Backend runs on port :8088 by default, Nuxt on `http://localhost:3000` in dev mode.

### Generate Swagger Documentation

```bash
mise run swag  # or: swag init -d pkg -g server/server.go -o pkg/docs --ot go
```

Run this after modifying handler comments with `@Summary`, `@Description`, etc. Swagger is served at `/api/swagger/`.

### Production Build

```bash
# 1. Build Nuxt frontend
cd web && bun run build  # generates web/.output/public/

# 2. Copy assets to embed location
# (Manual step - copy web/.output/public/* to pkg/web/public/)

# 3. Build Go binary
cd pkg && go build -o zeedzad
```

The Go binary will contain the embedded web assets and can run standalone.

## Project-Specific Conventions

### API Response Pattern

All API responses use the generic `model.APIResponse[T]` wrapper:

```go
type APIResponse[T any] struct {
    Data T     `json:"data"`
    Meta *Meta `json:"meta,omitempty"`  // pagination metadata
}
```

Use the `handler.Response()` helper to construct responses with optional metadata.

### Database Query Pattern

- **Never write raw SQL** - use go-jet/jet's type-safe SQL builder
- Repository methods accept `context.Context` as first parameter
- Use `repository.WithTx()` for transaction-based operations
- Helper functions: `NullString()`, `NullInt16()`, `ILIKE()` for nullable/special queries
- `TotalItems()` standardizes counting with optional WHERE clauses

### Fiber Handler Registration

Handlers are methods on `handler.Handler` struct:

```go
func (h *Handler) EndpointName(c fiber.Ctx) error {
    ctx := c.RequestCtx()  // Extract fasthttp context for database operations
    // Use h.repo for database access with ctx
}
```

Register in `server.NewRouter()` using `app.Get/Post()` with handler method references.

### Error Handling

- Repository errors: wrap with `repository.FormatError(prefix, err)`
- HTTP errors: return `c.Status(code).JSON(response)` directly
- Database health check pattern: see `handler/health.go` for status degradation

### Swagger Documentation

Every exported handler must have godoc comments:

```go
// EndpointName godoc
// @Summary Brief description
// @Description Detailed description
// @Tags category
// @Accept json
// @Produce json
// @Param id path string true "ID description"
// @Success 200 {object} model.ResponseType
// @Router /endpoint [get]
func (h *Handler) EndpointName(c fiber.Ctx) error
```

Run `mise run swag` after changes.

## Tool Management

This project uses **mise** (see `mise.toml`) for version management:

- Go 1.25
- Bun 1.3 (for frontend package management)
- Custom task: `mise run swag` for Swagger generation

## Key Files

- `pkg/server/server.go` - Fiber app configuration, middleware, route registration
- `pkg/web/embed.go` - Frontend asset embedding (critical for production builds)
- `pkg/repository/repository.go` - Database helpers and query patterns
- `pkg/handler/handler.go` - Common handler utilities (GetOffset, Response, etc.)
- `web/nuxt.config.ts` - Nuxt configuration (minimal default setup)

## Code Style

- **Indentation**: Tabs (width 2) - enforced by `.editorconfig`
- **Line endings**: CRLF (Windows-style)
- **Package naming**: Single-word lowercase (matches directory name)
- **Error messages**: Lowercase, no ending punctuation

## Notes

- Health check endpoint at `/api/databasez` (note the 'z')
- Startup health check at root: uses Fiber's built-in `healthcheck` middleware
- Static file serving includes custom 404 handler with `404.html` fallback
