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
	g.POST("/sync/backfill", h.Backfill)
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

// BackfillRequest holds parameters for a manual backfill operation.
type BackfillRequest struct {
	RestaurantID uuid.UUID `json:"restaurant_id"`
	Entity       string    `json:"entity"`           // orders, payments
	From         string    `json:"from"`             // YYYY-MM-DD
	To           string    `json:"to,omitempty"`     // YYYY-MM-DD (optional, for logging)
}

// Backfill performs a manual historical sync for a specific date range.
func (h *Handler) Backfill(c echo.Context) error {
	var req BackfillRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if req.RestaurantID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, "restaurant_id is required")
	}
	if req.Entity == "" || req.From == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "entity and from are required")
	}

	cursor, err := time.Parse("2006-01-02", req.From)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid from date, expected YYYY-MM-DD")
	}

	ctx := c.Request().Context()
	var newCursor time.Time
	var count int
	var syncErr error

	switch req.Entity {
	case "orders":
		newCursor, count, syncErr = h.engine.SyncOrders(ctx, req.RestaurantID, cursor)
	case "payments":
		newCursor, count, syncErr = h.engine.SyncPayments(ctx, req.RestaurantID, cursor)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "invalid entity: must be 'orders' or 'payments'")
	}

	status := "success"
	var details []byte
	if syncErr != nil {
		status = "failed"
		details = []byte(`{"error": "` + syncErr.Error() + `"}`)
	}

	log := &Log{
		RestaurantID:   req.RestaurantID,
		Entity:           SyncEntity(req.Entity),
		Status:           status,
		RecordsProcessed: int32(count),
		CursorFrom:       &cursor,
		TriggeredBy:      "backfill",
	}
	if !newCursor.IsZero() {
		log.CursorTo = &newCursor
	}
	if req.To != "" {
		toDate, _ := time.Parse("2006-01-02", req.To)
		if !toDate.IsZero() {
			log.Details = []byte(fmt.Sprintf(`{"from": "%s", "to": "%s"}`, req.From, req.To))
		}
	}
	if syncErr != nil {
		log.Details = details
	}
	_ = h.logRepo.CreateLog(ctx, log)

	if syncErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, syncErr.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "backfill completed",
		"entity":       req.Entity,
		"from":         req.From,
		"to":           req.To,
		"count":        count,
		"cursor_from":  cursor.Format(time.RFC3339),
		"cursor_to":    newCursor.Format(time.RFC3339),
	})
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
