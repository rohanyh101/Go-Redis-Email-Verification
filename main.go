package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
)

func Init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	redisClient = Connect2db()
}

func Connect2db() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		DB:   0,
	})
}

func main() {

	Init()

	app := fiber.New()

	// Don't use it on prod. You should start verification inside your signup route. Also don't use GET method for this.
	app.Get("/send-verification-email", sendEmail)

	// If you're on production, don't send this link in mail or don't use GET method for verification. Do this stuff on clientside.
	app.Get("/verify-email", verifyEmail)
	app.Listen(":3000")
}

func sendEmail(c *fiber.Ctx) error {
	email := c.Query("email")
	varificationLink, err := SendVerificationMail(email)
	if err != nil {
		fmt.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "dude it's not working...",
		})
	}

	fmt.Printf("Here's your link: %v", varificationLink)
	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Sent!",
	})
}

func verifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	email, err := redisClient.Get(ctx, token).Result()
	if err == redis.Nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Invalid credentials",
		})
	}

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "error while getting token",
		})
	}

	fmt.Println(email)
	fmt.Println("Now you can update emailVerified column to true")

	_, err = redisClient.Del(ctx, token).Result()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "error while deleting token",
		})
	}

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Your email has been verified!!!",
	})
}
