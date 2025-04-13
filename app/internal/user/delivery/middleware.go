package delivery

import (
	"net/http"
	"strings"

	userUC "github.com/RLutsuk/Service-for-pickup-points/app/internal/user/usecase"
	"github.com/RLutsuk/Service-for-pickup-points/app/models"
	"github.com/labstack/echo"
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
				return echo.NewHTTPError(http.StatusUnauthorized, models.ErrBadAuthorizated)
			}

			if len(headerParts[1]) == 0 {
				return echo.NewHTTPError(http.StatusUnauthorized, models.ErrBadAuthorizated)
			}

			userID, userRole, err := userUC.Parsetoken(headerParts[1])

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}

			isAllowed := false
			for _, r := range allowedRoles {
				if r == userRole {
					isAllowed = true
					break
				}
			}
			if !isAllowed {
				return echo.NewHTTPError(http.StatusForbidden, models.ErrAccessDenied)
			}

			c.Set("userID", userID)
			c.Set("userRole", userRole)

			return next(c)
		}
	}
}
