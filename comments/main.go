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

var comments = make(map[string]Comment)

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

	return c.JSON(comment)
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(comments)
	})

	app.Post("/", post)

	app.Listen(":8083")
}