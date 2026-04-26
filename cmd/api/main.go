package main

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/saleemlawal/social/internal/db"
	"github.com/saleemlawal/social/internal/env"
	"github.com/saleemlawal/social/internal/store"
)

const version = "0.0.1"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := &config{
		addr: ":" + env.GetString("PORT", "3000"),
		db: dbConfig{
			addr: env.GetString("DB_URL", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleTime:  time.Duration(env.GetInt("DB_MAX_OPEN_CONNS", 30)),
		},
		env: env.GetString("env", "dev"),
	}

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)

	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Printf("database connection pool established")

	store := store.NewStorage(db)
	app := &application{
		config: *cfg,
		store: store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
