package dal

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/google/uuid"
	"github.com/guilhermeCoutinho/worlds-api/models"
)

type WorldsTransferJobsDAL interface {
	GetWorldsTransferJob(jobId uuid.UUID) (*models.WorldsTransferJob, error)
	UpsertJob(job *models.WorldsTransferJob) error

	GetWorldsTransferByJob(jobId uuid.UUID) ([]models.WorldTransferJob, error)
	UpsertWorldTransferJob(worldTransferJob *models.WorldTransferJob) error
}

type WorldsTransferJobsDALImpl struct {
	db *pg.DB
}

func NewWorldsTransferJobsDAL(db *pg.DB) *WorldsTransferJobsDALImpl {
	return &WorldsTransferJobsDALImpl{db: db}
}

func (d *WorldsTransferJobsDALImpl) GetWorldsTransferJob(jobId uuid.UUID) (*models.WorldsTransferJob, error) {
	job := &models.WorldsTransferJob{}
	err := d.db.Model(job).Where("id = ?", jobId).Select()
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (d *WorldsTransferJobsDALImpl) UpsertJob(job *models.WorldsTransferJob) error {
	_, err := d.db.Model(job).
		OnConflict("(id) DO UPDATE").
		Set("status = ?", job.Status).
		Set("updated_at = ?", time.Now()).
		Insert()
	return err
}

func (d *WorldsTransferJobsDALImpl) GetWorldsTransferByJob(jobId uuid.UUID) ([]models.WorldTransferJob, error) {
	worlds := []models.WorldTransferJob{}
	err := d.db.Model(&worlds).Where("job_id = ?", jobId).Select()
	if err == pg.ErrNoRows {
		return []models.WorldTransferJob{}, nil
	}
	if err != nil {
		return nil, err
	}
	return worlds, nil
}

func (d *WorldsTransferJobsDALImpl) UpsertWorldTransferJob(worldTransferJob *models.WorldTransferJob) error {
	_, err := d.db.Model(worldTransferJob).
		OnConflict("(world_id) DO UPDATE").
		Set("status = ?", worldTransferJob.Status).
		Set("updated_at = ?", time.Now()).
		Insert()
	return err
}
