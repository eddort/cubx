package registry

import (
	"log"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func FetchManifests(imageName string) *[]v1.Descriptor {
	ref, err := name.ParseReference(imageName)
	if err != nil {
		log.Fatalf("Error parsing image name: %v", err)
	}

	desc, err := remote.Get(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		log.Fatalf("Error fetching image description: %v", err)
	}

	imageIndex, err := desc.ImageIndex()
	if err != nil {
		log.Fatalf("Error fetching image index: %v", err)
	}

	indexManifest, err := imageIndex.IndexManifest()
	if err != nil {
		log.Fatalf("Error getting index manifest: %v", err)
	}

	return &indexManifest.Manifests
}
