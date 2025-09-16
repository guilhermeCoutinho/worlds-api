package dal

import (
	"github.com/go-pg/pg"
	"github.com/guilhermeCoutinho/worlds-api/models"
)

type WorldsDAL interface {
	GetWorlds() ([]models.World, error)
}

type WorldsDALImpl struct {
	db *pg.DB
}

func NewWorldsDAL(db *pg.DB) *WorldsDALImpl {
	return &WorldsDALImpl{db: db}
}

func (d *WorldsDALImpl) GetWorlds() ([]models.World, error) {
	var worlds []models.World
	err := d.db.Model(&worlds).Select()
	return worlds, err
}
