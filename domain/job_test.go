package domain_test

import (
	"github.com/google/uuid"
	"github.com/matheusvidal21/microservice-encoder/domain"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewJob(t *testing.T) {
	video := domain.NewVideo()
	video.ID = uuid.New().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	job, err := domain.NewJob("path", "Converted", video)
	require.Nil(t, err)
	require.NotNil(t, job)
}
