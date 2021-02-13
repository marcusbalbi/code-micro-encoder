package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

// Job executa uma operação em um video
type Job struct {
	ID               string    `valid:"uuid"`
	OutputBucketPath string    `valid:"notnull"`
	Status           string    `valid:"notnull"`
	Video            *Video    `valid:"-"`
	VideoID          string    `valid:"-"`
	Error            string    `valid:"-"`
	CreatedAt        time.Time `valid:"-"`
	UpdatedAt        time.Time `valid:"-"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

// NewJob cria um job
func NewJob(output string, status string, video *Video) (j *Job, err error) {
	job := Job{
		OutputBucketPath: output,
		Status:           status,
		Video:            video,
	}

	job.prepare()

	e := job.Validate()

	if e != nil {
		return nil, e
	}

	return &job, nil
}

// Prepare prepara um job com valores iniciais
func (j *Job) prepare() {
	j.ID = uuid.NewV4().String()
	j.CreatedAt = time.Now()
	j.UpdatedAt = time.Now()
}

// Validate valida um Job
func (j *Job) Validate() error {
	_, err := govalidator.ValidateStruct(j)

	if err != nil {
		return err
	}
	return nil
}
