package services

import (
	"encoder/application/repositories"
	"encoder/domain"
	"errors"
	"os"
	"strconv"
)

//JobService do the steps to prepare a video
type JobService struct {
	Job           *domain.Job
	JobRepository repositories.JobRepository
	VideoService  VideoService
}

//Start the process to prepare a video
func (j *JobService) Start() error {

	// DOWNLOADING
	//--------------------------------------------------------------
	err := j.changeJobStatus("DOWNLOADING")
	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Download(os.Getenv("inputBucketName"))
	if err != nil {
		return j.failJob(err)
	}
	//--------------------------------------------------------------

	// FRAGMENTING
	//--------------------------------------------------------------
	err = j.changeJobStatus("FRAGMENTING")
	if err != nil {
		return j.failJob(err)
	}
	err = j.VideoService.Fragment()
	if err != nil {
		return j.failJob(err)
	}
	//--------------------------------------------------------------

	// ENCODING
	//--------------------------------------------------------------
	err = j.changeJobStatus("ENCODING")
	if err != nil {
		return j.failJob(err)
	}
	err = j.VideoService.Encode()
	if err != nil {
		return j.failJob(err)
	}
	//--------------------------------------------------------------

	// UPLOADING
	//--------------------------------------------------------------

	err = j.performUpload()
	if err != nil {
		return j.failJob(err)
	}
	//--------------------------------------------------------------

	// FINISHING
	//--------------------------------------------------------------
	err = j.changeJobStatus("FINISHING")
	if err != nil {
		return j.failJob(err)
	}
	err = j.VideoService.Finish()
	if err != nil {
		return j.failJob(err)
	}
	//--------------------------------------------------------------

	// COMPLETED
	//--------------------------------------------------------------
	err = j.changeJobStatus("COMPLETED")
	if err != nil {
		return j.failJob(err)
	}
	//--------------------------------------------------------------
	return nil
}

func (j *JobService) performUpload() error {
	err := j.changeJobStatus("UPLOADING")
	if err != nil {
		return j.failJob(err)
	}

	videoUpload := NewVideoUpload()
	videoUpload.OutputBucket = os.Getenv("outputBucketName")
	videoUpload.VideoPath = os.Getenv("localstoragepath") + "/" + j.VideoService.Video.ID
	concurrency, _ := strconv.Atoi(os.Getenv("CONCURRENCY_UPLOAD"))
	doneUpload := make(chan string)

	go videoUpload.ProcessUpload(concurrency, doneUpload)

	var uploadResult string

	uploadResult = <-doneUpload

	if uploadResult != "Upload Completed" {
		return j.failJob(errors.New(uploadResult))
	}

	return err

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

func (j *JobService) failJob(err error) error {
	j.Job.Status = "FAILED"
	j.Job.Error = err.Error()

	_, saveErr := j.JobRepository.Update(j.Job)

	if saveErr != nil {
		return saveErr
	}

	return err
}
