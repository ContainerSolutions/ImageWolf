package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func dockerCliInit() {
	var err error
	CTX = context.Background()
	*CLI, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}
}

// dockerPull pulls an image:tag from a given registry
func dockerPull(registry, repo, tag string) error {

	imageName := fmt.Sprintf("%s/%s:%s", registry, repo, tag)

	out, err := CLI.ImagePull(CTX, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	defer out.Close()
	log.Printf("Pulled image %s/%s:%s\n", registry, repo, tag)

	return err
}

// dockerSave retrieves an image from the docker host as an io.ReadCloser.
func dockerSave(registry, repo, tag string) (io.ReadCloser, error) {

	images := []string{fmt.Sprintf("%s:%s", repo, tag)}

	out, err := CLI.ImageSave(CTX, images)
	if err != nil {
		panic(err)
	}

	return out, err
}

// dockerLoad loads an image in the docker host from the client host.
func dockerLoad(file string) error {

	in, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	log.Println("Loading image from file", file)

	_, err = CLI.ImageLoad(CTX, bufio.NewReader(in), true)

	return err
}

/*
 * END OF FILE!
 */
