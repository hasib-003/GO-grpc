package main

import (
	"User_Service/repository"
	"User_Service/service"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/uptrace/bun"
	"log"
	"strconv"
	"syscall"

	"net/http"
	"os"
	"os/signal"
	"time"

	"User_Service/config"
	"User_Service/config/database"
	v1 "User_Service/delivery/v1"
	"User_Service/entity"
	"User_Service/lib/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title ERP core Service API Documentation.
// @version 1.0
// @description This is a sample api documentation.

// @host localhost:1327
// @BasePath /api/v1

func main() {
	conf := config.NewConfig("config.env")
	appLogger := logger.NewApiLogger(conf)

	appLogger.InitLogger()
	appLogger.Info("Starting the API Server")
	db := database.NewDB(conf)

	e := echo.New()
	//e.Logger.SetLevel(log.INFO)
	// Enable HTTP compression
	e.Use(middleware.Gzip())

	// Recover from panics
	e.Use(middleware.Recover())

	// Allow requests from *
	e.Use(middleware.CORS())
	// Print http request and response log to stdout if debug is enabled
	if conf.Debug {
		e.Use(middleware.Logger())
	}

	// JWT Middleware
	jwtConfig := middleware.JWTConfig{
		Claims:       &entity.JwtClaims{},
		SigningKey:   []byte(conf.JwtSecret),
		ErrorHandler: v1.InvalidJwt,
	}

	v1.SetupRouters(e, conf, db, jwtConfig, appLogger)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	go SubscribeToRedis(db)
	go httpServer(e, conf.HTTP)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	log.Println("Shutting down HTTP server...")
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	log.Println("HTTP server stopped!")
}

func httpServer(e *echo.Echo, httpConfig config.HTTP) {
	if err := e.Start(httpConfig.HTTPAddress); err != nil && !errors.Is(err, http.ErrServerClosed) {
		e.Logger.Fatal("shutting down the server")
	}
}

func SubscribeToRedis(db *bun.DB) {
	var ctx = context.Background()
	var redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	subscriber := redisClient.Subscribe(ctx, "request")
	defer subscriber.Close()

	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			fmt.Println("Error receiving message:", err)
			continue
		}
		fmt.Printf("%T\n", msg.Payload)
		id, _ := strconv.Atoi(msg.Payload)
		fmt.Println(" ...........................", id)
		userRepository := repository.NewUserRepository(db)
		userService := service.NewUserService(userRepository)
		user, err := userService.GetAUserRedis(id)
		Name := user.FirstName + " " + user.LastName
		fmt.Println(Name)
		cx := context.Background()
		result := redisClient.Publish(cx, "response", Name)
		if result.Err() != nil {

			fmt.Println("Error publishing message:", result.Err())
		} else {
		
			fmt.Println("Message published successfully")
		}

	}
}
