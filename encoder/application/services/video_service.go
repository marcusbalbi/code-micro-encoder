package services

import (
	"context"
	"encoder/application/repositories"
	"encoder/domain"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
)

//VideoService to manpulate Videos Downloads
type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
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

//Fragment Fragmenta o Video
func (v VideoService) Fragment() error {
	err := os.Mkdir(os.Getenv("localstoragepath")+"/"+v.Video.ID, os.ModePerm)
	if err != nil {
		return err
	}
	source := os.Getenv("localstoragepath") + "/" + v.Video.ID + ".mp4"
	target := os.Getenv("localstoragepath") + "/" + v.Video.ID + ".frag"

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	PrintOutput(output)

	return nil
}

//Encode Encode a video
func (v VideoService) Encode() error {
	cmdArgs := []string{}

	cmdArgs = append(cmdArgs, os.Getenv("localstoragepath")+"/"+v.Video.ID+".frag")
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, os.Getenv("localstoragepath")+"/"+v.Video.ID)
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin")

	cmd := exec.Command("mp4dash", cmdArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	PrintOutput(output)

	return nil
}

//Finish Clean everything
func (v VideoService) Finish() error {
	err := os.Remove(os.Getenv("localstoragepath") + "/" + v.Video.ID + ".mp4")
	if err != nil {
		log.Println("Error removing mp4 ", v.Video.ID, ".mp4")
		return err
	}

	err = os.Remove(os.Getenv("localstoragepath") + "/" + v.Video.ID + ".frag")
	if err != nil {
		log.Println("Error removing frag ", v.Video.ID, ".frag")
		return err
	}

	err = os.RemoveAll(os.Getenv("localstoragepath") + "/" + v.Video.ID)
	if err != nil {
		log.Println("Error removing dir ", v.Video.ID)
		return err
	}

	log.Println("Files have been removed!", v.Video.ID)
	return nil
}

func (v VideoService) InsertVideo() error {
	_, err := v.VideoRepository.Insert(v.Video)
	if err != nil {
		return err
	}
	return nil
}

//PrintOutput Exibe no console o resultado de um output caso exista
func PrintOutput(out []byte) {
	if len(out) > 0 {
		log.Printf("========> Output: %s\n", string(out))
	}
}
