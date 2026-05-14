package merchant

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Handler exposes merchant HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a new merchant HTTP handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers merchant routes on the given Echo group.
func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.GET("/merchants", h.List)
	g.GET("/merchants/:id", h.GetByID)
	g.POST("/merchants/bootstrap", h.Bootstrap)
	g.POST("/merchants/:id/revoke", h.Revoke)
}

func (h *Handler) List(c echo.Context) error {
	merchants, err := h.service.List(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, merchants)
}

func (h *Handler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	m, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if m == nil {
		return echo.NewHTTPError(http.StatusNotFound, "merchant not found")
	}
	return c.JSON(http.StatusOK, m)
}

// BootstrapRequest holds tokens for manual insertion.
type BootstrapRequest struct {
	Name             string `json:"name"`
	CloverMerchantID string `json:"clover_merchant_id"`
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
}

func (h *Handler) Bootstrap(c echo.Context) error {
	var req BootstrapRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if req.CloverMerchantID == "" || req.AccessToken == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "clover_merchant_id and access_token are required")
	}

	m, err := h.service.Bootstrap(
		c.Request().Context(),
		req.Name,
		req.CloverMerchantID,
		req.AccessToken,
		req.RefreshToken,
		nil, nil,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, m)
}

func (h *Handler) Revoke(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	if err := h.service.Revoke(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
