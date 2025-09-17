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

type Event interface {
	GetType() string
	GetLogMetadata() map[string]interface{}
}

type WorldTransferRequestedEvent struct {
	Type              string    `json:"type"`
	WorldID           uuid.UUID `json:"world_id"`
	UserID            uuid.UUID `json:"user_id"`
	WorldVersion      int       `json:"world_version"`
	TargetEnvironment string    `json:"target_environment"`
	Timestamp         time.Time `json:"timestamp"`
}

func (e *WorldTransferRequestedEvent) GetType() string {
	return e.Type
}

func (e *WorldTransferRequestedEvent) GetLogMetadata() map[string]interface{} {
	return map[string]interface{}{
		"world_id":           e.WorldID,
		"user_id":            e.UserID,
		"world_version":      e.WorldVersion,
		"target_environment": e.TargetEnvironment,
	}
}

type EventPublisher interface {
	PublishWorldCreated(ctx context.Context, world *models.World)
	PublishWorldUpdated(ctx context.Context, world *models.World)
	PublishWorldTransferRequested(ctx context.Context, worldTransferRequestedEvent *WorldTransferRequestedEvent)
}

type WorldEvent struct {
	Type      string      `json:"type"`
	WorldID   uuid.UUID   `json:"world_id"`
	UserID    uuid.UUID   `json:"user_id"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

func (e *WorldEvent) GetType() string {
	return e.Type
}

func (e *WorldEvent) GetLogMetadata() map[string]interface{} {
	return map[string]interface{}{
		"world_id": e.WorldID,
		"user_id":  e.UserID,
	}
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

	p.publishEvent(ctx, "worlds", &event)
}

func (p *RedisAsyncEventPublisher) PublishWorldTransferRequested(ctx context.Context, worldTransferRequestedEvent *WorldTransferRequestedEvent) {
	p.publishEvent(ctx, "worlds", worldTransferRequestedEvent)
}

func (p *RedisAsyncEventPublisher) PublishWorldUpdated(ctx context.Context, world *models.World) {
	event := WorldEvent{
		Type:      "world.updated",
		WorldID:   world.ID,
		UserID:    world.UserID,
		Data:      world,
		Timestamp: time.Now(),
	}

	p.publishEvent(ctx, "worlds", &event)
}

func (p *RedisAsyncEventPublisher) publishEvent(ctx context.Context, channel string, event Event) {
	utils.SafeGo(ctx, func() {
		logger := p.logger.WithFields(logrus.Fields{
			"channel":  channel,
			"type":     event.GetType(),
			"metadata": event.GetLogMetadata(),
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
