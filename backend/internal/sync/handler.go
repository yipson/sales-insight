package sync

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Handler exposes sync status HTTP endpoints.
type Handler struct {
	logRepo    LogRepository
	errorRepo  ErrorRepository
	engine     *Engine
}

// NewHandler creates a new sync HTTP handler.
func NewHandler(logRepo LogRepository, errorRepo ErrorRepository, engine *Engine) *Handler {
	return &Handler{
		logRepo:   logRepo,
		errorRepo: errorRepo,
		engine:    engine,
	}
}

// RegisterRoutes registers sync routes on the given Echo group.
func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.GET("/sync/status", h.Status)
	g.GET("/sync/logs", h.Logs)
	g.POST("/sync/trigger", h.Trigger)
	g.GET("/sync/errors", h.Errors)
}

func (h *Handler) Status(c echo.Context) error {
	restaurantID, err := parseRestaurantID(c)
	if err != nil {
		return err
	}

	logs, err := h.logRepo.GetLatestByEntity(c.Request().Context(), restaurantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, logs)
}

func (h *Handler) Logs(c echo.Context) error {
	restaurantID, err := parseRestaurantID(c)
	if err != nil {
		return err
	}

	limit := int32(50)
	if l := c.QueryParam("limit"); l != "" {
		var parsed int
		if _, err := fmt.Sscanf(l, "%d", &parsed); err == nil && parsed > 0 {
			limit = int32(parsed)
		}
	}

	logs, err := h.logRepo.ListByRestaurant(c.Request().Context(), restaurantID, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, logs)
}

func (h *Handler) Trigger(c echo.Context) error {
	var req struct {
		RestaurantID uuid.UUID `json:"restaurant_id"`
		Entity       string    `json:"entity"` // orders, items, employees, payments
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	var err error
	switch req.Entity {
	case "orders":
		_, _, err = h.engine.SyncOrders(ctx, req.RestaurantID, time.Time{})
	case "items":
		_, _, err = h.engine.SyncItems(ctx, req.RestaurantID)
	case "employees":
		_, err = h.engine.SyncEmployees(ctx, req.RestaurantID)
	case "payments":
		_, _, err = h.engine.SyncPayments(ctx, req.RestaurantID, time.Time{})
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "invalid entity")
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusAccepted)
}

func (h *Handler) Errors(c echo.Context) error {
	restaurantID, err := parseRestaurantID(c)
	if err != nil {
		return err
	}

	limit := int32(50)
	if l := c.QueryParam("limit"); l != "" {
		var parsed int
		if _, err := fmt.Sscanf(l, "%d", &parsed); err == nil && parsed > 0 {
			limit = int32(parsed)
		}
	}

	errors, err := h.errorRepo.ListUnresolved(c.Request().Context(), restaurantID, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, errors)
}

func parseRestaurantID(c echo.Context) (uuid.UUID, error) {
	restaurantIDStr := c.QueryParam("restaurant_id")
	if restaurantIDStr == "" {
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "restaurant_id is required")
	}
	return uuid.Parse(restaurantIDStr)
}
