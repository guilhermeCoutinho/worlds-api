package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/guilhermeCoutinho/worlds-api/dal"
	"github.com/guilhermeCoutinho/worlds-api/models"
	"github.com/sirupsen/logrus"
)

type WorldsImporterService struct {
	eventPublisher EventPublisher
	dal            *dal.DAL
	logger         logrus.FieldLogger
}

func NewWorldsImporterService(eventPublisher EventPublisher, dal *dal.DAL) *WorldsImporterService {
	return &WorldsImporterService{eventPublisher: eventPublisher, dal: dal}
}

func (s *WorldsImporterService) MakeRequest(ctx context.Context, url string) (*models.World, error) {
	return nil, nil
}

func (s *WorldsImporterService) GetEnvironmentURL(ctx context.Context, worldId uuid.UUID, targetEnvironment string) string {
	return ""
}

func (s *WorldsImporterService) CreateImportWorldsJob(ctx context.Context, userId uuid.UUID, worlds []uuid.UUID, targetEnvironment string) (*models.WorldTransferJobStatusDTO, error) {
	response := &models.WorldTransferJobStatusDTO{
		JobId:           uuid.New(),
		Status:          models.WorldTransferJobStatusCreated,
		StatusByWorldID: make(map[uuid.UUID]models.WorldTransferJobStatus),
	}

	allWorldsUpToDate := true
	worldCurrentVersion := make(map[uuid.UUID]int)
	for _, worldID := range worlds {
		world, err := s.dal.WorldsDAL.GetWorldByID(worldID)
		if err != nil {
			return nil, err
		}

		worldCurrentVersion[worldID] = world.Version

		targetEnvironmentWorld, err := s.MakeRequest(ctx, s.GetEnvironmentURL(ctx, worldID, targetEnvironment))
		if err != nil {
			return nil, err
		}

		if targetEnvironmentWorld.Version <= world.Version {
			s.logger.WithField("world_id", worldID).Info("World version is already up to date")
			response.StatusByWorldID[worldID] = models.WorldTransferJobStatusCompleted
			continue
		}

		s.eventPublisher.PublishWorldTransferRequested(ctx, &WorldTransferRequestedEvent{
			WorldID:           world.ID,
			UserID:            world.UserID,
			WorldVersion:      world.Version,
			TargetEnvironment: targetEnvironment,
		})
		response.StatusByWorldID[worldID] = models.WorldTransferJobStatusCreated
		allWorldsUpToDate = false
	}

	err := s.dal.WorldsTransferJobsDAL.UpsertJob(&models.WorldsTransferJob{
		ID:        response.JobId,
		Status:    response.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	if allWorldsUpToDate {
		response.Status = models.WorldTransferJobStatusCompleted
	} else {
		for worldId, worldStatus := range response.StatusByWorldID {
			s.dal.WorldsTransferJobsDAL.UpsertWorldTransferJob(&models.WorldTransferJob{
				WorldID:      worldId,
				JobId:        response.JobId,
				WorldVersion: worldCurrentVersion[worldId],
				Status:       worldStatus,
			})
		}
	}

	return response, nil
}

func (s *WorldsImporterService) GetAndUpdateWorldsTransferJobStatus(ctx context.Context, jobId uuid.UUID) (*models.WorldTransferJobStatusDTO, error) {
	job, err := s.dal.WorldsTransferJobsDAL.GetWorldsTransferJob(jobId)
	if err != nil {
		return nil, err
	}

	// if completed, no need to fetch individual world statuses, just return
	if job.Status == models.WorldTransferJobStatusCompleted {
		return &models.WorldTransferJobStatusDTO{
			JobId:           jobId,
			Status:          job.Status,
			StatusByWorldID: nil,
		}, nil
	}

	statusByWorldID := make(map[uuid.UUID]models.WorldTransferJobStatus)
	worldTransferJobs, err := s.dal.WorldsTransferJobsDAL.GetWorldsTransferByJob(jobId)
	if err != nil {
		return nil, err
	}
	for _, worldTransferJob := range worldTransferJobs {
		statusByWorldID[worldTransferJob.WorldID] = worldTransferJob.Status
		if worldTransferJob.Status == models.WorldTransferJobStatusCompleted {
			continue
		}

		worldAtTargetEnvironment, err := s.MakeRequest(ctx, s.GetEnvironmentURL(ctx, worldTransferJob.WorldID, job.TargetEnvironment))
		if err != nil {
			s.logger.WithField("world_id", worldTransferJob.WorldID).Error("Error making request to target environment")
			return nil, err
		}
		if worldAtTargetEnvironment.Version <= worldTransferJob.WorldVersion {
			statusByWorldID[worldTransferJob.WorldID] = models.WorldTransferJobStatusCompleted
		}
	}

	areAllWorldsUpToDate := true
	for _, worldStatus := range statusByWorldID {
		if worldStatus != models.WorldTransferJobStatusCompleted {
			areAllWorldsUpToDate = false
			break
		}
	}

	// save it as completed
	if areAllWorldsUpToDate {
		job.Status = models.WorldTransferJobStatusCompleted
		job.UpdatedAt = time.Now()
		err = s.dal.WorldsTransferJobsDAL.UpsertJob(job)
		if err != nil {
			s.logger.WithField("job_id", jobId).Error("Error updating job status")
		}
	}

	return &models.WorldTransferJobStatusDTO{
		JobId:           jobId,
		Status:          job.Status,
		StatusByWorldID: statusByWorldID,
	}, nil
}
