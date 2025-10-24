# Migration Guide: SQLite to Cloudflare D1

This document explains the migration from local SQLite to Cloudflare D1 database.

## What Changed

### Database Layer
- **Before**: Local SQLite database using `modernc.org/sqlite`
- **After**: Cloudflare D1 (serverless SQLite) via REST API using `github.com/cloudflare/cloudflare-go/v6`

### Environment Variables
- **Removed**: `SQLITE_PATH`
- **Added**:
  - `D1_ACCOUNT_ID` - Your Cloudflare account ID
  - `D1_DATABASE_ID` - Your D1 database ID (UUID)
  - `CLOUDFLARE_API_TOKEN` - API token with D1 permissions

## Setup Instructions

### 1. Create Cloudflare D1 Database

```bash
# Install Wrangler CLI (Cloudflare's CLI tool)
npm install -g wrangler

# Login to Cloudflare
wrangler login

# Create a new D1 database
wrangler d1 create zeedzad

# Note the database ID from the output
```

### 2. Migrate Your Schema

If you have an existing SQLite database, export the schema:

```bash
# Export schema from local SQLite
sqlite3 $SQLITE_PATH .dump > schema.sql

# Import to D1 (replace DATABASE_ID with your actual ID)
wrangler d1 execute zeedzad --file=pkg/db/schema.sql
```

Or use the schema file directly:

```bash
wrangler d1 execute zeedzad --file=pkg/db/schema.sql
```

### 3. Get Your Cloudflare Credentials

#### Account ID
1. Go to https://dash.cloudflare.com/
2. Your account ID is in the URL or right sidebar

#### Database ID
- You received this when creating the database with `wrangler d1 create`
- Or find it in the Cloudflare dashboard under Workers & Pages → D1

#### API Token
1. Go to https://dash.cloudflare.com/profile/api-tokens
2. Click "Create Token"
3. Use "Edit Cloudflare Workers" template or create custom token
4. Required permissions:
   - Account → D1 → Read
   - Account → D1 → Write
5. Copy the generated token

### 4. Update Environment Variables

Copy `.env.example` to `.env` and fill in your values:

```bash
cp .env.example .env
```

Edit `.env`:

```bash
D1_ACCOUNT_ID=your_cloudflare_account_id
D1_DATABASE_ID=your_d1_database_id
CLOUDFLARE_API_TOKEN=your_cloudflare_api_token
YOUTUBE_API_KEY=your_youtube_api_key
IGDB_CLIENT_ID=your_igdb_client_id
IGDB_CLIENT_SECRET=your_igdb_client_secret
```

### 5. Run the Application

```bash
cd pkg
go run main.go
```

## Key Differences

### Performance Considerations
- **Latency**: D1 REST API has network latency (~50-200ms per request)
- **Best for**: Applications where network latency is acceptable
- **Not ideal for**: High-frequency, low-latency operations

### Transaction Support
- D1 REST API doesn't support traditional SQL transactions
- Each query is executed atomically
- Batch operations are supported via semicolon-separated statements

### Query Execution
- Queries are sent as HTTP requests to Cloudflare's API
- Results are returned as JSON
- The database layer implements a `database/sql` compatible driver wrapper

### Foreign Key Handling
- **Write operations** (INSERT, UPDATE, DELETE) automatically include `PRAGMA defer_foreign_keys = on`
- Parameters are **inlined** into SQL for write operations (batch queries don't support `?` placeholders in D1)
- Read operations use standard parameterized queries for security
- This allows NULL foreign keys while still enforcing constraints at transaction end

### DateTime Handling
- D1 returns datetime values as RFC3339 strings (e.g., `"2025-10-22T09:00:27Z"`)
- The database layer automatically converts these to SQLite datetime format (`"2006-01-02 15:04:05"`)
- This ensures compatibility with go-jet's SQLite expectations
- Both RFC3339 and SQLite formats are supported on read

## Advantages of D1

1. **No Local Database**: No need to manage database files
2. **Serverless**: Scales automatically
3. **Global Distribution**: Data replicated across Cloudflare's network
4. **Backup & Recovery**: Built-in Time Travel (30 days)
5. **No CGO**: Pure Go implementation (cloudflare-go SDK)

## Troubleshooting

### Connection Issues
- Verify your API token has D1 permissions
- Check account ID and database ID are correct
- Ensure you have network connectivity

### Query Errors
- D1 uses SQLite syntax (same as before)
- Check query syntax in Wrangler: `wrangler d1 execute zeedzad --command="SELECT * FROM videos LIMIT 1"`

### Foreign Key Constraint Errors
- The database layer automatically adds `PRAGMA defer_foreign_keys = on` for write operations (INSERT, UPDATE, DELETE)
- This allows NULL foreign keys (e.g., videos without assigned games) to be inserted temporarily
- Constraints are still enforced at the end of each transaction

### Performance Issues
- Consider batching queries when possible
- Use indexed columns for searches
- Monitor query duration in response metadata

## Rolling Back

If you need to revert to local SQLite:

1. Checkout the previous commit before this migration
2. Restore `.env` with `SQLITE_PATH`
3. Run `go mod tidy` to restore old dependencies

## Additional Resources

- [Cloudflare D1 Documentation](https://developers.cloudflare.com/d1/)
- [D1 API Reference](https://developers.cloudflare.com/api/resources/d1/)
- [Wrangler CLI Docs](https://developers.cloudflare.com/workers/wrangler/)
