package controller

import (
	"example/event-bus/model"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type SubscribeRequest struct {
	Host      string `json:"host"`
	EventType string `json:"event_type"`
}

func (controller *Controller) Subscribe(c *fiber.Ctx) error {
	var req SubscribeRequest
	err := c.BodyParser(&req)

	if err != nil {
		return fiber.NewError(400, "Invalid request body")
	}

	if req.EventType == "" || req.Host == "" {
		return fiber.NewError(400, "Event Type and Host are required")
	}

	eStore := controller.eventStore
	if eStore.Subscribers[req.EventType] == nil {
		return fiber.NewError(400, "Event Type does not exist")
	}

	eStore.Subscribers[req.EventType] = append(eStore.Subscribers[req.EventType], req.Host)
	controller.eventStore = eStore
	return c.Status(200).JSON("Subscribed successfully")
}

func (controller *Controller) Publish(c *fiber.Ctx) error {
	unknownEvent := model.UnknownEvent{}
	err := c.BodyParser(&unknownEvent)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if unknownEvent.Type == "" || unknownEvent.Data == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Event Type and Data are required")
	}

	controller.eventStore.PublishEvent(unknownEvent.Type, unknownEvent)
	if controller.eventStore.Subscribers[unknownEvent.Type] == nil {
		return fiber.NewError(400, "Event Type does not exist")
	}

	for _, s := range controller.eventStore.Subscribers[unknownEvent.Type] {
		go func(s string) {
			agent := fiber.Post(s + "/events")
			agent.JSON(unknownEvent)
			_, _, errs := agent.Bytes()
			if len(errs) > 0 {
				fmt.Println("Error publishing event to: " + s)
			}
		}(s)
	}

	return c.Status(fiber.StatusOK).JSON("Event published successfully")
}
