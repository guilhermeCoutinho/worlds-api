package dal

import (
	"github.com/go-pg/pg"
	"github.com/guilhermeCoutinho/worlds-api/models"
)

type UserDAL interface {
	CreateUser(user *models.User) error
}

type UserDALImpl struct {
	db *pg.DB
}

func NewUserDAL(db *pg.DB) *UserDALImpl {
	return &UserDALImpl{db: db}
}

func (d *UserDALImpl) CreateUser(user *models.User) error {
	_, err := d.db.Model(user).Insert()
	return err
}
