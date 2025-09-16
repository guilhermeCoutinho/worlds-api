package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/guilhermeCoutinho/worlds-api/models"
	"github.com/guilhermeCoutinho/worlds-api/utils"
	"github.com/sirupsen/logrus"
)

type EventPublisher interface {
	PublishWorldCreated(ctx context.Context, world *models.World)
	PublishWorldUpdated(ctx context.Context, world *models.World)
}

type WorldEvent struct {
	Type      string      `json:"type"`
	WorldID   uuid.UUID   `json:"world_id"`
	UserID    uuid.UUID   `json:"user_id"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

type RedisAsyncEventPublisher struct {
	client *redis.Client
	logger logrus.FieldLogger
}

func NewRedisEventPublisher(client *redis.Client, logger logrus.FieldLogger) *RedisAsyncEventPublisher {
	return &RedisAsyncEventPublisher{
		client: client,
		logger: logger,
	}
}

func (p *RedisAsyncEventPublisher) PublishWorldCreated(ctx context.Context, world *models.World) {
	event := WorldEvent{
		Type:      "world.created",
		WorldID:   world.ID,
		UserID:    world.UserID,
		Data:      world,
		Timestamp: time.Now(),
	}

	p.publishEvent(ctx, "worlds", event)
}

func (p *RedisAsyncEventPublisher) PublishWorldUpdated(ctx context.Context, world *models.World) {
	event := WorldEvent{
		Type:      "world.updated",
		WorldID:   world.ID,
		UserID:    world.UserID,
		Data:      world,
		Timestamp: time.Now(),
	}

	p.publishEvent(ctx, "worlds", event)
}

func (p *RedisAsyncEventPublisher) publishEvent(ctx context.Context, channel string, event WorldEvent) {
	utils.SafeGo(ctx, func() {
		logger := p.logger.WithFields(logrus.Fields{
			"channel":  channel,
			"type":     event.Type,
			"world_id": event.WorldID,
		})

		logger.Debug("Publishing event")
		eventJSON, err := json.Marshal(event)
		if err != nil {
			logger.WithError(err).Error("Failed to marshal event")
			return
		}

		err = p.client.Publish(ctx, channel, eventJSON).Err()
		if err != nil {
			logger.WithError(err).Error("Failed to publish event")
			return
		}
	})

}
