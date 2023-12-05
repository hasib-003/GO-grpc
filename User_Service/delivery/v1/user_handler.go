package v1

import (
	"User_Service/entity"
	"User_Service/entity/httpentity"
	"User_Service/lib/logger"
	"User_Service/service"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type UserHandler struct {
	UserService *service.UserService
	Logger      logger.Logger
	RedisClient redis.Client
}

func NewUserHandler(userService *service.UserService, logger logger.Logger, red *redis.Client) *UserHandler {
	return &UserHandler{
		UserService: userService,
		Logger:      logger,
		RedisClient: red,
	}
}

func (h *UserHandler) MapUserRoutes(userGroup *echo.Group, authenticated echo.MiddlewareFunc) {
	userGroup.POST("", h.CreateUser)
	userGroup.GET("", h.GetAllUser, authenticated, authorized(RoleAdmin))
	userGroup.GET("/single", h.GetAUser, authenticated, authorized(RoleUser, RoleAdmin))
	userGroup.GET("/df/:email", h.GetUserByEmail, authenticated)
	userGroup.PUT("/:id", h.UpdateUser, authenticated, authorized(RoleUser, RoleAdmin))
	userGroup.DELETE("/delete", h.DeleteUser, authenticated, authorized(RoleUser, RoleAdmin))
	userGroup.POST("/login", h.Login)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	data := httpentity.CreateUserRegistration{}

	err := c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &entity.Response{
			Success: false,
			Message: "Invalid Payload",
			Data:    err,
		})
	}
	validationError := data.Validate()
	if validationError != nil {
		return c.JSON(http.StatusBadRequest, &entity.Response{
			Success: false,
			Message: "Validation Error",
			Data:    err.Error(),
		})
	}
	err = h.UserService.CreateUser(c.Request().Context(), data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &entity.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, &entity.Response{
		Success: true,
		Message: "New User Created Successfully",
	})
}

func (h *UserHandler) GetAllUser(c echo.Context) error {

	cacheKey := "all_users"
	cachedData, err := h.RedisClient.Get(context.Background(), cacheKey).Result()
	//if err == nil {
	//	// Data found in cache, return cached data
	//	var response entity.GetAllUserResponses
	//	if err := json.Unmarshal([]byte(cachedData), &response); err != nil {
	//		return c.JSON(http.StatusInternalServerError, &entity.Response{
	//			Success: false,
	//			Message: "Error unmarshalling cached data",
	//		})
	//	}
	//	return c.JSON(http.StatusOK, &entity.Response{
	//		Success: true,
	//		Message: "Getting done (from cache)",
	//		Data:    response,
	//	})
	//}
	if err != nil {
		return c.JSON(http.StatusNotFound, &entity.Response{
			Success: false,
			Message: "No data available to chache",
		})
	}
	var response entity.GetAllUserResponses
	if err := json.Unmarshal([]byte(cachedData), &response); err != nil {
		return c.JSON(http.StatusInternalServerError, &entity.Response{
			Success: false,
			Message: "error unmarshaling data",
		})
	}
	return c.JSON(http.StatusOK, &entity.Response{
		Success: true,
		Message: "getting done (from cache)",
		Data:    response,
	})

	filter := entity.UserFilter{}
	err = c.Bind(&filter)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &entity.Response{
			Success: false,
			Message: "Internal Server error",
		})
	}
	res, count, err := h.UserService.GetAllUser(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &entity.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	go func() {
		response := entity.GetAllUserResponses{
			Total: count,
			Page:  filter.Page,
			Users: res,
		}
		fmt.Println(response)
		responseJSON, err := json.Marshal(response)
		if err != nil {
			return
		}
		if err == nil {
			h.RedisClient.Set(context.Background(), cacheKey, responseJSON, time.Minute)
		}
	}()

	return c.JSON(http.StatusOK, &entity.Response{
		Success: true,
		Message: "Getting done",
		Data:    res,
	})
}

func (h *UserHandler) GetAUser(c echo.Context) error {
	id := sessionData(c).UserId
	fmt.Println("userid...........", id)
	res, err := h.UserService.GetAUser(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &entity.Response{
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
func (h *UserHandler) GetUserByEmail(c echo.Context) error {
	email := c.Param("email")
	res, err := h.UserService.GetUserByEmail(c.Request().Context(), email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &entity.Response{
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

func (h *UserHandler) UpdateUser(c echo.Context) error {
	id := sessionData(c).UserId
	data := entity.UserRegistration{}
	err := c.Bind(&data)
	if err != nil {
		err := c.JSON(http.StatusBadRequest, &entity.Response{
			Success: false,
			Message: "invalid request",
		})
		if err != nil {
			return err
		}
	}
	validationError := data.Validate()
	if validationError != nil {
		err := c.JSON(http.StatusBadRequest, &entity.Response{
			Success: false,
			Message: "Validation Error",
		})
		if err != nil {
			return err
		}
	}
	err = h.UserService.UpdateUser(c.Request().Context(), data, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &entity.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, &entity.Response{
		Success: true,
		Message: "successfully Updated",
	})
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := sessionData(c).UserId
	err := h.UserService.DeleteUser(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &entity.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, &entity.Response{
		Success: true,
		Message: "successfully deleted User",
	})
}

func (h *UserHandler) Login(c echo.Context) error {
	login := entity.LoginInfo{}
	err := c.Bind(&login)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &entity.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	user, err := h.UserService.GetUserByEmail(c.Request().Context(), login.Email)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &entity.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		return c.JSON(http.StatusBadRequest, &entity.Response{
			Success: false,
			Message: "Wrong password",
		})
	}
	expirationTime := time.Now().Add(24 * time.Hour)
	token, err := user.GetJwt(expirationTime, user.UserId, user.Role)
	fmt.Println("useid>>>.....userrole...........", user.UserId, user.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &entity.Response{
			Success: false,
			Message: "Error signing token",
		})
	}
	return c.JSON(http.StatusOK, &entity.Response{
		Success: true,
		Message: "All ok ,Good to go ",
		Data:    token,
	})
}
