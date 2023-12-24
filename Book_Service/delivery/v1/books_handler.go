package v1

import (
	"Book_Service/entity"
	"Book_Service/lib/logger"
	"Book_Service/service"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type BooksHandler struct {
	Service     *service.Service
	Logger      logger.Logger
	RedisClient *redis.Client
}

var ctx = context.Background()
var redisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

func NewBooksHandler(booksService *service.Service, logger logger.Logger, red *redis.Client) *BooksHandler {
	return &BooksHandler{
		Service:     booksService,
		Logger:      logger,
		RedisClient: red,
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
	cacheKey := "all-books"
	cachedData, err := h.RedisClient.Get(context.Background(), cacheKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return c.JSON(http.StatusInternalServerError, &entity.Response{
			Success: false,
			Message: "Error retrieving data from the cache",
		})
	}
	if errors.Is(err, redis.Nil) {
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
		go func() {

			response := entity.ListAllBooksResponse{
				Total: count,
				Page:  filter.Page,
				Books: res,
			}
			responseJSON, err := json.Marshal(response)
			if err != nil {
				return
			}
			if err := h.RedisClient.Set(context.Background(), cacheKey, responseJSON, time.Minute).Err(); err != nil {
				return
			}
		}()

		return c.JSON(http.StatusOK, &entity.Response{
			Success: true,
			Message: "Successfully get all books",
			Data:    res,
		})

	}
	var response entity.ListAllBooksResponse
	if err := json.Unmarshal([]byte(cachedData), &response); err != nil {
		return c.JSON(http.StatusInternalServerError, entity.Response{
			Success: false,
			Message: "error unmarshalling data",
		})
	}
	return c.JSON(http.StatusOK, &entity.Response{
		Success: true,
		Message: "getting from cache",
		Data:    response,
	})

}

func (h *BooksHandler) GetABook(c echo.Context) error {
	bookId := c.Param("id")
	//	fmt.Println("...................................", bookId)

	res, err := h.Service.GetABook(c.Request().Context(), bookId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &entity.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	message := res.UserId
	if err := redisClient.Publish(ctx, "request", message).Err(); err != nil {
		return err
	}

	cx := context.Background()
	subscriber := redisClient.Subscribe(cx, "response")
	msg, err := subscriber.ReceiveMessage(cx)
	if err != nil {
		fmt.Println("Error receiving message")
		return nil
	}
	fmt.Println("received message......", msg.Payload)
	var resp entity.RedisResponse
	resp.Title = res.Title
	resp.Id = res.Id
	resp.UserName = msg.Payload
	resp.PublicationYear = res.PublicationYear

	return c.JSON(http.StatusOK, &entity.Response{
		Success: true,
		Message: "Successfully get a book",
		Data:    resp,
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
