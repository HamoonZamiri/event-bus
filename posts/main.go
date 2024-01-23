package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
)

type UnknownEvent struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

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

type CommentEvent = Event[Comment]

type Post struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Comments []Comment `json:"comments"`
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

	p.Comments = []Comment{}
	posts[id] = *p

	res := Response [Post]{
		Message: "Post created successfully",
		Data: *p,
	}
	publishEvent("post_created", p)
	return c.JSON(res)
}

func handleCommentModerated(c *fiber.Ctx) error {
	commentEvent := new(CommentEvent)
	c.BodyParser(commentEvent)

	postId := commentEvent.Data.PostID
	post, ok := posts[postId]

	if !ok {
		return fiber.NewError(404, "Post not found")
	}

	post.Comments = append(post.Comments, commentEvent.Data)
	posts[postId] = post

	return nil
}

func handleEvent(c *fiber.Ctx) error {
	unknownEvent := new(UnknownEvent)
	if err := c.BodyParser(unknownEvent); err != nil {
		return err
	}

	switch unknownEvent.Type {
		case "comment_moderated":
			if err := handleCommentModerated(c); err != nil {
				fmt.Println(err)
			}
		default:
			return nil
	}
	return nil
}

func subscribeToEvents(events []string) error {
	for _, event := range events {
		agent := fiber.Post("http://localhost:8080/subscribe")
		body := fiber.Map{
			"host": "http://localhost:8081",
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

func main() {
	// subscribe to events
	err := subscribeToEvents([]string{"post_created", "comment_moderated"})
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	// middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(posts)
	})

	app.Post("/", post)
	app.Post("/events", handleEvent)

	app.Listen(":8081")
}