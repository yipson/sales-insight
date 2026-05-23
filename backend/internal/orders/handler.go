package orders

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Handler exposes order HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a new order HTTP handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers order routes on the given Echo group.
func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.GET("/orders", h.List)
	g.GET("/orders/:id", h.GetByID)
}

func (h *Handler) List(c echo.Context) error {
	restaurantIDStr := c.QueryParam("restaurant_id")
	if restaurantIDStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "restaurant_id is required")
	}
	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid restaurant_id")
	}

	fromStr := c.QueryParam("from")
	toStr := c.QueryParam("to")
	if fromStr == "" || toStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "from and to dates are required")
	}
	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid from date format (RFC3339)")
	}
	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid to date format (RFC3339)")
	}

	orders, err := h.service.ListByRestaurantAndDate(c.Request().Context(), restaurantID, from, to)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, orders)
}

func (h *Handler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	detail, err := h.service.GetOrderWithDetails(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, detail)
}
