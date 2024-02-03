package controller

import "example/event-bus/service"

type Controller struct {
	EventStore service.EventService
}

func NewController(eventStore service.EventService) *Controller {
	return &Controller{EventStore: eventStore}
}
