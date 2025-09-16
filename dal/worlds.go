package dal

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/google/uuid"
	"github.com/guilhermeCoutinho/worlds-api/models"
)

type WorldsDAL interface {
	GetWorlds() ([]models.World, error)
	GetWorldByID(id uuid.UUID) (*models.World, error)
	GetWorldsByOwnerID(ownerID uuid.UUID) ([]models.World, error)
	CreateWorld(world *models.World) error
	UpdateWorld(world *models.World) error
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

func (d *WorldsDALImpl) GetWorldByID(id uuid.UUID) (*models.World, error) {
	world := &models.World{}
	err := d.db.Model(world).Where("id = ?", id).Select()
	if err != nil {
		return nil, err
	}
	return world, nil
}

func (d *WorldsDALImpl) GetWorldsByOwnerID(ownerID uuid.UUID) ([]models.World, error) {
	var worlds []models.World
	err := d.db.Model(&worlds).Where("user_id = ?", ownerID).Select()
	return worlds, err
}

func (d *WorldsDALImpl) CreateWorld(world *models.World) error {
	world.CreatedAt = time.Now()
	world.UpdatedAt = time.Now()
	_, err := d.db.Model(world).Insert()
	return err
}

func (d *WorldsDALImpl) UpdateWorld(world *models.World) error {
	world.UpdatedAt = time.Now()
	_, err := d.db.Model(world).Where("id = ?", world.ID).Update()
	return err
}
