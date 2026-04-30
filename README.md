```
 ________  ________  ________  ___  ________  ___          ________  ________  ___
|\   ____\|\   __  \|\   ____\|\  \|\   __  \|\  \        |\   __  \|\   __  \|\  \
\ \  \___|\ \  \|\  \ \  \___|\ \  \ \  \|\  \ \  \       \ \  \|\  \ \  \|\  \ \  \
 \ \_____  \ \  \\\  \ \  \    \ \  \ \   __  \ \  \       \ \   __  \ \   ____\ \  \
  \|____|\  \ \  \\\  \ \  \____\ \  \ \  \ \  \ \  \____   \ \  \ \  \ \  \___|\ \  \
    ____\_\  \ \_______\ \_______\ \__\ \__\ \__\ \_______\   \ \__\ \__\ \__\    \ \__\
   |\_________\|_______|\|_______|\|__|\|__|\|__|\|_______|    \|__|\|__|\|__|     \|__|
   \|_________|
```

---

```
  [ REST API ]  ·  [ Go 1.25 ]  ·  [ PostgreSQL 17 ]  ·  [ Docker ]  ·  [ Swagger ]
```

---

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Environment Variables](#environment-variables)
- [API Reference](#api-reference)
- [Database Migrations](#database-migrations)
- [Development](#development)
- [Make Commands](#make-commands)

---

## Overview

**Social API** is a RESTful backend service built in Go that powers a social networking platform. It supports user registration with email-based account activation, post creation, commenting, social following, and a personalized feed — all served through a clean, versioned HTTP API.

---

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                        Client                           │
└────────────────────────┬────────────────────────────────┘
                         │  HTTP
┌────────────────────────▼────────────────────────────────┐
│                    chi Router /v1                        │
│                                                         │
│   Middleware Stack                                      │
│   ─────────────────────────────────────                 │
│   Logger · RequestID · RealIP · Recoverer · Timeout     │
│                                                         │
│   Routes                                                │
│   ─────────────────────────────────────                 │
│   /health          /posts           /users              │
│   /swagger         /authentication  /users/feeds        │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│                  Store Layer (internal/store)            │
│                                                         │
│   UserStore · PostStore · CommentStore · FollowerStore  │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│                   PostgreSQL 17                          │
└─────────────────────────────────────────────────────────┘
```

---

## Features

```
  USERS                          POSTS                         SOCIAL
  ─────────────────────────      ────────────────────────      ─────────────────────────
  + Register with token          + Create post                 + Follow a user
  + Email activation flow        + Get post by ID              + Unfollow a user
  + Get user profile             + Update post (PATCH)         + Personalized feed
                                 + Delete post
                                 + Add comment to post
```

---

## Tech Stack

| Layer          | Technology                              |
| -------------- | --------------------------------------- |
| Language       | Go 1.25                                 |
| Router         | chi v5                                  |
| Database       | PostgreSQL 17 (Alpine)                  |
| Migrations     | golang-migrate v4                       |
| Validation     | go-playground/validator v10             |
| Documentation  | Swagger (swaggo/swag + http-swagger)    |
| Logging        | Uber Zap (sugared logger)               |
| Hot Reload     | Air                                     |
| Containerisation | Docker + Docker Compose               |
| Password Hashing | bcrypt                                |
| Token Hashing  | SHA-256                                 |

---

## Project Structure

```
social/
├── cmd/
│   ├── api/
│   │   ├── main.go          # Entry point, config bootstrap
│   │   ├── api.go           # Application struct, router, server
│   │   ├── auth.go          # Registration handler
│   │   ├── users.go         # User, follow, unfollow, activate handlers
│   │   ├── posts.go         # Post CRUD handlers
│   │   ├── errors.go        # Centralised error helpers
│   │   └── json.go          # JSON read/write helpers
│   └── migrate/
│       ├── migrations/      # Sequential SQL migration files
│       └── seed/            # Database seed data
├── internal/
│   ├── store/
│   │   ├── storage.go       # Storage interface
│   │   ├── users.go         # User & follower DB operations
│   │   ├── posts.go         # Post DB operations
│   │   └── comments.go      # Comment DB operations
│   ├── db/                  # DB connection setup
│   └── env/                 # Environment variable helpers
├── docs/                    # Auto-generated Swagger docs
├── scripts/                 # DB init scripts
├── compose.yml              # Production-like Compose setup
├── compose.override.yml     # Dev overrides (hot reload, volume mounts)
├── Dockerfile               # Multi-stage: builder / dev / runtime
└── Makefile
```

---

## Getting Started

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and Docker Compose
- [Go 1.25+](https://go.dev/dl/) (for running outside Docker)
- [golang-migrate CLI](https://github.com/golang-migrate/migrate) (for running migrations locally)

### Run with Docker

```bash
# Clone the repository
git clone https://github.com/saleemlawal/social.git
cd social

# Copy the example env file and fill in your values
cp .env.example .env

# Build and start all services (API + DB + migrations)
make run
```

The API will be available at `http://localhost:8080/v1`.
Swagger UI is available at `http://localhost:8080/v1/swagger/index.html`.

---

## Environment Variables

Create a `.env` file in the project root:

```env
PORT=8080
ENV=dev
API_URL=localhost:8080

# PostgreSQL connection string
DB_URL=postgres://admin:adminpassword@localhost/social?sslmode=disable

# Connection pool
DB_MAX_OPEN_CONNS=30
DB_MAX_IDLE_CONNS=30
DB_MAX_IDLE_TIME=30
```

---

## API Reference

All routes are prefixed with `/v1`. Full interactive documentation is available via Swagger at `/v1/swagger/index.html`.

### Health

```
GET    /v1/health
```

### Authentication

```
POST   /v1/authentication/user       Register a new user (returns activation token)
```

### Users

```
PUT    /v1/users/activate/{token}    Activate a user account
GET    /v1/users/{id}                Get a user by ID
PUT    /v1/users/{id}/follow         Follow a user
PUT    /v1/users/{id}/unfollow       Unfollow a user
GET    /v1/users/feeds               Get the personalized feed
```

### Posts

```
POST   /v1/posts                     Create a post
GET    /v1/posts/{postId}            Get a post by ID
PATCH  /v1/posts/{postId}            Update a post
DELETE /v1/posts/{postId}            Delete a post
POST   /v1/posts/{postId}/comments   Add a comment to a post
```

### Registration Flow

```
  1.  POST /v1/authentication/user
        └── Creates user (inactive)
        └── Generates SHA-256 hashed token stored in user_invitations
        └── Returns plain token in response

  2.  PUT /v1/users/activate/{token}
        └── Looks up token (checks expiry)
        └── Sets user.activated = true
        └── Deletes the invitation record
```

---

## Database Migrations

Migrations live in `cmd/migrate/migrations/` and are named sequentially:

```
000001  create users table
000002  create posts table
000003  add comments table
000004  delete cascade constraints
000005  add version to posts
000006  alter comments table
000007  add followers table
000008  add indexes (GIN, btree)
000009  add user invitations
000010  add activated column to users
000011  alter user invitations
000012  rename followers.user_id -> followed_id
```

```bash
# Run all pending migrations
make migrate-up

# Roll back N steps
make migrate-down 1

# Create a new migration file
make migration <name>

# Force a version (recover from dirty state)
make migrate-force <version>
```

---

## Development

The project uses a multi-stage Dockerfile:

| Stage     | Purpose                                              |
| --------- | ---------------------------------------------------- |
| `builder` | Compiles the binary with `CGO_ENABLED=0`             |
| `dev`     | Hot reload via Air, installs `make` and `swag`       |
| `runtime` | Minimal Alpine image, runs the pre-built binary      |

The `compose.override.yml` activates the `dev` target automatically and mounts the source directory into the container so Air can detect changes and rebuild.

```bash
# Start in development mode with live reload
make run

# Seed the database
make seed

# Regenerate Swagger docs
make swagger

# Format Go source files
make format
```

---

## Make Commands

| Command              | Description                                      |
| -------------------- | ------------------------------------------------ |
| `make run`           | Build and start services via Docker Compose      |
| `make migrate-up`    | Apply all pending migrations                     |
| `make migrate-down`  | Roll back migrations (pass a number)             |
| `make migrate-force` | Force migration version and re-run up            |
| `make migration`     | Create a new named migration file                |
| `make seed`          | Seed the database with sample data               |
| `make swagger`       | Generate and format Swagger documentation        |
| `make format`        | Run `gofmt` across the project                   |

---

```
  Built with Go · PostgreSQL · Docker
  github.com/saleemlawal/social
```
