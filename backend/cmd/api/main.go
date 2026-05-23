package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/sales-insight/backend/internal/analytics"
	analyticssqlc "github.com/sales-insight/backend/internal/analytics/sqlc"
	"github.com/sales-insight/backend/internal/auth"
	"github.com/sales-insight/backend/internal/clover"
	"github.com/sales-insight/backend/internal/employees"
	employeessqlc "github.com/sales-insight/backend/internal/employees/sqlc"
	"github.com/sales-insight/backend/internal/merchant"
	merchantsqlc "github.com/sales-insight/backend/internal/merchant/sqlc"
	"github.com/sales-insight/backend/internal/orders"
	orderssqlc "github.com/sales-insight/backend/internal/orders/sqlc"
	paymentssqlc "github.com/sales-insight/backend/internal/payments/sqlc"
	"github.com/sales-insight/backend/internal/platform/config"
	"github.com/sales-insight/backend/internal/platform/db"
	"github.com/sales-insight/backend/internal/platform/logger"
	"github.com/sales-insight/backend/internal/platform/scheduler"
	"github.com/sales-insight/backend/internal/platform/security"
	"github.com/sales-insight/backend/internal/products"
	productssqlc "github.com/sales-insight/backend/internal/products/sqlc"
	"github.com/sales-insight/backend/internal/sync"
	"github.com/sales-insight/backend/internal/token_cache"
)

func main() {
	// 1. Config
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. Logger
	log := logger.New(cfg.AppEnv)
	slog.SetDefault(log)

	// 3. Database
	postgres, err := db.NewPostgres(cfg.DatabaseDSN())
	if err != nil {
		log.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer postgres.Close()

	// 4. Security
	encrypter, err := security.NewEncrypter(cfg.Security.EncryptionKey)
	if err != nil {
		log.Error("failed to initialize encrypter", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// 5. Repositories
	merchantRepo := merchantsqlc.NewSQLCRepository(postgres.DB, encrypter)
	employeeRepo := employeessqlc.NewSQLCRepository(postgres.DB)
	productRepo := productssqlc.NewProductSQLCRepository(postgres.DB)
	categoryRepo := productssqlc.NewCategorySQLCRepository(postgres.DB)
	analyticCategoryRepo := productssqlc.NewAnalyticCategorySQLCRepository(postgres.DB)
	mappingRepo := productssqlc.NewCategoryMappingSQLCRepository(postgres.DB)
	orderRepo := orderssqlc.NewOrderSQLCRepository(postgres.DB)
	orderItemRepo := orderssqlc.NewOrderItemSQLCRepository(postgres.DB)
	categorySummaryRepo := orderssqlc.NewCategorySummarySQLCRepository(postgres.DB)
	paymentRepo := paymentssqlc.NewSQLCRepository(postgres.DB)
	analyticsRepo := analyticssqlc.NewSQLCRepository(postgres.DB)

	// 6. Services
	merchantSvc := merchant.NewService(merchantRepo)
	orderSvc := orders.NewService(orderRepo, orderItemRepo, categorySummaryRepo)
	employeeSvc := employees.NewService(employeeRepo, orderSvc)
	productSvc := products.NewService(productRepo, categoryRepo, analyticCategoryRepo, mappingRepo)
	// analyticsSvc will be wired into dashboard.Service in Phase 5
	_ = analytics.NewService(analyticsRepo)

	// 7. Auth & Clover
	tokenCache := tokencache.NewCache()
	oauthClient := clover.NewOAuthClient(cfg.Clover.Env)
	authService := auth.NewService(
		oauthClient,
		merchantRepo,
		encrypter,
		tokenCache,
		cfg.Clover.ClientID,
		cfg.Clover.ClientSecret,
	)
	authHandler := auth.NewHandler(authService, merchantSvc, oauthClient, cfg.Clover.ClientID, cfg.Server.FrontendURL)

	// 8. Rebuild token cache from DB on startup
	if err := authService.RebuildCache(context.Background()); err != nil {
		log.Warn("failed to rebuild token cache", slog.String("error", err.Error()))
	}

	// 9. Sync Engine & Scheduler
	cloverClient := clover.NewClient(cfg.Clover.Env)
		syncEngine := sync.NewEngine(
		cloverClient,
		tokenCache,
		merchantRepo,
		orderRepo,
		orderItemRepo,
		categorySummaryRepo,
		productRepo,
		categoryRepo,
		employeeRepo,
		paymentRepo,
	)
	syncEngine.SetBatchSize(100)

	cronScheduler := scheduler.NewScheduler(syncEngine, merchantRepo, log)
	cronScheduler.Start(context.Background())

	// 10. HTTP Server
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} ${uri} ${status} ${latency_human}\n",
	}))

	// Health check
	e.GET("/health", func(c echo.Context) error {
		if err := postgres.Health(c.Request().Context()); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "unhealthy", "db": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
	})

	// API v1
	v1 := e.Group("/api/v1")

	// Merchant routes
	merchantHandler := merchant.NewHandler(merchantSvc)
	merchantHandler.RegisterRoutes(v1)

	// Employee routes
	employeeHandler := employees.NewHandler(employeeSvc)
	employeeHandler.RegisterRoutes(v1)

	// Product routes
	productHandler := products.NewHandler(productSvc)
	productHandler.RegisterRoutes(v1)

	// Order routes
	orderHandler := orders.NewHandler(orderSvc)
	orderHandler.RegisterRoutes(v1)

	// Auth routes
	authHandler.RegisterRoutes(v1)

	// Start server in a goroutine
	go func() {
		addr := ":" + cfg.Server.Port
		log.Info("starting server", slog.String("addr", addr))
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Error("server error", slog.String("error", err.Error()))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server")

	// Stop scheduler before shutting down HTTP server
	cronScheduler.Stop()

	if cfg.AppEnv == "development" {
		// In development, close immediately to free the port right away
		if err := e.Close(); err != nil {
			log.Error("server close error", slog.String("error", err.Error()))
		}
	} else {
		// In production, graceful shutdown with timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(shutdownCtx); err != nil {
			log.Error("server shutdown error", slog.String("error", err.Error()))
		}
	}

	log.Info("server stopped")
}
