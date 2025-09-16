package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/guilhermeCoutinho/worlds-api/dal"
	"github.com/guilhermeCoutinho/worlds-api/models"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type WorldsService struct {
	dal            *dal.DAL
	logger         logrus.FieldLogger
	config         *viper.Viper
	eventPublisher EventPublisher
}

func NewWorldsService(
	config *viper.Viper,
	dal *dal.DAL,
	logger logrus.FieldLogger,
	eventPublisher EventPublisher,
) *WorldsService {
	return &WorldsService{
		dal:            dal,
		logger:         logger,
		config:         config,
		eventPublisher: eventPublisher,
	}
}

func (s *WorldsService) GetWorlds() ([]models.World, error) {
	return s.dal.WorldsDAL.GetWorlds()
}

func (s *WorldsService) GetWorldByID(id string) (*models.World, error) {
	return s.dal.WorldsDAL.GetWorldByID(id)
}

func (s *WorldsService) GetWorldsByOwnerID(ownerID string) ([]models.World, error) {
	return s.dal.WorldsDAL.GetWorldsByOwnerID(ownerID)
}

func (s *WorldsService) CreateWorld(name, description, ownerID string) (*models.World, error) {

	world := &models.World{
		ID:          uuid.New().String(),
		UserID:      ownerID,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	err := s.dal.WorldsDAL.CreateWorld(world)
	if err != nil {
		return nil, err
	}

	// Publish event
	ctx := context.Background()
	if err := s.eventPublisher.PublishWorldCreated(ctx, world); err != nil {
		s.logger.WithError(err).Error("Failed to publish world created event")
		// Don't fail the operation if event publishing fails
	}

	return world, nil
}

func (s *WorldsService) UpdateWorld(id, name, description string) (*models.World, error) {
	world, err := s.dal.WorldsDAL.GetWorldByID(id)
	if err != nil {
		return nil, err
	}

	world.Name = name
	world.Description = description
	world.UpdatedAt = time.Now().Format(time.RFC3339)

	err = s.dal.WorldsDAL.UpdateWorld(world)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	// dispatch event
	SafeGo(ctx, func() {
		if err := s.eventPublisher.PublishWorldUpdated(ctx, world); err != nil {
			s.logger.WithError(err).Error("Failed to publish world updated event")
		}
	})

	return world, nil
}
