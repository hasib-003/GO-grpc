package v1

import (
	"Book_Service/entity"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func sessionData(c echo.Context) *entity.JwtClaim {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*entity.JwtClaim)
	return claims

}
