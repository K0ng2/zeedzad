package config

import (
	"os"
)

var (
	SQLITE_PATH = os.Getenv("SQLITE_PATH")
)
