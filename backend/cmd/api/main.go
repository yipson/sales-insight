package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sales-insight/backend/internal/analytics"
	analyticssqlc "github.com/sales-insight/backend/internal/analytics/sqlc"
	"github.com/sales-insight/backend/internal/api"
	"github.com/sales-insight/backend/internal/auth"
	"github.com/sales-insight/backend/internal/clover"
	"github.com/sales-insight/backend/internal/dashboard"
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
	syncsqlc "github.com/sales-insight/backend/internal/sync/sqlc"
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
	syncRepo := syncsqlc.NewSQLCRepository(postgres.DB)

	// 6. Services
	merchantSvc := merchant.NewService(merchantRepo)
	orderSvc := orders.NewService(orderRepo, orderItemRepo, categorySummaryRepo)
	employeeSvc := employees.NewService(employeeRepo, orderSvc)
	productSvc := products.NewService(productRepo, categoryRepo, analyticCategoryRepo, mappingRepo)
	analyticsSvc := analytics.NewService(analyticsRepo)
	dashboardSvc := dashboard.NewService(analyticsSvc, merchantSvc)

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
		cfg.Security.JWTSecret,
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

	cronScheduler := scheduler.NewScheduler(syncEngine, merchantRepo, syncRepo, log)
	cronScheduler.Start(context.Background())

	// 10. Handlers
	merchantHandler := merchant.NewHandler(merchantSvc)
	employeeHandler := employees.NewHandler(employeeSvc)
	productHandler := products.NewHandler(productSvc)
	orderHandler := orders.NewHandler(orderSvc)
	dashboardHandler := dashboard.NewHandler(dashboardSvc)
	syncHandler := sync.NewHandler(syncRepo, syncRepo, syncEngine)

	// 11. HTTP Server (centralized setup)
	server := api.NewServer(
		postgres,
		authHandler,
		merchantHandler,
		employeeHandler,
		productHandler,
		orderHandler,
		dashboardHandler,
		syncHandler,
		cfg.Server.FrontendURL,
		cfg.Security.JWTSecret,
	)

	go func() {
		addr := ":" + cfg.Server.Port
		log.Info("starting server", slog.String("addr", addr))
		if err := server.Start(addr); err != nil {
			log.Error("server error", slog.String("error", err.Error()))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server")

	cronScheduler.Stop()

	if cfg.AppEnv == "development" {
		if err := server.Close(); err != nil {
			log.Error("server close error", slog.String("error", err.Error()))
		}
	} else {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Error("server shutdown error", slog.String("error", err.Error()))
		}
	}

	log.Info("server stopped")
}
