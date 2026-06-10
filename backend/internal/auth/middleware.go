package auth

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWTAuth middleware validates JWT tokens on protected routes.
// It extracts merchant_id and clover_merchant_id from the token claims
// and stores them in the echo context for downstream handlers.
func JWTAuth(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			tokenString := strings.TrimPrefix(auth, "Bearer ")
			if tokenString == auth {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			// Extract claims and store in context for handlers
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if merchantID, ok := claims["merchant_id"].(string); ok {
					c.Set("merchant_id", merchantID)
				}
				if cloverMerchantID, ok := claims["clover_merchant_id"].(string); ok {
					c.Set("clover_merchant_id", cloverMerchantID)
				}
			}

			return next(c)
		}
	}
}
