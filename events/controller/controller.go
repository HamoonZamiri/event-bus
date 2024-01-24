package controller

import "example/event-bus/service"

type Controller struct {
	EventStore service.EventService
}

func NewController(eventStore *service.EventStore) *Controller {
	return &Controller{EventStore: eventStore}
}
