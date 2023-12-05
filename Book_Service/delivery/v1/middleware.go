package v1

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"

	"Book_Service/entity/httpentity"
)

func InvalidJwt(err error) error {
	if err != nil {
		log.Println("Got from echo")
	}
	return echo.NewHTTPError(http.StatusUnauthorized,
		&httpentity.ErrorResponse{
			Success:      false,
			ErrorCode:    "unauthorized",
			ErrorMessage: "You are not authorized to access this resource",
		})
}
