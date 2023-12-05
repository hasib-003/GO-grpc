package v1

import (
	"User_Service/config"
	"User_Service/lib/logger"
	"User_Service/repository"
	"User_Service/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/uptrace/bun"
)

// Setup all routers
func SetupRouters(c *echo.Echo, conf *config.Config, db *bun.DB, jwtConfig middleware.JWTConfig, logger logger.Logger) {

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := NewUserHandler(userService, logger)

	authenticated := middleware.JWTWithConfig(jwtConfig)

	v1 := c.Group("/api/v1")

	userGroup := v1.Group("/user_registration")
	userHandler.MapUserRoutes(userGroup, authenticated)

}
