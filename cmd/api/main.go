package main

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/saleemlawal/social/internal/auth"
	"github.com/saleemlawal/social/internal/db"
	"github.com/saleemlawal/social/internal/env"
	"github.com/saleemlawal/social/internal/mailer"
	"github.com/saleemlawal/social/internal/store"
	"github.com/saleemlawal/social/internal/store/cache"
	"go.uber.org/zap"
)

const version = "0.0.2"

//	@title			Social API
//	@description	A RESTful social networking API supporting user registration, posts, comments, following, and personalized feeds.

//	@contact.name	Saleem Lawal
//	@contact.url	https://github.com/saleemlawal/social

//	@license.name	MIT

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and JWT token.
func main() {
	appEnv := env.GetString("env", "dev")
	var logger *zap.SugaredLogger
	if appEnv == "production" {
		logger = zap.Must(zap.NewProduction()).Sugar()
	} else {
		logger = zap.Must(zap.NewDevelopment()).Sugar()
	}
	defer logger.Sync()

	if err := godotenv.Load(); err != nil {
		logger.Warn("no .env file found, relying on process environment")
	}

	cfg := &config{
		addr:        ":" + env.GetString("PORT", "8080"),
		env:         env.GetString("env", "dev"),
		apiURL:      env.GetString("API_URL", "localhost:8080"),
		frontendUrl: env.GetString("FRONTEND_URL", "http://localhost:3000"),
		db: dbConfig{
			addr:         env.GetString("DB_URL", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleTime:  time.Duration(env.GetInt("DB_MAX_OPEN_CONNS", 30)),
		},
		redisCfg: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:      env.GetString("REDIS_PW", ""),
			db:      env.GetInt("REDIS_DB", 0),
			enabled: env.GetBool("REDIS_ENABLED", true),
		},
		mail: struct {
			fromEmail string
			sendGrid  sendGridConfig
			mailtrap  mailtrapConfig
			exp       time.Duration
		}{
			fromEmail: env.GetString("FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
			mailtrap: mailtrapConfig{
				username: env.GetString("MAILTRAP_USERNAME", ""),
				password: env.GetString("MAILTRAP_PASSWORD", ""),
			},
			exp: time.Hour * 24 * 3,
		},
		auth: authConfig{
			basic: basicAuthConfig{
				username: env.GetString("BASIC_AUTH_USERNAME", "admin"),
				password: env.GetString("BASIC_AUTH_PASSWORD", "password"),
			},
			token: tokenConfig{
				secret:   env.GetString("JWT_SECRET", "example-secret"),
				audience: env.GetString("JWT_AUDIENCE", "social"),
				exp:      time.Hour * 24 * 3, // 3 days
				iss:      env.GetString("JWT_ISS", "social"),
			},
		},
	}

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)

	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("database connection pool established")

	// redis
	var redisClient *redis.Client = nil
	if cfg.redisCfg.enabled {
		redisClient = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)
		logger.Info("redis client established")
		defer redisClient.Close()
	}

	mailer := mailer.NewMailtrapMailer(cfg.mail.fromEmail, cfg.mail.mailtrap.username, cfg.mail.mailtrap.password)

	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.audience, cfg.auth.token.iss)

	store := store.NewStorage(db)
	app := &application{
		config:        *cfg,
		store:         store,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuthenticator,
		cacheStorage:  cache.NewRedisStorage(redisClient),
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
