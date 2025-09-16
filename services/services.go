package services

import (
	"github.com/guilhermeCoutinho/worlds-api/dal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Services struct {
	WorldsService *WorldsService
}

func NewServices(
	config *viper.Viper,
	dal *dal.DAL,
	logger logrus.FieldLogger,
	eventPublisher EventPublisher,
) *Services {
	worldsService := NewWorldsService(config, dal, logger, eventPublisher)

	return &Services{
		WorldsService: worldsService,
	}
}
