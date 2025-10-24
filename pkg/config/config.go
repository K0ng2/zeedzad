package config

import (
	"os"
)

var (
	D1_ACCOUNT_ID        = os.Getenv("D1_ACCOUNT_ID")
	D1_DATABASE_ID       = os.Getenv("D1_DATABASE_ID")
	CLOUDFLARE_API_TOKEN = os.Getenv("CLOUDFLARE_API_TOKEN")
	YOUTUBE_API_KEY      = os.Getenv("YOUTUBE_API_KEY")
	IGDB_CLIENT_ID       = os.Getenv("IGDB_CLIENT_ID")
	IGDB_CLIENT_SECRET   = os.Getenv("IGDB_CLIENT_SECRET")
)
