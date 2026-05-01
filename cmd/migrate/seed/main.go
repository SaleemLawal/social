package main

import (
	"log"
	"os"
	"time"

	"github.com/saleemlawal/social/internal/db"
	"github.com/saleemlawal/social/internal/store"
)

func main() {
	// Enough connections for bounded parallel seed work (see db.Seed); a tiny pool + huge fan-out hits store query timeouts.
	conn, err := db.New(os.Getenv("DB_URL"), 25, 10, 5*time.Minute)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	store := store.NewStorage(conn)
	db.Seed(store, conn)
}
