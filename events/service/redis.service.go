package service

import (
	"context"
	"encoding/json"
	"fmt"

	"example/event-bus/model"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// Subscribe(eventType string, host string) error
// PublishEvent(eventType string, event *model.UnknownEvent) error
// GetSubscribers(eventType string) ([]string, error)
// GetEventsByType(eventType string) ([]*model.UnknownEvent, error)
// DeleteSubscriber(eventType string, host string) error
// ListEventTypes() ([]string, error)

type RedisClient struct {
	rdb *redis.Client
	ctx context.Context
}

func NewRedisClient(ctx context.Context) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: "cache:6379",
		DB:   0,
	})
	return &RedisClient{rdb: rdb, ctx: ctx}
}

func (r *RedisClient) Subscribe(eventType string, host string) error {
	err := r.rdb.SAdd(r.ctx, eventType+":sub", host).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisClient) PublishEvent(eventType string, event *model.UnknownEvent) error {
	val, err := json.Marshal(event)
	fmt.Printf("val: %s\n", string(val))
	if err != nil {
		return err
	}

	err = r.rdb.LPush(r.ctx, eventType+":pub", string(val)).Err()

	if err != nil {
		return err
	}

	var subs []string
	subs, err = r.GetSubscribers(eventType)
	if err != nil {
		return err
	}

	for _, s := range subs {
		go func(s string) {
			agent := fiber.Post(s + "/events")
			agent.JSON(event)
			_, _, errs := agent.Bytes()
			if len(errs) > 0 {
				fmt.Println("Error publishing event to: " + s + " " + errs[0].Error())
			}
		}(s)
	}

	return nil
}

func (r *RedisClient) GetSubscribers(eventType string) ([]string, error) {
	res := r.rdb.SMembers(r.ctx, eventType+":sub")
	return res.Result()
}

func (r *RedisClient) GetEventsByType(eventType string) ([]*model.UnknownEvent, error) {
	res := r.rdb.LRange(r.ctx, eventType+":pub", 0, -1)
	vals, err := res.Result()
	if err != nil {
		return nil, err
	}

	var events []*model.UnknownEvent
	for _, val := range vals {
		var event model.UnknownEvent
		err = json.Unmarshal([]byte(val), &event)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, nil
}

func (r *RedisClient) DeleteSubscriber(eventType string, host string) error {
	err := r.rdb.SRem(r.ctx, eventType+":sub", host).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisClient) ListEventTypes() ([]string, error) {
	res := r.rdb.Keys(r.ctx, "*")
	return res.Val(), nil
}
