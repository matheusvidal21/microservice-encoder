package services

import (
	"cloud.google.com/go/storage"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type VideoUpload struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewVideoUpload() *VideoUpload {
	return &VideoUpload{}
}

func (v *VideoUpload) UploadObject(objectPath string, client *storage.Client, ctx context.Context) error {
	//caminho/x/b/arquivo.mp4
	// split: caminho/x/b/
	// [0] caminho/x/b/ [1] arquivo.mp4
	path := strings.Split(objectPath, os.Getenv("LOCAL_STORAGE_PATH")+"/")

	f, err := os.Open(objectPath)
	if err != nil {
		return err
	}
	defer f.Close()

	wc := client.Bucket(v.OutputBucket).Object(path[1]).NewWriter(ctx)

	if _, err = io.Copy(wc, f); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}

func (v *VideoUpload) loadPaths() error {
	err := filepath.Walk(v.VideoPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			v.Paths = append(v.Paths, path)
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (v *VideoUpload) ProcessUpload(concurrency int, doneUpload chan string) error {
	in := make(chan int, runtime.NumCPU()) // arquivo baseado na posição do slice Paths
	returnChannel := make(chan string)

	err := v.loadPaths()
	if err != nil {
		return err
	}

	uploadClient, ctx, err := getClientUpload()
	if err != nil {
		return err
	}

	// Quantidade de workers que irão trabalhar
	for process := 0; process < concurrency; process++ {
		go v.uploadWorker(in, returnChannel, uploadClient, ctx)
	}

	// Fica enviando os arquivos para os workers através do canal in
	go func() {
		for x := 0; x < len(v.Paths); x++ {
			in <- x
		}
		close(in)
	}()

	// Fica recebendo os retornos dos workers
	for r := range returnChannel {
		if r != "" {
			doneUpload <- r
			break
		}
	}

	return nil
}

func (v *VideoUpload) uploadWorker(in chan int, returnChan chan string, uploadClient *storage.Client, ctx context.Context) {
	for x := range in {
		err := v.UploadObject(v.Paths[x], uploadClient, ctx)
		if err != nil {
			v.Errors = append(v.Errors, v.Paths[x])
			log.Printf("Error during the upload: %v. Error: %v", v.Paths[x], err)
			returnChan <- err.Error()
		}
		returnChan <- ""
	}
	returnChan <- "upload completed"
}

func getClientUpload() (*storage.Client, context.Context, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}
