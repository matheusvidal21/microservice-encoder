package repositories_test

import (
	"github.com/google/uuid"
	"github.com/matheusvidal21/microservice-encoder/application/repositories"
	"github.com/matheusvidal21/microservice-encoder/domain"
	"github.com/matheusvidal21/microservice-encoder/framework/database"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestJobRepositoryDb_Insert(t *testing.T) {
	db := database.NewDatabaseTest()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	video := domain.NewVideo()
	video.ID = uuid.New().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.NewVideoRepositoryDb(db)
	repo.Insert(video)

	job, err := domain.NewJob("output", "Pending", video)
	require.Nil(t, err)

	repoJob := repositories.NewJobRepositoryDb(db)
	repoJob.Insert(job)

	j, err := repoJob.Find(job.ID)
	require.NotEmpty(t, j.ID)
	require.Nil(t, err)
	require.Equal(t, j.ID, job.ID)
	require.Equal(t, j.VideoID, video.ID)
}

func TestJobRepositoryDb_Update(t *testing.T) {
	db := database.NewDatabaseTest()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	video := domain.NewVideo()
	video.ID = uuid.New().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.NewVideoRepositoryDb(db)
	repo.Insert(video)

	job, err := domain.NewJob("output", "Pending", video)
	require.Nil(t, err)

	repoJob := repositories.NewJobRepositoryDb(db)
	repoJob.Insert(job)

	job.Status = "Complete"

	repoJob.Update(job)

	j, err := repoJob.Find(job.ID)
	require.NotEmpty(t, j.ID)
	require.Nil(t, err)
	require.Equal(t, j.ID, job.ID)
	require.Equal(t, j.VideoID, video.ID)
	require.Equal(t, j.Status, job.Status)
}
