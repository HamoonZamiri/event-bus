package model

type UnknownEvent struct {
	Data any    `json:"data"`
	Type string `json:"type"`
}

type Event[T any] struct {
	Data T      `json:"data"`
	Type string `json:"type"`
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
