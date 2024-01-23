package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

type Post struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Comments []Comment `json:"comments"`
}

type CommentEvent = Event[Comment]
type PostEvent = Event[Post]

var comments = make(map[string]Comment)

func subscribeToEvents(events []string) error {
	for _, event := range events {
		agent := fiber.Post("http://localhost:8080/subscribe")
		body := fiber.Map{
			"host": "http://localhost:8083",
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

func handleEvent(c *fiber.Ctx) error {
	unknownEvent := new(UnknownEvent)
	if err := c.BodyParser(unknownEvent); err != nil {
		return err
	}

	switch unknownEvent.Type {
		case "comment_moderated":
			commentEvent := new(CommentEvent)
			if err := c.BodyParser(commentEvent); err != nil {
				return err
			}
			id := commentEvent.Data.ID
			if comments[id] == (Comment{}) {
				return fiber.NewError(404, "Moderated comment was not found!")
			}
			comments[id] = commentEvent.Data
		case "post_created":
			postEvent := new(PostEvent)
			if err := c.BodyParser(postEvent); err != nil {
				return err
			}
		default:
			return nil
	}
	return nil
}

func publishEvent(eventType string, data interface{}) error {
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

func post(c *fiber.Ctx) error {
	var commentID = uuid.New().String()
	comment := new(Comment)
	if err := c.BodyParser(comment); err != nil {
		return err
	}

	if comment.Content == "" || comment.PostID == "" {
		return c.Status(400).JSON("Content and PostID are required")
	}

	comment.ID = commentID

	comment.Status = "pending"
	comments[commentID] = *comment

	// Publish an event to the event bus
	if err := publishEvent("comment_created", comment); err != nil {
		return err
	}

	return c.JSON(comment)
}

func main() {
	subscribeToEvents([]string{"comment_moderated", "post_created"})
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(comments)
	})

	app.Post("/", post)
	app.Post("/events", handleEvent)

	app.Listen(":8083")
}