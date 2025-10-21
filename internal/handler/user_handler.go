package handler

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/bugude101/user-age-api/internal/repository"
	"github.com/bugude101/user-age-api/internal/service"
)

type UserHandler struct {
	repo     *repository.UserRepo
	validate *validator.Validate
	logger   *zap.Logger
}

func NewUserHandler(repo *repository.UserRepo, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		repo:     repo,
		validate: validator.New(),
		logger:   logger,
	}
}

type createUpdateUserReq struct {
	Name string `json:"name" validate:"required"`
	DOB  string `json:"dob" validate:"required,datetime=2006-01-02"`
}

func (h *UserHandler) RegisterRoutes(app fiber.Router) {
	app.Post("/users", h.CreateUser)
	app.Get("/users", h.ListUsers)
	app.Get("/users/:id", h.GetUserByID)
	app.Put("/users/:id", h.UpdateUser)
	app.Delete("/users/:id", h.DeleteUser)
}

func (h *UserHandler) parseDOB(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req createUpdateUserReq
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("invalid body", zap.Error(err))
		return fiber.ErrBadRequest
	}
	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	dob, err := h.parseDOB(req.DOB)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "dob must be YYYY-MM-DD"})
	}

	ctx := context.Background()
	u, err := h.repo.Create(ctx, req.Name, dob)
	if err != nil {
		h.logger.Error("create user error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":   u.ID,
		"name": u.Name,
		"dob":  u.Dob.Format("2006-01-02"),
	})
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fiber.ErrBadRequest
	}

	ctx := context.Background()
	u, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.SendStatus(fiber.StatusNotFound)
		}
		h.logger.Error("get user error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user",
		})
	}

	age := service.CalculateAge(u.Dob, time.Now())
	return c.JSON(fiber.Map{
		"id":   u.ID,
		"name": u.Name,
		"dob":  u.Dob.Format("2006-01-02"),
		"age":  age,
	})
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fiber.ErrBadRequest
	}

	var req createUpdateUserReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	dob, err := h.parseDOB(req.DOB)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "dob must be YYYY-MM-DD"})
	}

	ctx := context.Background()
	u, err := h.repo.Update(ctx, id, req.Name, dob)
	if err != nil {
		h.logger.Error("update user error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}
	return c.JSON(fiber.Map{
		"id":   u.ID,
		"name": u.Name,
		"dob":  u.Dob.Format("2006-01-02"),
	})
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fiber.ErrBadRequest
	}
	ctx := context.Background()
	if err := h.repo.Delete(ctx, id); err != nil {
		h.logger.Error("delete user error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	limit := 100
	offset := 0
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	ctx := context.Background()
	users, err := h.repo.List(ctx, limit, offset)
	if err != nil {
		h.logger.Error("list users error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list users",
		})
	}

	resp := make([]fiber.Map, 0, len(users))
	now := time.Now()
	for _, u := range users {
		resp = append(resp, fiber.Map{
			"id":   u.ID,
			"name": u.Name,
			"dob":  u.Dob.Format("2006-01-02"),
			"age":  service.CalculateAge(u.Dob, now),
		})
	}
	return c.JSON(resp)
}
