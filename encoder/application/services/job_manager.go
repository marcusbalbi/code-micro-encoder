package services

import (
	"encoder/domain"
	utils "encoder/framework/util"

	"github.com/streadway/amqp"
)

type JobWorkerResult struct {
	Job     domain.Job
	Message *amqp.Delivery
	Error   error
}

func JobWorker(messageChannel chan amqp.Delivery, returnChannel chan JobWorkerResult, jobService JobService, workerID int) {

	for message := range messageChannel {
		err := utils.IsJSON(string(message.Body))
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		// pegar body da mensagem
		// verificar se o json e valido
		// validar o vido
		// inserir o video no banco de dados
		// start
	}
}

func returnJobResult(job domain.Job, message amqp.Delivery, err error) JobWOJobWorkerResult {
	result := JobJobWorkerResult{
		Job:     job,
		Message: &message,
		Error:   err,
	}
	return result
}
