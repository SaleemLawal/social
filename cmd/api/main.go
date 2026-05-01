package main

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/saleemlawal/social/internal/auth"
	"github.com/saleemlawal/social/internal/db"
	"github.com/saleemlawal/social/internal/env"
	"github.com/saleemlawal/social/internal/mailer"
	"github.com/saleemlawal/social/internal/store"
	"go.uber.org/zap"
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

	mailer := mailer.NewMailtrapMailer(cfg.mail.fromEmail, cfg.mail.mailtrap.username, cfg.mail.mailtrap.password)

	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.audience, cfg.auth.token.iss)

	store := store.NewStorage(db)
	app := &application{
		config:        *cfg,
		store:         store,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuthenticator,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
