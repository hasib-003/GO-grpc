package v1

import (
	"User_Service/entity"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func sessionData(c echo.Context) *entity.JwtClaims {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*entity.JwtClaims)
	return claims
}
