package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/bugude101/user-age-api/internal/logger"
	"github.com/bugude101/user-age-api/internal/middleware"
	"github.com/bugude101/user-age-api/internal/repository"
	"github.com/bugude101/user-age-api/internal/routes"
)

func main() {
	// load config from env (simple)
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:bugude@localhost:5432/postgres?sslmode=disable"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// create logger
	logg, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logg.Sync()

	// open DB
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logg.Fatal("failed to open db", zap.Error(err))
	}
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	// repository
	repo := repository.NewUserRepo(db)

	app := fiber.New(fiber.Config{
		Prefork: false,
	})

	// middlewares
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger(logg))

	routes.Register(app, repo, logg)

	addr := fmt.Sprintf(":%s", port)
	logg.Info("starting server", zap.String("addr", addr))
	if err := app.Listen(addr); err != nil {
		logg.Fatal("server exited", zap.Error(err))
	}
}
