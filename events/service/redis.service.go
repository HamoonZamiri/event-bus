package service

import (
	"context"
	"encoding/json"

	"example/event-bus/model"

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
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return &RedisClient{rdb: rdb, ctx: ctx}
}

func (r *RedisClient) Subscribe(eventType string, host string) error {
	err := r.rdb.SAdd(r.ctx, eventType, host).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisClient) PublishEvent(eventType string, event *model.UnknownEvent) error {
	val, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = r.rdb.LPush(r.ctx, eventType, val).Err()

	if err != nil {
		return err
	}

	return nil
}
