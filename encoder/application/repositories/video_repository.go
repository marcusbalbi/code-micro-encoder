package repository

import (
	"encoder/domain"
	"fmt"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// VideoRepository Repositório para armazenar Videos
type VideoRepository interface {
	Insert(video *domain.Video) (*domain.Video, error)
	Find(id string) (*domain.Video, error)
}

//VideoRepositoryDb um Repositorio de Videos
type VideoRepositoryDb struct {
	Db *gorm.DB
}

//NewVideoRepository Cria um novo repositorio
func NewVideoRepository(db *gorm.DB) *VideoRepositoryDb {
	return &VideoRepositoryDb{Db: db}
}

//Insert Insere um novo Video ao Repositorio
func (repo VideoRepositoryDb) Insert(video *domain.Video) (*domain.Video, error) {
	if video.ID == "" {
		video.ID = uuid.NewV4().String()
	}
	err := repo.Db.Create(video).Error

	if err != nil {
		return nil, err
	}
	return video, nil
}

//Find encontra um video pelo seu ID
func (repo VideoRepositoryDb) Find(id string) (*domain.Video, error) {
	var video domain.Video
	repo.Db.First(&video, "id = ?", id)

	if video.ID == "" {
		return nil, fmt.Errorf("Video does not exist")
	}
	return &video, nil
}
