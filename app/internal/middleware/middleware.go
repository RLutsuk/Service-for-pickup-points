package middleware

import (
	"fmt"
	"net/http"
	"strings"

	userUC "github.com/RLutsuk/Service-for-pickup-points/app/internal/user/usecase"
	"github.com/RLutsuk/Service-for-pickup-points/app/models"
	monitoring "github.com/RLutsuk/Service-for-pickup-points/app/monitoring"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

func AuthWithRole(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, models.ErrUnAuthorizated.Error())
			}

			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, models.ErrBadAuthorizated.Error())
			}

			if len(headerParts[1]) == 0 {
				return echo.NewHTTPError(http.StatusUnauthorized, models.ErrBadAuthorizated.Error())
			}

			userID, userRole, err := userUC.Parsetoken(headerParts[1])

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			isAllowed := false
			for _, r := range allowedRoles {
				if r == userRole {
					isAllowed = true
					break
				}
			}
			if !isAllowed {
				return echo.NewHTTPError(http.StatusForbidden, models.ErrAccessDenied.Error())
			}

			c.Set("userRole", userRole)
			c.Set("userID", userID)

			return next(c)
		}
	}
}

func PrometheusMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Path()
		method := c.Request().Method

		timer := prometheus.NewTimer(monitoring.HTTPRequestDuration.WithLabelValues(path))
		defer timer.ObserveDuration()

		err := next(c)

		status := fmt.Sprintf("%d", c.Response().Status)
		monitoring.HTTPRequestsTotal.WithLabelValues(path, method, status).Inc()

		return err
	}
}
