package routes

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/bugude99/user-age-api/internal/handler"
	"github.com/bugude99/user-age-api/internal/repository"
)

func Register(app *fiber.App, repo *repository.UserRepo, logger *zap.Logger) {
	logger.Info("hello")
	uh := handler.NewUserHandler(repo, logger)
	api := app.Group("/api")
	uh.RegisterRoutes(api)
}
