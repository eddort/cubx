package registry

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func FetchManifests(imageName string) (*[]v1.Descriptor, error) {
	ref, err := name.ParseReference(imageName)
	if err != nil {
		return nil, fmt.Errorf("error parsing image name: %w", err)
	}

	desc, err := remote.Get(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return nil, fmt.Errorf("error fetching image description: %w", err)
	}

	imageIndex, err := desc.ImageIndex()
	if err != nil {
		return nil, fmt.Errorf("error fetching image index: %w", err)
	}

	indexManifest, err := imageIndex.IndexManifest()
	if err != nil {
		return nil, fmt.Errorf("error getting index manifest: %w", err)
	}

	return &indexManifest.Manifests, nil
}
