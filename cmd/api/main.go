package main

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/saleemlawal/social/internal/db"
	"github.com/saleemlawal/social/internal/env"
	"github.com/saleemlawal/social/internal/store"
)

const version = "0.0.2"

//	@title			Swagger Example API
//	@description	This is a sample server Social API.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and JWT token.
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on process environment")
	}

	cfg := &config{
		addr: ":" + env.GetString("PORT", "8080"),
		db: dbConfig{
			addr:         env.GetString("DB_URL", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleTime:  time.Duration(env.GetInt("DB_MAX_OPEN_CONNS", 30)),
		},
		env:    env.GetString("env", "dev"),
		apiURL: env.GetString("API_URL", "localhost:8080"),
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
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
