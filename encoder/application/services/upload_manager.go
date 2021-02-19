package services

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"cloud.google.com/go/storage"
)

//VideoUpload keek data to upload a video
type VideoUpload struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

//NewVideoUpload create a new VideoUpload
func NewVideoUpload() *VideoUpload {
	return &VideoUpload{}
}

//UploadObject teste
func (vu VideoUpload) UploadObject(ctx context.Context, objectPath string, client *storage.Client) error {
	path := strings.Split(objectPath, os.Getenv("localstoragepath")+"/")

	f, err := os.Open(objectPath)
	if err != nil {
		return err
	}

	defer f.Close()

	wc := client.Bucket(vu.OutputBucket).Object(path[1]).NewWriter(ctx)
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	_, err = io.Copy(wc, f)
	if err != nil {
		return err
	}

	err = wc.Close()
	if err != nil {
		return err
	}

	return nil
}

func (vu VideoUpload) loadPaths() error {
	err := filepath.Walk(vu.VideoPath, func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() {
			vu.Paths = append(vu.Paths, path)
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

//ProcessUpload process an upload
func (vu VideoUpload) ProcessUpload(concurency int, doneUpload chan string) error {
	in := make(chan int, runtime.NumCPU())
	returnChannel := make(chan string)

	err := vu.loadPaths()
	if err != nil {
		return err
	}

	uploadClient, ctx, err := getClientUpload()
	if err != nil {
		return err
	}

	for process := 0; process < concurency; process++ {
		go vu.uploadWorker(ctx, in, returnChannel, uploadClient)
	}

	go func() {
		for x := 0; x < len(vu.Paths); x++ {
			in <- x
		}
		close(in)
	}()

	for r := range returnChannel {
		if r != "" {
			doneUpload <- r
			break
		}
	}

	return nil
}

func (vu *VideoUpload) uploadWorker(ctx context.Context, in chan int, returnChan chan string, uploadClient *storage.Client) {
	for x := range in {
		err := vu.UploadObject(ctx, vu.Paths[x], uploadClient)

		if err != nil {
			vu.Errors = append(vu.Errors, vu.Paths[x])
			log.Printf("Erro during the upload: %v. Error: %v", vu.Paths[x], err)
			returnChan <- err.Error()
		}
		returnChan <- ""
	}
	returnChan <- "Upload Completed"
}

func getClientUpload() (*storage.Client, context.Context, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}
