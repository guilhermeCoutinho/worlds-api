package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/guilhermeCoutinho/worlds-api/models"
	"github.com/sirupsen/logrus"
)

type EventPublisher interface {
	PublishWorldCreated(ctx context.Context, world *models.World) error
	PublishWorldUpdated(ctx context.Context, world *models.World) error
}

type WorldEvent struct {
	Type      string      `json:"type"`
	WorldID   string      `json:"world_id"`
	UserID    string      `json:"user_id"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

type RedisEventPublisher struct {
	client *redis.Client
	logger logrus.FieldLogger
}

func NewRedisEventPublisher(client *redis.Client, logger logrus.FieldLogger) *RedisEventPublisher {
	return &RedisEventPublisher{
		client: client,
		logger: logger,
	}
}

func (p *RedisEventPublisher) PublishWorldCreated(ctx context.Context, world *models.World) error {
	event := WorldEvent{
		Type:      "world.created",
		WorldID:   world.ID,
		UserID:    world.UserID,
		Data:      world,
		Timestamp: time.Now(),
	}

	return p.publishEvent(ctx, "worlds", event)
}

func (p *RedisEventPublisher) PublishWorldUpdated(ctx context.Context, world *models.World) error {
	event := WorldEvent{
		Type:      "world.updated",
		WorldID:   world.ID,
		UserID:    world.UserID,
		Data:      world,
		Timestamp: time.Now(),
	}

	return p.publishEvent(ctx, "worlds", event)
}

func (p *RedisEventPublisher) publishEvent(ctx context.Context, channel string, event WorldEvent) error {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		p.logger.WithError(err).Error("Failed to marshal event")
		return err
	}

	err = p.client.Publish(ctx, channel, eventJSON).Err()
	if err != nil {
		p.logger.WithError(err).WithField("channel", channel).Error("Failed to publish event")
		return err
	}

	p.logger.WithFields(logrus.Fields{
		"channel":  channel,
		"type":     event.Type,
		"world_id": event.WorldID,
	}).Info("Event published successfully")

	return nil
}
