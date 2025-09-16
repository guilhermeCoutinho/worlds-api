package dal

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/guilhermeCoutinho/worlds-api/models"
)

type WorldsDAL interface {
	GetWorlds() ([]models.World, error)
	GetWorldByID(id uuid.UUID) (*models.World, error)
	GetWorldsByOwnerID(ownerID uuid.UUID) ([]models.World, error)
	CreateWorld(world *models.World) error
	UpdateWorld(world *models.World) error
	JoinWorld(ctx context.Context, userID, worldID uuid.UUID) error
	GetUserCurrentWorld(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
}

type WorldsDALImpl struct {
	db    *pg.DB
	redis *redis.Client
}

func NewWorldsDAL(db *pg.DB, redisClient *redis.Client) *WorldsDALImpl {
	return &WorldsDALImpl{db: db, redis: redisClient}
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

func (d *WorldsDALImpl) JoinWorld(ctx context.Context, userID, worldID uuid.UUID) error {
	scriptPath := filepath.Join(".", "dal", "join_world.lua")
	scriptContent, err := ioutil.ReadFile(scriptPath)
	if err != nil {
		return err
	}

	script := redis.NewScript(string(scriptContent))

	_, err = script.Run(ctx, d.redis, []string{userID.String()}, worldID.String()).Result()
	if err != nil {
		return err
	}

	return nil
}

func (d *WorldsDALImpl) GetUserCurrentWorld(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	worldIDStr, err := d.redis.Get(ctx, "user:"+userID.String()+":world").Result()
	if err != nil {
		if err == redis.Nil {
			return uuid.Nil, nil
		}
		return uuid.Nil, err
	}

	return uuid.Parse(worldIDStr)
}
