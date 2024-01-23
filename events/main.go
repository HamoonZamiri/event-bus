package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Event[T any] struct {
	Type string `json:"type"`
	Data T      `json:"data"`
}
type Event2 struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Comment struct {
	ID      string `json:"id"`
	PostID  string `json:"post_id"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

type SubscribeRequest struct {
	Host      string `json:"host"`
	EventType string `json:"event_type"`
}

type CommentEvent = Event[Comment]

var subscribers = make(map[string][]string)

func initializeSubscribers() {
	subscribers["comment_created"] = []string{}
	subscribers["post_created"] = []string{}
	subscribers["comment_moderated"] = []string{}
}

func Subscribe(c *fiber.Ctx) error {
	var subscribeRequest = new(SubscribeRequest)
	c.BodyParser(&subscribeRequest)

	if subscribeRequest.EventType == "" || subscribeRequest.Host == "" {
		return fiber.NewError(400, "Event Type and Host are required")
	}

	if subscribers[subscribeRequest.EventType] == nil {
		return fiber.NewError(400, "Event Type does not exist")
	}

	subscribers[subscribeRequest.EventType] = append(subscribers[subscribeRequest.EventType], subscribeRequest.Host)

	fmt.Println(subscribers)
	return nil
}

func findEventType(body []byte) string {
	var t string

	for i, c := range body {

		if c == 't' {
			j := i
			for body[j] != ' ' {
				j++
			}
			j += 2
			for body[j] != '"' {
				t += string(body[j])
				j++
			}
			break
		}
	}
	return t
}

func FindEventDomain(event string) string {
	return strings.Split(event, "_")[0]
}

func handleEvent(eventType string, event Event2) []error {
	currErrors := make([]error, 0)
	if subscribers[eventType] == nil {
		return append(currErrors, errors.New("event type does not exist"))
	}

	for _, s := range subscribers[eventType] {
		agent := fiber.Post(s + "/events")
		agent.JSON(event)

		_, _, errs := agent.Bytes()
		if len(errs) > 0 {
			message := fmt.Errorf("error sending event to: %s", s)
			currErrors = append(currErrors, message)
		}

	}

	return currErrors
}

func postEvent(c *fiber.Ctx) error {
	body := c.Body()

	event := new(Event2)
	err := c.BodyParser(&event)
	if err != nil {
		return err
	}

	eventType := event.Type

	errs := handleEvent(eventType, *event)
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "Event sent successfully",
		"data":    body,
		"receivers": subscribers[eventType],
	})

}

func main() {
	initializeSubscribers()
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON("Hello World")
	})
	app.Post("/publish", postEvent)
	app.Post("/subscribe", Subscribe)

	app.Listen(":8080")
}
