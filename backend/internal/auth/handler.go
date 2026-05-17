package auth

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sales-insight/backend/internal/clover"
	"github.com/sales-insight/backend/internal/merchant"
)

// Handler exposes auth HTTP endpoints.
type Handler struct {
	service      *Service
	merchantSvc  *merchant.Service
	oauthClient  *clover.OAuthClient
	clientID     string
	frontendURL  string
}

// NewHandler creates a new auth HTTP handler.
func NewHandler(service *Service, merchantSvc *merchant.Service, oauthClient *clover.OAuthClient, clientID, frontendURL string) *Handler {
	return &Handler{
		service:     service,
		merchantSvc: merchantSvc,
		oauthClient: oauthClient,
		clientID:    clientID,
		frontendURL: frontendURL,
	}
}

// RegisterRoutes registers auth routes on the given Echo group.
func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.GET("/auth/clover", h.InitOAuth)
	g.GET("/auth/clover/callback", h.Callback)
	g.POST("/auth/bootstrap", h.Bootstrap)
	g.GET("/auth/status", h.Status)
	g.POST("/auth/revoke", h.Revoke)
}

// InitOAuth redirects the merchant to Clover's OAuth authorization page.
func (h *Handler) InitOAuth(c echo.Context) error {
	url := clover.AuthorizeURL("sandbox", h.clientID, h.frontendURL+"/auth/clover/callback")
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// Callback handles the OAuth callback from Clover.
func (h *Handler) Callback(c echo.Context) error {
	code := c.QueryParam("code")
	merchantID := c.QueryParam("merchant_id")

	if code == "" || merchantID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing code or merchant_id")
	}

	m, err := h.service.HandleCallback(c.Request().Context(), code, merchantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "connected",
		"merchant_id":  m.ID,
		"clover_env":   m.CloverEnv,
	})
}

// BootstrapRequest holds tokens for manual insertion (local dev).
type BootstrapRequest struct {
	Name             string `json:"name"`
	CloverMerchantID string `json:"clover_merchant_id"`
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
}

// Bootstrap manually inserts tokens for local development.
func (h *Handler) Bootstrap(c echo.Context) error {
	var req BootstrapRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if req.CloverMerchantID == "" || req.AccessToken == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "clover_merchant_id and access_token are required")
	}

	m, err := h.merchantSvc.Bootstrap(
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

// Status returns the connection status of the authenticated merchant.
func (h *Handler) Status(c echo.Context) error {
	// For now, list all merchants (in production this would use JWT session)
	merchants, err := h.merchantSvc.List(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"merchants": merchants,
	})
}

// RevokeRequest holds the merchant ID to revoke.
type RevokeRequest struct {
	MerchantID string `json:"merchant_id"`
}

// Revoke disconnects a merchant from Clover.
func (h *Handler) Revoke(c echo.Context) error {
	var req RevokeRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	id, err := uuid.Parse(req.MerchantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid merchant_id")
	}

	if err := h.merchantSvc.Revoke(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
