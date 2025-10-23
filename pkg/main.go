package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/K0ng2/zeedzad/config"
	"github.com/K0ng2/zeedzad/db"
	"github.com/K0ng2/zeedzad/handler"
	"github.com/K0ng2/zeedzad/server"
)

var port string

func main() {
	flag.StringVar(&port, "port", ":8088", "Server port")
	flag.Parse()

	database, err := db.NewDatabase(config.SQLITE_PATH)
	if err != nil {
		log.Fatalf("Failed to connect to SQLite: %v", err)
	}
	defer database.Close()

	handler := handler.NewHandler(database)

	// show config
	fmt.Println("Using configuration:")
	fmt.Printf("  SQLITE_PATH: %s\n", config.SQLITE_PATH)
	fmt.Printf("  YOUTUBE_API_KEY: %s\n", config.YOUTUBE_API_KEY)

	// Setup and start the router
	r := server.NewRouter(handler)
	if err := r.Listen(port); err != nil {
		panic(err)
	}
}
