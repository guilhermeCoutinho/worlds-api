package services

import (
	"github.com/guilhermeCoutinho/worlds-api/dal"
	"github.com/guilhermeCoutinho/worlds-api/models"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type WorldsService struct {
	dal    *dal.DAL
	logger logrus.FieldLogger
	config *viper.Viper
}

func NewWorldsService(config *viper.Viper, dal *dal.DAL, logger logrus.FieldLogger) *WorldsService {
	return &WorldsService{
		dal:    dal,
		logger: logger,
		config: config,
	}
}

func (s *WorldsService) GetWorlds() ([]models.World, error) {
	return s.dal.WorldsDAL.GetWorlds()
}
