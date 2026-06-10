package dashboard

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Handler exposes dashboard HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a new dashboard HTTP handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers dashboard routes on the given Echo group.
func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.GET("/dashboard/summary", h.Summary)
	g.GET("/dashboard/sales-by-employee", h.SalesByEmployee)
	g.GET("/dashboard/top-products", h.TopProducts)
	g.GET("/dashboard/category-coverage", h.CategoryCoverage)
	g.GET("/dashboard/ticket-ideal", h.TicketIdeal)
}

func getRestaurantIDFromContext(c echo.Context) (uuid.UUID, error) {
	// Try JWT context first
	if merchantID, ok := c.Get("merchant_id").(string); ok && merchantID != "" {
		return uuid.Parse(merchantID)
	}
	// Fallback to query param (for backward compatibility or public routes)
	restaurantIDStr := c.QueryParam("restaurant_id")
	if restaurantIDStr == "" {
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "restaurant_id is required")
	}
	return uuid.Parse(restaurantIDStr)
}

func parseDateRange(c echo.Context) (uuid.UUID, time.Time, time.Time, error) {
	restaurantID, err := getRestaurantIDFromContext(c)
	if err != nil {
		return uuid.Nil, time.Time{}, time.Time{}, err
	}

	fromStr := c.QueryParam("from")
	toStr := c.QueryParam("to")
	if fromStr == "" || toStr == "" {
		return uuid.Nil, time.Time{}, time.Time{}, echo.NewHTTPError(http.StatusBadRequest, "from and to dates are required (YYYY-MM-DD)")
	}

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		return uuid.Nil, time.Time{}, time.Time{}, echo.NewHTTPError(http.StatusBadRequest, "invalid from date format (YYYY-MM-DD)")
	}
	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		return uuid.Nil, time.Time{}, time.Time{}, echo.NewHTTPError(http.StatusBadRequest, "invalid to date format (YYYY-MM-DD)")
	}
	// Inclusive date range: from start of day to end of day
	to = to.Add(24*time.Hour - time.Nanosecond)

	return restaurantID, from, to, nil
}

func (h *Handler) Summary(c echo.Context) error {
	restaurantID, from, to, err := parseDateRange(c)
	if err != nil {
		return err
	}

	summary, err := h.service.GetSummary(c.Request().Context(), restaurantID, from, to)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, summary)
}

func (h *Handler) SalesByEmployee(c echo.Context) error {
	restaurantID, from, to, err := parseDateRange(c)
	if err != nil {
		return err
	}

	result, err := h.service.GetSalesByEmployee(c.Request().Context(), restaurantID, from, to)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}

func (h *Handler) TopProducts(c echo.Context) error {
	restaurantID, from, to, err := parseDateRange(c)
	if err != nil {
		return err
	}

	limit := int32(10)
	if l := c.QueryParam("limit"); l != "" {
		// Simple parsing; Echo's binder could also be used
		var parsed int
		if _, err := fmt.Sscanf(l, "%d", &parsed); err == nil && parsed > 0 {
			limit = int32(parsed)
		}
	}

	result, err := h.service.GetTopProducts(c.Request().Context(), restaurantID, from, to, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}

func (h *Handler) CategoryCoverage(c echo.Context) error {
	restaurantID, from, to, err := parseDateRange(c)
	if err != nil {
		return err
	}

	result, err := h.service.GetCategoryCoverage(c.Request().Context(), restaurantID, from, to)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}

func (h *Handler) TicketIdeal(c echo.Context) error {
	restaurantID, from, to, err := parseDateRange(c)
	if err != nil {
		return err
	}

	result, err := h.service.GetTicketIdeal(c.Request().Context(), restaurantID, from, to)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}
