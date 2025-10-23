# Go-Jet Model Generation Instructions

## Prerequisites

1. Ensure you have a SQLite database file created
2. Apply the schema from `pkg/db/schema.sql` to your database

## Generate Go-Jet Models

Run the following command from the `pkg` directory:

```bash
# Set your database path
export SQLITE_PATH="/path/to/your/database.db"

# Apply the schema
sqlite3 $SQLITE_PATH < db/schema.sql

# Generate Go-Jet models
jet -source=sqlite -dsn=$SQLITE_PATH -path=./gen
```

This will generate type-safe Go models in `pkg/gen/` directory based on your database schema.

## Expected Output Structure

After running the command, you should have:

```
pkg/gen/
  zeedzad/        # or your database name
    public/
      table/
        games.go
        videos.go
      model/
        games.go
        videos.go
```

## Import in Your Code

```go
import (
	"github.com/K0ng2/zeedzad/gen/zeedzad/public/table"
	. "github.com/go-jet/jet/v2/sqlite"
)

// Use in queries
var Games = table.Games
var Videos = table.Videos
```

## Alternative: Manual Installation of jet CLI

If `jet` command is not available:

```bash
go install github.com/go-jet/jet/v2/cmd/jet@latest
```

Then run the generation command above.
