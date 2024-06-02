package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"crypto/sha256"
	"cubx/internal/streams"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
)

// Function for calculating the hash sum of a directory
func HashDirectory(dir string) (string, error) {
	hash := sha256.New()

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Add file information to the hash
		hash.Write([]byte(info.Name()))
		hash.Write([]byte(info.Mode().String()))
		hash.Write([]byte(info.ModTime().String()))

		// If it's a file, add its contents to the hash
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(hash, file); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

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
	fmt.Println("context build: ", contextDir)
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	// currentHash, err := hashDirectory(contextDir)
	// if err != nil {
	// 	return err
	// }

	// hashFilePath := filepath.Join(contextDir, ".context_hash")

	// previousHash, err := os.ReadFile(hashFilePath)
	// if err != nil && !os.IsNotExist(err) {
	// 	return err
	// }

	// if string(previousHash) != currentHash {
	// fmt.Println("Context has changed, creating new tar archive...")

	tarBuf, err := tarDirectory(contextDir)
	if err != nil {
		return err
	}
	// printTarContents(tarBuf)
	buildOptions := types.ImageBuildOptions{
		Dockerfile: filepath.Base(dockerfilePath),
		Tags:       []string{imageTag},
		Remove:     true,
		// NoCache:    true,
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
	// Save the new hash to a file
	// if err := os.WriteFile(hashFilePath, []byte(currentHash), 0644); err != nil {
	// 	return err
	// }

	// Reading the response from the image builder
	// _, err = io.Copy(os.Stdout, buildResponse.Body)
	// if err != nil {
	// 	return err
	// }
	// } else {
	// 	fmt.Println("Context has not changed, skipping tar archive creation.")
	// }

	return nil
}
