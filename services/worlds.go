package services

import (
	"context"
	"errors"
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
	worlds, err := s.dal.WorldsDAL.GetWorlds()
	if err != nil {
		return nil, err
	}
	if worlds == nil {
		return []models.World{}, nil
	}
	return worlds, nil
}

func (s *WorldsService) GetWorldByID(id uuid.UUID) (*models.World, error) {
	world, err := s.dal.WorldsDAL.GetWorldByID(id)
	if err != nil {
		return nil, err
	}
	return world, nil
}

func (s *WorldsService) GetWorldsByOwnerID(ownerID uuid.UUID) ([]models.World, error) {
	worlds, err := s.dal.WorldsDAL.GetWorldsByOwnerID(ownerID)
	if err != nil {
		return nil, err
	}
	if worlds == nil {
		return []models.World{}, nil
	}
	return worlds, nil
}

func (s *WorldsService) CreateWorld(ownerID uuid.UUID, name, description string) (*models.World, error) {
	world := &models.World{
		ID:          uuid.New(),
		UserID:      ownerID,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.dal.WorldsDAL.CreateWorld(world)
	if err != nil {
		return nil, err
	}

	s.eventPublisher.PublishWorldCreated(context.Background(), world)

	return world, nil
}

func (s *WorldsService) UpdateWorld(userId, worldId uuid.UUID, name, description string) (*models.World, error) {
	world, err := s.dal.WorldsDAL.GetWorldByID(worldId)
	if err != nil {
		return nil, err
	}

	world.Name = name
	world.Description = description
	world.UpdatedAt = time.Now()

	if world.UserID != userId {
		return nil, errors.New("unauthorized")
	}

	err = s.dal.WorldsDAL.UpdateWorld(world)
	if err != nil {
		return nil, err
	}

	s.eventPublisher.PublishWorldUpdated(context.Background(), world)

	return world, nil
}
