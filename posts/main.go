package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)
type Post struct {
	ID string `json:"id"`
	Title string `json:"title"`
}
type Response[T any] struct {
	Message string `json:"message"`
	Data T `json:"data"`
}

var posts = make(map[string]Post)

func post(c *fiber.Ctx) error {
	id := uuid.New().String()
	p := new(Post)
	p.ID = id
	if err := c.BodyParser(p); err != nil {
		return err
	}

	if p.Title == "" {
		return c.Status(400).JSON(Response[string]{
			Message: "Title is required",
		})
	}
	posts[id] = *p
	res := Response [Post]{
		Message: "Post created successfully",
		Data: *p,
	}
	return c.JSON(res)
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(posts)
	})

	app.Post("/", post)

	app.Listen(":8081")
}