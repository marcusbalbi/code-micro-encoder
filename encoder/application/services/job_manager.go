package services

import (
	"encoder/application/repositories"
	"encoder/domain"
	"encoder/framework/queue"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
)

//JobManager Controls the JobWorkers to encode videos from queue
type JobManager struct {
	Db               *gorm.DB
	Domain           domain.Job
	MessageChannel   chan amqp.Delivery
	JobReturnChannel chan JobWorkerResult
	RabbitMQ         *queue.RabbitMQ
}

//JobNotificationError Message Type
type JobNotificationError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

//NewJobManager Creates a new Job Manager
func NewJobManager(db *gorm.DB, rabbitMQ *queue.RabbitMQ, jobReturnChannel chan JobWorkerResult, messageChannel chan amqp.Delivery) *JobManager {

	return &JobManager{
		Db:               db,
		Domain:           domain.Job{},
		MessageChannel:   messageChannel,
		JobReturnChannel: jobReturnChannel,
		RabbitMQ:         rabbitMQ,
	}

}

//Start start taking jobs
func (j *JobManager) Start(ch *amqp.Channel) {
	videoService := NewVideoService()
	videoService.VideoRepository = repositories.VideoRepositoryDb{Db: j.Db}

	jobService := JobService{
		JobRepository: repositories.JobRepositoryDb{Db: j.Db},
		VideoService:  videoService,
	}

	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY_WORKERS"))

	if err != nil {
		log.Fatalf("Error Loading var: CONCURRENCY_WORKERS")
	}

	for qtdProcess := 0; qtdProcess < concurrency; qtdProcess++ {
		go JobWorker(j.MessageChannel, j.JobReturnChannel, jobService, j.Domain, qtdProcess)
	}

	for jobResult := range j.JobReturnChannel {
		if jobResult.Error != nil {
			err = j.checkParseErrors(jobResult)
		} else {
			err = j.notifySuccess(jobResult, ch)
		}

		if err != nil {
			jobResult.Message.Reject(false)
		}
	}

}

func (j *JobManager) notifySuccess(jobResult JobWorkerResult, ch *amqp.Channel) error {
	/**
		{
	 		"resource_id": "e3a42dd6-788b-11eb-9439-0242ac130002",
	 		"file_path": "convite.mp4"
	}
		**/
	Mutex.Lock()
	jobJSON, err := json.Marshal(jobResult.Job)
	Mutex.Unlock()

	if err != nil {
		return err
	}

	err = j.notify(jobJSON)
	if err != nil {
		return err
	}

	err = jobResult.Message.Ack(false)
	if err != nil {
		return err
	}

	return nil
}

func (j JobManager) checkParseErrors(jobResult JobWorkerResult) error {
	if jobResult.Job.ID != "" {
		log.Printf("MessageID: %v. Error During the Job: %v whit Video %v. Error: %v ",
			jobResult.Message.DeliveryTag,
			jobResult.Job.ID,
			jobResult.Job.Video.ID,
			jobResult.Error.Error())
	} else {
		log.Printf("MessageID: %v. Error Parsing Message: %v", jobResult.Message.DeliveryTag, jobResult.Error.Error())
	}

	errorMsg := JobNotificationError{
		Message: string(jobResult.Message.Body),
		Error:   jobResult.Error.Error(),
	}

	jobJSON, err := json.Marshal(errorMsg)
	if err != nil {
		return err
	}

	err = j.notify(jobJSON)
	if err != nil {
		return err
	}

	err = jobResult.Message.Reject(false)
	if err != nil {
		return err
	}

	return nil

}

func (j *JobManager) notify(jobJSON []byte) error {
	err := j.RabbitMQ.Notify(
		string(jobJSON),
		"application/json",
		os.Getenv("RABBITMQ_NOTIFICATION_EX"),
		os.Getenv("RABBITMQ_NOTIFICATION_ROUTING_KEY"),
	)

	if err != nil {
		return err
	}

	return nil

}
