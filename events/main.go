package main

import (
	"context"

	"example/event-bus/controller"
	"example/event-bus/router"
	"example/event-bus/service"

	"github.com/gofiber/fiber/v2"
)

var ctx = context.Background()

func main() {
	app := router.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	redisStore := service.NewRedisClient(ctx)
	// eventStore := service.NewEventStore([]string{"comment_created", "post_created", "comment_moderated"})
	c := controller.NewController(redisStore)
	c.Register(app)

	app.Listen(":8080")
}
