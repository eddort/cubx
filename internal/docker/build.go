package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"cubx/internal/streams"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
)

// Function for creating a tar archive from a directory
func tarDirectory(dir string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	err := filepath.Walk(dir, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(file)
		hdr, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		// hdr.Name = file
		relPath, err := filepath.Rel(dir, file)
		if err != nil {
			return err
		}
		hdr.Name = filepath.ToSlash(relPath)
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		if fi.Mode().IsRegular() {
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return buf, nil
}
func PrintTarContents(buf *bytes.Buffer) error {
	tr := tar.NewReader(bytes.NewReader(buf.Bytes()))
	fmt.Println("Contents of tar:")
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(hdr.Name)
	}
	return nil
}

// Basic function to build a Docker image with hash validation
func BuildImage(dockerfilePath, imageTag, contextDir string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	// TODO: force rebuild
	found, err := imageExists(ctx, cli, imageTag)
	if err != nil {
		return err
	}

	if found {
		return nil
	}

	fmt.Println("context build:", contextDir)

	tarBuf, err := tarDirectory(contextDir)
	if err != nil {
		return err
	}
	// printTarContents(tarBuf)
	buildOptions := types.ImageBuildOptions{
		Dockerfile: filepath.Base(dockerfilePath),
		Tags:       []string{imageTag},
		Remove:     true,
	}

	buildResponse, err := cli.ImageBuild(ctx, tarBuf, buildOptions)
	if err != nil {
		return err
	}
	defer buildResponse.Body.Close()
	err = jsonmessage.DisplayJSONMessagesToStream(buildResponse.Body, streams.NewOut(), nil)
	if err != nil {
		return err
	}

	return nil
}
