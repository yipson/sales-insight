package employees

import (
	"net/http"
	"time"

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
	g.GET("/employees/:id/orders", h.GetEmployeeOrders)
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

func (h *Handler) GetEmployeeOrders(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
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

	orders, err := h.service.GetEmployeeOrders(c.Request().Context(), id, from, to)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, orders)
}
