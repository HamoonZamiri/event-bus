package main

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Comment struct {
	ID string `json:"id"`
	PostID string `json:"post_id"`
	Content string `json:"content"`
	Status string `json:"status"`
}

type Event[T any] struct {
	Type string `json:"type"`
	Data T      `json:"data"`
}

type UnknownEvent struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type CommentEvent = Event[Comment]

var comments = make([]Comment, 0)

func publishEvent(eventType string, data interface{}) error {
	// Publish an event to the event bus

	agent := fiber.Post("http://localhost:8080/publish")
	body := fiber.Map{
		"type": eventType,
		"data": data,
	}

	agent.JSON(body)
	_, _, errs := agent.Bytes()

	if errs != nil {
		return errs[0]
	}
	return nil;

}

func subscribeToEvents(events []string) error {
	for _, event := range events {
		agent := fiber.Post("http://localhost:8080/subscribe")
		body := fiber.Map{
			"host": "http://localhost:8082",
			"event_type": event,
		}

		agent.JSON(body)
		_, _, errs := agent.Bytes()

		if errs != nil {
			return errs[0]
		}
	}
	return nil;
}

func isValidComment(comment string) bool {
	if strings.Contains(comment, "orange") {
		return false
	}
	return true
}

func handleEvent(c *fiber.Ctx) error {
	unknownEvent := new(UnknownEvent)
	c.BodyParser(unknownEvent)

	if unknownEvent.Type == "comment_created" {
		commentEvent := new(CommentEvent)
		c.BodyParser(commentEvent)

		if isValidComment(commentEvent.Data.Content) {
			moderatedComment := commentEvent.Data
			moderatedComment.Status = "approved"

			comments = append(comments, moderatedComment)

			err := publishEvent("comment_moderated", moderatedComment)
			if err != nil {
				fmt.Println(err)
			}
		}
		return nil
	}
	return nil
}

func main() {
	subscribeToEvents([]string{"comment_created"})
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON("Hello World")
	})

	app.Post("/events", handleEvent)

	app.Listen(":8082")
}