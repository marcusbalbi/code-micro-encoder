package repositories

import (
	"encoder/domain"
	"fmt"

	"github.com/jinzhu/gorm"
)

// JobRepository Repositório para armazenar Jobs
type JobRepository interface {
	Insert(video *domain.Job) (*domain.Job, error)
	Find(id string) (*domain.Job, error)
	Update(job *domain.Job) (*domain.Job, error)
}

//JobRepositoryDb um Repositorio de Jobs
type JobRepositoryDb struct {
	Db *gorm.DB
}

//NewJobRepository Cria um novo repositorio de Jobs
func NewJobRepository(db *gorm.DB) *JobRepositoryDb {
	return &JobRepositoryDb{Db: db}
}

//Insert Insere um novo Job ao Repositorio
func (repo JobRepositoryDb) Insert(job *domain.Job) (*domain.Job, error) {

	err := repo.Db.Create(job).Error

	if err != nil {
		return nil, err
	}
	return job, nil
}

//Find encontra um Job pelo seu ID
func (repo JobRepositoryDb) Find(id string) (*domain.Job, error) {
	var job domain.Job
	repo.Db.Preload("Video").First(&job, "id = ?", id)

	if job.ID == "" {
		return nil, fmt.Errorf("Job does not exist")
	}
	return &job, nil
}

//Update atualiza um Job
func (repo JobRepositoryDb) Update(job *domain.Job) (*domain.Job, error) {
	err := repo.Db.Save(&job).Error

	if err != nil {
		return nil, err
	}
	return job, nil
}
