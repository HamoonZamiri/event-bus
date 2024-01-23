package controller

import "example/event-bus/service"

type Controller struct {
	eventStore service.EventService
}

func NewController(eventStore *service.EventStore) *Controller {
	return &Controller{eventStore: eventStore}
}
