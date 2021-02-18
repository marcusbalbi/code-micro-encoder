package services

import (
	"context"
	"io"
	"os"
	"path/filepath"
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

func getClientUpload() (*storage.Client, context.Context, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}
