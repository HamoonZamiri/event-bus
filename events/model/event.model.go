package model

type UnknownEvent struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Event[T any] struct {
	Type string `json:"type"`
	Data T      `json:"data"`
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