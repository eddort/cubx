package registry

import (
	"context"
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func FetchTags(repository string) ([]string, error) {
	// Parse the repository name.
	repo, err := name.NewRepository(repository)
	if err != nil {
		return nil, fmt.Errorf("parsing repository name: %w", err)
	}

	// List all tags in the repository.
	tags, err := remote.List(repo, remote.WithContext(context.Background()))
	if err != nil {
		return nil, fmt.Errorf("listing repository tags: %w", err)
	}

	return tags, nil
}
