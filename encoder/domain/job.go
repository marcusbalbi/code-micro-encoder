package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

// Job executa uma operação em um video
type Job struct {
	ID               string    `json:"job_id" valid:"uuid" gorm:"type:uuid;primary_key"`
	OutputBucketPath string    `json:"output_bucket_path" valid:"notnull"`
	Status           string    `json:"status" valid:"notnull"`
	Video            *Video    `json:"video" valid:"-"`
	VideoID          string    `json:"-" valid:"-" gorm:"column:video_id;type:uuid;notnull"`
	Error            string    `valid:"-"`
	CreatedAt        time.Time `json:"created_at" valid:"-"`
	UpdatedAt        time.Time `json:"updated_at" valid:"-"`
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
