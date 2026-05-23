package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/sales-insight/backend/internal/auth"
	"github.com/sales-insight/backend/internal/dashboard"
	"github.com/sales-insight/backend/internal/merchant"
	"github.com/sales-insight/backend/internal/orders"
	"github.com/sales-insight/backend/internal/platform/db"
	"github.com/sales-insight/backend/internal/products"
	"github.com/sales-insight/backend/internal/employees"
	"github.com/sales-insight/backend/internal/sync"
)

// Server holds the Echo instance and all handlers.
type Server struct {
	e          *echo.Echo
	postgres   *db.Postgres
}

// NewServer creates and configures a new Echo server.
func NewServer(
	postgres *db.Postgres,
	authHandler *auth.Handler,
	merchantHandler *merchant.Handler,
	employeeHandler *employees.Handler,
	productHandler *products.Handler,
	orderHandler *orders.Handler,
	dashboardHandler *dashboard.Handler,
	syncHandler *sync.Handler,
	frontendURL string,
) *Server {
	e := echo.New()
	e.HideBanner = true

	// Middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} ${uri} ${status} ${latency_human}\n",
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{frontendURL},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
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

	merchantHandler.RegisterRoutes(v1)
	authHandler.RegisterRoutes(v1)
	employeeHandler.RegisterRoutes(v1)
	productHandler.RegisterRoutes(v1)
	orderHandler.RegisterRoutes(v1)
	dashboardHandler.RegisterRoutes(v1)
	syncHandler.RegisterRoutes(v1)

	return &Server{
		e:        e,
		postgres: postgres,
	}
}

// Start begins listening on the given address.
func (s *Server) Start(addr string) error {
	return s.e.Start(addr)
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}

// Close immediately closes the server.
func (s *Server) Close() error {
	return s.e.Close()
}
