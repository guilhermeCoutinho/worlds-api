package services

import (
	"github.com/guilhermeCoutinho/worlds-api/dal"
	"github.com/guilhermeCoutinho/worlds-api/models"
)

type UserService struct {
	dal *dal.DAL
}

func NewUserService(dal *dal.DAL) *UserService {
	return &UserService{dal: dal}
}

func (s *UserService) CreateUser(user *models.User) error {
	return s.dal.UserDAL.CreateUser(user)
}
