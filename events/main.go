package main

import (
	"example/event-bus/controller"
	"example/event-bus/router"
	"example/event-bus/service"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := router.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	eventStore := service.NewEventStore([]string{"comment_created", "post_created", "comment_moderated"})
	c := controller.NewController(eventStore)
	c.Register(app)

	app.Listen(":8080")
}
