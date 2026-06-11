package handler

import (
	"errors"
	"strconv"

	"github.com/Karthisgowda/Ainyx/internal/models"
	"github.com/Karthisgowda/Ainyx/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type UserHandler struct {
	service   service.UserService
	validate  *validator.Validate
	logger    *zap.Logger
}

func NewUserHandler(service service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		service:  service,
		validate: validator.New(),
		logger:   logger,
	}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var request models.UserRequest
	if err := c.BodyParser(&request); err != nil {
		return badRequest(c, "invalid JSON body")
	}
	if err := h.validate.Struct(request); err != nil {
		return badRequest(c, err.Error())
	}

	user, err := h.service.Create(c.UserContext(), request)
	if err != nil {
		return h.handleError(c, err)
	}

	h.logger.Info("created user", zap.Int32("user_id", user.ID))
	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badRequest(c, "invalid user id")
	}

	user, err := h.service.GetByID(c.UserContext(), id)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(user)
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	limit := queryInt32(c, "limit", 20)
	offset := queryInt32(c, "offset", 0)
	if limit < 1 || limit > 100 {
		return badRequest(c, "limit must be between 1 and 100")
	}
	if offset < 0 {
		return badRequest(c, "offset cannot be negative")
	}

	users, err := h.service.List(c.UserContext(), limit, offset)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(users)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badRequest(c, "invalid user id")
	}

	var request models.UserRequest
	if err := c.BodyParser(&request); err != nil {
		return badRequest(c, "invalid JSON body")
	}
	if err := h.validate.Struct(request); err != nil {
		return badRequest(c, err.Error())
	}

	user, err := h.service.Update(c.UserContext(), id, request)
	if err != nil {
		return h.handleError(c, err)
	}

	h.logger.Info("updated user", zap.Int32("user_id", user.ID))
	return c.JSON(user)
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return badRequest(c, "invalid user id")
	}

	if err := h.service.Delete(c.UserContext(), id); err != nil {
		return h.handleError(c, err)
	}

	h.logger.Info("deleted user", zap.Int32("user_id", id))
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *UserHandler) handleError(c *fiber.Ctx, err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "user not found"})
	}
	if errors.Is(err, service.ErrFutureDOB) {
		return badRequest(c, err.Error())
	}

	h.logger.Error("request failed", zap.Error(err))
	return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "internal server error"})
}

func parseID(c *fiber.Ctx) (int32, error) {
	id, err := strconv.ParseInt(c.Params("id"), 10, 32)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return int32(id), nil
}

func queryInt32(c *fiber.Ctx, key string, fallback int32) int32 {
	value := c.Query(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return fallback
	}
	return int32(parsed)
}

func badRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: message})
}
