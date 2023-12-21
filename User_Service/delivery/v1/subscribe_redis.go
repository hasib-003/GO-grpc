package v1

import (
	"User_Service/entity"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
)

var userHandler = &UserHandler{}

func SubscribeToRedis() {
	var ctx = context.Background()
	var redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	subscriber := redisClient.Subscribe(ctx, "publisher")
	defer subscriber.Close()
	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			fmt.Println("Error receiving message:", err)
			continue
		}
		fmt.Printf("%T\n", msg.Payload)
		id, _ := strconv.Atoi(msg.Payload)
		user, err := userHandler.GetAUserByID(id)
		fmt.Println(user)

	}
}

func (h *UserHandler) GetAUserByID(id int) (entity.UserRegistration, error) {
	fmt.Println("userid...........", id)

	// Call the UserService with the provided ID
	res, err := h.UserService.GetAUserRedis(id)
	if err != nil {
		return entity.UserRegistration{}, err
	}

	return res, nil
}
