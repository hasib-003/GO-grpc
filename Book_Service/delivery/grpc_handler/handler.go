package grpc_handler

import (
	"Book_Service/entity"
	"Book_Service/lib/logger"
	"Book_Service/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type BooksHandler struct {
	Service *service.Service
	Logger  logger.Logger
}

func NewBooksHandler(booksService *service.Service, logger logger.Logger) *BooksHandler {
	return &BooksHandler{
		Service: booksService,
		Logger:  logger,
	}
}

func (h *BooksHandler) GetBooksByUserID(c echo.Context) error {
	userId := c.Param("user_id")
	res, err := h.Service.GetBooksByUserID(c.Request().Context(), entity.MessageRequest{
		UserID: userId,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &entity.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, &entity.Response{
		Success: true,
		Message: "Successfully get a user",
		Data:    res,
	})

}
