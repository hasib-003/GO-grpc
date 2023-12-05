package v1

import (
	"User_Service/entity"
	"User_Service/entity/httpentity"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func InvalidJwt(err error) error {
	if err != nil {
		log.Println("Got from echo")
	}
	return echo.NewHTTPError(http.StatusUnauthorized,
		&httpentity.ErrorResponse{
			Success:      false,
			ErrorCode:    "unauthorized",
			ErrorMessage: "You are not authorized to sthis resource",
		})
}

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func authorized(roles ...Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, &entity.Response{
					Success: false,
					Message: "JWT TOKEN MISSING",
				})
			}

			claim, ok := token.Claims.(*entity.JwtClaims)
			fmt.Printf("claim user id>>>>>>> %+v\n", claim)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, &entity.Response{
					Success: false,
					Message: "JWT CLAIMS MISSING",
				})
			}

			for _, role := range roles {
				fmt.Println("role.....................", role)
				fmt.Println("clainm >>>>>>>", claim.UserId)
				if checkPermission(role, claim.Role) {
					return next(c)
				}
			}
			return echo.NewHTTPError(http.StatusUnauthorized, &entity.Response{
				Success: false,
				Message: "You are not authorized to access this resource",
			})
		}
	}
}
func checkPermission(requiredRole Role, userRole string) bool {
	switch requiredRole {
	case RoleAdmin:
		return userRole == entity.ACCOUNT_TYPE_ADMIN
	case RoleUser:
		return userRole == entity.ACCOUNT_TYPE_ADMIN || userRole == entity.ACCOUNT_TYPE_USER

	default:
		return false
	}
}
