package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"go-users-api/config"
	db "go-users-api/db/sqlc"
	"go-users-api/internal/handler"
	"go-users-api/internal/logger"
	"go-users-api/internal/middleware"
	"go-users-api/internal/repository"
	"go-users-api/internal/routes"
	"go-users-api/internal/service"
)

func main() {
	// Load config (DATABASE_URL, SERVER_ADDR)
	cfg := config.Load()

	// Create pgx connection pool
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}

	// Create sqlc Queries
	q := db.New(pool)

	// Logger
	logg, err := logger.New()
	if err != nil {
		log.Fatal(err)
	}
	defer logg.Sync()

	// Layers
	repo := repository.NewUserRepository(q)
	svc := service.NewUserService(repo)
	uh := handler.NewUserHandler(svc, logg)

	// Fiber app
	app := fiber.New()

	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger(logg))

	routes.Register(app, uh)

	if err := app.Listen(cfg.ServerAddr); err != nil {
		log.Fatal(err)
	}
}
