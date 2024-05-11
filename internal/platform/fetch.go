package platform

import (
	"log"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func FetchPlatforms(imageName string) *PlatformMap {
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

	var platforms = NewPlatformMap()

	for _, manifest := range indexManifest.Manifests {
		// fmt.Printf("OS: %s, Architecture: %s, OS Version: %s\n", manifest.Platform.OS, manifest.Platform.Architecture, manifest.Platform.OSVersion)
		os := manifest.Platform.OS
		arch := manifest.Platform.Architecture
		platforms.Add(os, arch)
	}

	return platforms
}
