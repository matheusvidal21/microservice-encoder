package services

import (
	"errors"
	"github.com/matheusvidal21/microservice-encoder/application/repositories"
	"github.com/matheusvidal21/microservice-encoder/domain"
	"os"
	"strconv"
)

const (
	STARTING    = "STARTING"
	FAILED      = "FAILED"
	DOWNLOADING = "DOWNLOADING"
	FRAGMENTING = "FRAGMENTING"
	ENCONDING   = "ENCONDING"
	UPLOADING   = "UPLOADING"
	FINISHING   = "FINISHING"
	COMPLETED   = "COMPLETED"
)

type JobService struct {
	Job           *domain.Job
	JobRepository repositories.JobRepository
	VideoService  VideoService
}

func NewJobService() JobService {
	return JobService{}
}

func (j *JobService) Start() error {
	err := j.changeJobStatus(DOWNLOADING)
	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Download(os.Getenv("INPUT_BUCKET_NAME"))
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus(FRAGMENTING)
	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Fragment()
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus(ENCONDING)
	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Encode()
	if err != nil {
		return j.failJob(err)
	}

	err = j.performUpload()
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus(FINISHING)
	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Finish()
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus(COMPLETED)
	if err != nil {
		return j.failJob(err)
	}

	return nil
}

func (j *JobService) performUpload() error {
	err := j.changeJobStatus(UPLOADING)
	if err != nil {
		return j.failJob(err)
	}

	videoUpload := NewVideoUpload()
	videoUpload.OutputBucket = os.Getenv("OUTPUT_BUCKET_NAME")
	videoUpload.VideoPath = os.Getenv("LOCAL_STORAGE_PATH") + "/" + j.VideoService.Video.ID
	concurrency, _ := strconv.Atoi(os.Getenv("CONCURRENCY_UPLOAD"))
	doneUpload := make(chan string)

	go videoUpload.ProcessUpload(concurrency, doneUpload)
	uploadResult := <-doneUpload

	if uploadResult != "upload completed" {
		return j.failJob(errors.New(uploadResult))
	}

	return nil
}

func (j *JobService) changeJobStatus(status string) error {
	var err error
	j.Job.Status = status
	j.Job, err = j.JobRepository.Update(j.Job)

	if err != nil {
		return j.failJob(err)
	}
	return nil
}

func (j *JobService) failJob(error error) error {
	j.Job.Status = FAILED
	j.Job.Error = error.Error()

	_, err := j.JobRepository.Update(j.Job)

	if err != nil {
		return err
	}
	return error
}
