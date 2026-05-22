package employees

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Handler exposes employee HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a new employee HTTP handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers employee routes on the given Echo group.
func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.GET("/employees", h.List)
	g.GET("/employees/:id", h.GetByID)
	// TODO: Add GET /employees/:id/orders once orders.Service is available (Phase 4.3)
	// g.GET("/employees/:id/orders", h.GetEmployeeOrders)
}

func (h *Handler) List(c echo.Context) error {
	// In a real multi-tenant setup the restaurantID would come from the JWT.
	// For now we require it as a query parameter for flexibility during development.
	restaurantIDStr := c.QueryParam("restaurant_id")
	if restaurantIDStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "restaurant_id is required")
	}
	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid restaurant_id")
	}

	emps, err := h.service.ListByRestaurant(c.Request().Context(), restaurantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, emps)
}

func (h *Handler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	emp, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if emp == nil {
		return echo.NewHTTPError(http.StatusNotFound, "employee not found")
	}
	return c.JSON(http.StatusOK, emp)
}
