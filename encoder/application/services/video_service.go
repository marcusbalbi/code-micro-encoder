package services

import (
	"context"
	"encoder/application/repositories"
	"encoder/domain"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"cloud.google.com/go/storage"
)

//VideoService to manpulate Videos Downloads
type VideoService struct {
	Video           *domain.Video
	VideoRepository *repositories.VideoRepository
}

//NewVideoService return a new videoService
func NewVideoService() VideoService {
	return VideoService{}
}

//Download videos to encode
func (v VideoService) Download(bucketName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bkt := client.Bucket(bucketName)
	obj := bkt.Object(v.Video.FilePath)
	r, err := obj.NewReader(ctx)

	if err != nil {
		return err
	}
	defer r.Close()

	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	f, err := os.Create(os.Getenv("localstoragepath") + "/" + v.Video.ID + ".mp4")
	if err != nil {
		return err
	}

	d, err := f.Write(body)
	if err != nil {
		return err
	}

	defer f.Close()
	if d <= 0 {
		return fmt.Errorf("Erro ao gravar arquivo")
	}
	log.Printf("Video %v has been stored", v.Video.ID)
	return nil
}
