package controller

import (
	"example/event-bus/model"

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

	err = controller.EventStore.Subscribe(req.EventType, req.Host)

	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	return c.Status(200).JSON("Subscribed successfully")
}

func (controller *Controller) Publish(c *fiber.Ctx) error {
	unknownEvent := new(model.UnknownEvent)
	es := controller.EventStore

	err := c.BodyParser(&unknownEvent)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if unknownEvent.Type == "" || unknownEvent.Data == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Event Type and Data are required")
	}

	err = es.PublishEvent(unknownEvent.Type, unknownEvent)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON("Event published successfully")
}

func (controller *Controller) GetSubscribers(c *fiber.Ctx) error {
	params := c.AllParams()

	event_type := params["type"]
	if event_type == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Event Type is required")
	}

	subscribers, err := controller.EventStore.GetSubscribers(event_type)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Subscribers retrieved successfully",
		"data":    subscribers,
	})
}

func (controller *Controller) GetEventsByType(c *fiber.Ctx) error {
	params := c.AllParams()

	event_type := params["type"]
	if event_type == "" {
		return fiber.NewError(400, "Event type is required")
	}

	events, err := controller.EventStore.GetEventsByType(event_type)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Events retrieved successfully",
		"data":    events,
	})
}

func (controller *Controller) DeleteSubscriber(c *fiber.Ctx) error {
	params := c.AllParams()
	eventType := params["type"]
	host := params["host"]

	if eventType == "" || host == "" {
		return fiber.NewError(400, "Event Type and Host are required")
	}

	err := controller.EventStore.DeleteSubscriber(eventType, host)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	return c.Status(fiber.StatusOK).SendString("Subscriber successfully deleted")
}

func (controller *Controller) ListEventTypes(c *fiber.Ctx) error {
	eventTypes, err := controller.EventStore.ListEventTypes()
	if err != nil {
		return fiber.NewError(400, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Event types retrieved successfully",
		"data":    eventTypes,
	})
}
