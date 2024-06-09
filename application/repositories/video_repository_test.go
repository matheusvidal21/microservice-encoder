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

func TestVideoRepositoryDb_Insert(t *testing.T) {
	db := database.NewDatabaseTest()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	video := domain.NewVideo()
	video.ID = uuid.New().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(video)

	v, err := repo.Find(video.ID)

	require.NotEmpty(t, v.ID)
	require.Nil(t, err)
	require.Equal(t, v.ID, video.ID)

}
