# PokéCLIpse Server

Text based online multiplayer Pokemon game server with auto-battling designed to be used with [PokéCLIpse Client](https://github.com/moceviciusda/pokeCLIpse-client).

Server exposes an API and utilizes websockets to handle user authentication, Pokemon battle simulation, and progression system. Server is built using Go and uses PostgreSQL for persistent storage.

## Table of Contents
- [Features](#features)
- [Installation](#installation)
  - [Prerequisites](#prerequisites)
  - [Setup](#setup)
- [Technologies](#technologies)
- [Deployment](#deployment)

## Features

- **User Authentication**: JWT-based authentication system
- **Battle System**: Auto-battling system for Pokemon
- **Progression System**: Catching, Leveling up, and evolving Pokemon
- **Database Integration**: PostgreSQL database for persistent storage
- **Caching**: In-memory caching for improved performance

## Installation

### Prerequisites

- Go 1.23 or higher
- PostgreSQL database

### Setup

1. Clone the repository
```bash
git clone https://github.com/yourusername/pokeclipse-server.git
cd pokeclipse-server
```

2. Configure environment variables by creating a .env file:
```
PORT=8080
DATABASE_URL=your_postgres_connection_string
JWT_SECRET=your_jwt_secret
```

3. Build the application
```bash
go build -o pokeclipse-server
```

4. Run database migrations
```bash
# Install Goose to run database migrations
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations
goose -dir sql/schema postgres "your_postgres_connection_string" up
```

5. Run the application
```bash
./pokeclipse-server
```

## Technologies

- [Go](https://golang.org/) - Programming language
- [Chi Router](https://github.com/go-chi/chi) - Lightweight HTTP router
- [JWT](https://github.com/golang-jwt/jwt) - JSON Web Token authentication
- [PostgreSQL](https://www.postgresql.org/) - Database
- [lib/pq](https://github.com/lib/pq) - PostgreSQL driver for Go
- [godotenv](https://github.com/joho/godotenv) - Environment variable loading

## Deployment

This project includes configuration for deployment to [Fly.io](https://fly.io) using the provided fly.toml file and Dockerfile that can be used with Prisma Postgres service
