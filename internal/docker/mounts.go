package docker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/mount"
)

func getCurrentDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve the current directory: %w", err)
	}
	return dir, nil
}

// createTempDir creates a temporary empty directory
func createTempDir() (string, error) {
	tempDir, err := os.MkdirTemp("", "empty")
	if err != nil {
		return "", fmt.Errorf("error creating temporary directory: %w", err)
	}
	return tempDir, nil
}

func parseMountsString(input string) mount.Mount {
	parts := strings.Split(input, ":")
	if len(parts) == 1 {
		return mount.Mount{
			Type:   mount.TypeBind,
			Source: parts[0],
			Target: parts[0],
		}
	}
	return mount.Mount{
		Type:   mount.TypeBind,
		Source: parts[0],
		Target: parts[1],
	}
}

func generateMounts(cwd string, ignores []string, user_mounts []string) ([]mount.Mount, error) {
	mounts := []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: cwd,
			Target: "/app",
		},
	}

	for _, s := range user_mounts {
		mount := parseMountsString(s)
		mounts = append(mounts, mount)
	}

	for _, ignore := range ignores {
		// Convert the path to an absolute path
		absPath, err := filepath.Abs(ignore)
		if err != nil {
			return nil, fmt.Errorf("error converting path to absolute: %w", err)
		}

		// Check if the path exists
		fileInfo, err := os.Stat(absPath)
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("path does not exist: %s", absPath)
		}

		if err != nil {
			return nil, fmt.Errorf("error stating path: %w", err)
		}

		var source string
		if fileInfo.IsDir() {
			// Create a temporary empty directory
			source, err = createTempDir()
			if err != nil {
				return nil, err
			}
		} else {
			// Use /dev/null for files
			source = "/dev/null"
		}

		mount := mount.Mount{
			Type:   mount.TypeBind,
			Source: source,
			Target: "/app/" + filepath.Base(absPath),
		}

		mounts = append(mounts, mount)
	}

	return mounts, nil
}
