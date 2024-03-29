package controller

import "github.com/gofiber/fiber/v2"

func (controller *Controller) Register(app *fiber.App) {
	core := app.Group("/api")
	core.Post("/publish", controller.Publish)
	core.Post("/subscribe", controller.Subscribe)
	core.Get("/subscribe/:type", controller.GetSubscribers)
	core.Get("/:type", controller.GetEventsByType)
	core.Delete("/:type/:host", controller.DeleteSubscriber)
	core.Get("/", controller.ListEventTypes)
}
