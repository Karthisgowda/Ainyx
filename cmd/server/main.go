package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Karthisgowda/Ainyx/config"
	"github.com/Karthisgowda/Ainyx/db/sqlc"
	"github.com/Karthisgowda/Ainyx/internal/handler"
	"github.com/Karthisgowda/Ainyx/internal/logger"
	"github.com/Karthisgowda/Ainyx/internal/middleware"
	"github.com/Karthisgowda/Ainyx/internal/repository"
	"github.com/Karthisgowda/Ainyx/internal/routes"
	"github.com/Karthisgowda/Ainyx/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	appLogger, err := logger.New(cfg.AppEnv)
	if err != nil {
		log.Fatalf("create logger: %v", err)
	}
	defer func() { _ = appLogger.Sync() }()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		appLogger.Fatal("connect database", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		appLogger.Fatal("ping database", zap.Error(err))
	}
	if err := runMigrations(ctx, pool); err != nil {
		appLogger.Fatal("run migrations", zap.Error(err))
	}

	queries := sqlc.New(pool)
	userRepo := repository.NewUserRepository(queries)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, appLogger)

	app := fiber.New(fiber.Config{
		AppName:      "Ainyx Users API",
		ErrorHandler: fiberErrorHandler(appLogger),
	})
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger(appLogger))
	routes.Register(app, userHandler)

	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			appLogger.Fatal("server stopped", zap.Error(err))
		}
	}()

	appLogger.Info("server started", zap.String("port", cfg.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		appLogger.Error("shutdown failed", zap.Error(err))
	}
}

func fiberErrorHandler(logger *zap.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if fiberErr, ok := err.(*fiber.Error); ok {
			code = fiberErr.Code
		}

		logger.Error("unhandled fiber error", zap.Error(err), zap.Int("status", code))
		return c.Status(code).JSON(fiber.Map{"error": err.Error()})
	}
}

func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	migration, err := os.ReadFile("db/migrations/001_create_users.sql")
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, string(migration))
	return err
}
