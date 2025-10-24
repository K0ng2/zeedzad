package config

import (
	"os"
)

var (
	SQLITE_PATH        = os.Getenv("SQLITE_PATH")
	YOUTUBE_API_KEY    = os.Getenv("YOUTUBE_API_KEY")
	IGDB_CLIENT_ID     = os.Getenv("IGDB_CLIENT_ID")
	IGDB_CLIENT_SECRET = os.Getenv("IGDB_CLIENT_SECRET")
)
