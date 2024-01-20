package main

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Event[T any] struct {
	Type string `json:"type"`
	Data T `json:"data"`
}
type Comment struct {
	ID string `json:"id"`
	PostID string `json:"post_id"`
	Content string `json:"content"`
	Status string `json:"status"`
}

type CommentEvent = Event[Comment]

func FindEventType(body []byte) string {
	var t string

	for i, c := range body {

		if c == 't' {
			j := i
			for body[j] != ' ' {
				j++
			}
			j+=2
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

func post(c *fiber.Ctx) error {
	fmt.Println(c.Body())
	body := c.Body()
	t := FindEventType(body)
	fmt.Println("\n"+t)
	fmt.Println(t == "comment_created")
	return nil
}

func main() {
	app := fiber.New()

	app.Post("/", post)
	app.Listen(":8080")
}