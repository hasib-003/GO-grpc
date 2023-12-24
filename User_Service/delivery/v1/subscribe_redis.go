package v1

import (
	"User_Service/repository"
	"User_Service/service"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/uptrace/bun"
	"strconv"
)

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
