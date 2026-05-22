package products

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Handler exposes product HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a new product HTTP handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers product routes on the given Echo group.
func (h *Handler) RegisterRoutes(g *echo.Group) {
	// Products
	g.GET("/products", h.ListProducts)
	g.GET("/products/:id", h.GetProductByID)
	// TODO: GET /products/top requires order_items analytics (Phase 4.4)
	// g.GET("/products/top", h.GetTopProducts)

	// Native categories
	g.GET("/products/categories", h.ListCategories)
	g.PUT("/products/categories/:id/map", h.MapCategory)

	// Analytic categories
	g.GET("/analytic-categories", h.ListAnalyticCategories)
	g.POST("/analytic-categories", h.CreateAnalyticCategory)
	g.PUT("/analytic-categories/:id", h.UpdateAnalyticCategory)
	g.DELETE("/analytic-categories/:id", h.DeleteAnalyticCategory)
}

// ============================================================
// Product handlers
// ============================================================

func (h *Handler) ListProducts(c echo.Context) error {
	restaurantIDStr := c.QueryParam("restaurant_id")
	if restaurantIDStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "restaurant_id is required")
	}
	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid restaurant_id")
	}

	products, err := h.service.ListProductsByRestaurant(c.Request().Context(), restaurantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, products)
}

func (h *Handler) GetProductByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	product, err := h.service.GetProductByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if product == nil {
		return echo.NewHTTPError(http.StatusNotFound, "product not found")
	}
	return c.JSON(http.StatusOK, product)
}

// ============================================================
// Category handlers
// ============================================================

func (h *Handler) ListCategories(c echo.Context) error {
	restaurantIDStr := c.QueryParam("restaurant_id")
	if restaurantIDStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "restaurant_id is required")
	}
	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid restaurant_id")
	}

	categories, err := h.service.ListCategoriesByRestaurant(c.Request().Context(), restaurantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, categories)
}

func (h *Handler) MapCategory(c echo.Context) error {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid category id")
	}

	var req struct {
		RestaurantID       uuid.UUID `json:"restaurant_id"`
		AnalyticCategoryID uuid.UUID `json:"analytic_category_id"`
		MappedBy           string    `json:"mapped_by"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	mapping := &CategoryMapping{
		ID:                 uuid.New(),
		RestaurantID:       req.RestaurantID,
		CategoryID:         categoryID,
		AnalyticCategoryID: req.AnalyticCategoryID,
		MappedBy:           req.MappedBy,
		MappedAt:           time.Now().UTC(),
		CreatedAt:          time.Now().UTC(),
		UpdatedAt:          time.Now().UTC(),
	}

	if err := h.service.UpsertCategoryMapping(c.Request().Context(), mapping); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, mapping)
}

// ============================================================
// Analytic Category handlers
// ============================================================

func (h *Handler) ListAnalyticCategories(c echo.Context) error {
	restaurantIDStr := c.QueryParam("restaurant_id")
	if restaurantIDStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "restaurant_id is required")
	}
	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid restaurant_id")
	}

	categories, err := h.service.ListAnalyticCategoriesByRestaurant(c.Request().Context(), restaurantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, categories)
}

func (h *Handler) CreateAnalyticCategory(c echo.Context) error {
	var req struct {
		RestaurantID uuid.UUID `json:"restaurant_id"`
		Name         string    `json:"name"`
		Slug         string    `json:"slug"`
		DisplayOrder int32     `json:"display_order"`
		Color        string    `json:"color"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ac, err := h.service.CreateAnalyticCategory(
		c.Request().Context(),
		req.RestaurantID,
		req.Name,
		req.Slug,
		req.DisplayOrder,
		req.Color,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, ac)
}

func (h *Handler) UpdateAnalyticCategory(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	var req struct {
		Name         string `json:"name"`
		Slug         string `json:"slug"`
		DisplayOrder int32  `json:"display_order"`
		Color        string `json:"color"`
		IsActive     bool   `json:"is_active"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ac, err := h.service.UpdateAnalyticCategory(
		c.Request().Context(),
		id,
		req.Name,
		req.Slug,
		req.DisplayOrder,
		req.Color,
		req.IsActive,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, ac)
}

func (h *Handler) DeleteAnalyticCategory(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	if err := h.service.DeleteAnalyticCategory(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
