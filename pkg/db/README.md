# Go-Jet Model Generation

This project keeps an SQL schema at `pkg/db/schema.sql`. Instead of requiring a live SQLite database for generation, you can generate Go-Jet models from the schema file directly (recommended).

## Generate models from schema file (recommended)

Run from the repository root or `pkg`:

```bash
cd pkg
# Generate Go-Jet tables/models using the schema file
jet -dsn=file://db/schema.sql -path=./repository/table
```

This places generated table and model files under `pkg/repository/table` which the codebase uses.

## Optional: Local SQLite workflow

If you prefer to generate models from a real SQLite DB (local testing), you can create a database and apply the schema:

```bash
mkdir -p data
export SQLITE_PATH="$PWD/data/zeedzad.db"
sqlite3 $SQLITE_PATH < db/schema.sql

# Then generate using the local DB
jet -source=sqlite -dsn=$SQLITE_PATH -path=./repository/table
```

## Expected structure

After generation you should have Go-Jet table and model files available under `pkg/repository/table` (or the path you supplied).

## Installing jet CLI

If `jet` is not installed:

```bash
go install github.com/go-jet/jet/v2/cmd/jet@latest
```

## Notes

- Use the `file://` schema approach when you're generating models in CI or on machines without a local SQLite binary.
- Keep generated files out of git if you prefer (this repo keeps them under `pkg/repository/table`).
