package v1

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

func (h *BooksHandler) MapBooksRoutes(booksGroup *echo.Group, authenticated echo.MiddlewareFunc) {
	booksGroup.POST("", h.Create)
	booksGroup.GET("", h.ListAllBooks)
	booksGroup.GET("/:id", h.GetABook)
	booksGroup.PUT("/:id", h.Update)
	booksGroup.DELETE("/:id", h.Delete)
}

func (h *BooksHandler) Create(c echo.Context) error {

	data := entity.Books{}
	err := c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &entity.Response{
			Success: false,
			Message: "Invalid request payload",
			Data:    err,
		})
	}
	validationErrors := data.Validate()
	if validationErrors != nil {
		return c.JSON(http.StatusBadRequest, &entity.Response{
			Success: false,
			Message: "Field validation error!",
			Data:    validationErrors,
		})
	}

	err = h.Service.Create(c.Request().Context(), data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &entity.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, &entity.Response{
		Success: true,
		Message: "New Book Created Successfully!",
	})
}

func (h *BooksHandler) ListAllBooks(c echo.Context) error {
	filter := entity.BooksFilter{}
	err := c.Bind(&filter)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &entity.Response{
			Success: false,
			Message: "Invalid request filter",
			Data:    err,
		})
	}

	res, count, err := h.Service.ListAllBooks(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &entity.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	response := entity.ListAllBooksResponse{
		Total: count,
		Page:  filter.Page,
		Books: res,
	}

	return c.JSON(http.StatusOK, &entity.Response{
		Success: true,
		Message: "Successfully get all books",
		Data:    response,
	})

}
func (h *BooksHandler) GetABook(c echo.Context) error {
	bookId := c.Param("id")
	res, err := h.Service.GetABook(c.Request().Context(), bookId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &entity.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &entity.Response{
		Success: true,
		Message: "Successfully get a book",
		Data:    res,
	})

}

func (h *BooksHandler) Update(c echo.Context) error {
	bookId := c.Param("id")
	data := entity.Books{}
	err := c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &entity.Response{
			Success: false,
			Message: "Invalid request payload",
			Data:    err,
		})
	}
	validationErrors := data.Validate()
	if validationErrors != nil {
		return c.JSON(http.StatusBadRequest, &entity.Response{
			Success: false,
			Message: "Field validation error!",
			Data:    validationErrors,
		})
	}

	err = h.Service.Update(c.Request().Context(), data, bookId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &entity.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &entity.Response{
		Success: true,
		Message: "Successfully Updated",
	})
}
func (h *BooksHandler) Delete(c echo.Context) error {
	bookId := c.Param("id")
	err := h.Service.Delete(c.Request().Context(), bookId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &entity.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &entity.Response{
		Success: true,
		Message: "Successfully Deleted",
	})
}
