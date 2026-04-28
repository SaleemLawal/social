package main

import (
	"log"
	"os"
	"time"

	"github.com/saleemlawal/social/internal/db"
	"github.com/saleemlawal/social/internal/store"
)

func main() {
	conn, err := db.New(os.Getenv("DB_URL"), 3, 3, 5 * time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	store := store.NewStorage(conn)
	db.Seed(store)
}