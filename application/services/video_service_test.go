package services_test

import (
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	"github.com/matheusvidal21/microservice-encoder/application/repositories"
	"github.com/matheusvidal21/microservice-encoder/application/services"
	"github.com/matheusvidal21/microservice-encoder/domain"
	"github.com/matheusvidal21/microservice-encoder/framework/database"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func prepare() (*domain.Video, repositories.VideoRepositoryDb) {
	db := database.NewDatabaseTest()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	video := domain.NewVideo()
	video.ID = uuid.New().String()
	video.FilePath = "musica.mp4"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}
	return video, repo
}

func TestVideoService_Download(t *testing.T) {
	video, repo := prepare()
	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download("video-storage-project")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	err = videoService.Finish()
	require.Nil(t, err)

}
