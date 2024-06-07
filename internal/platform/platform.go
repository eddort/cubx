package platform

import (
	"fmt"
	"github.com/eddort/cubx/internal/registry"
)

type OsArch struct {
	Os   string
	Arch string
}

// Remove any unsupported Os/Arch combo
// list of valid Os/Arch values (see "Optional Environment Variables" section
// of https://go.dev/doc/install/source
// Added linux/s390x as we know System z support already exists
// Keep in sync with _docker_manifest_annotate in contrib/completion/bash/docker
var validOsArches = map[OsArch]bool{
	{Os: "darwin", Arch: "386"}:      true,
	{Os: "darwin", Arch: "amd64"}:    true,
	{Os: "darwin", Arch: "arm"}:      true,
	{Os: "darwin", Arch: "arm64"}:    true,
	{Os: "dragonfly", Arch: "amd64"}: true,
	{Os: "freebsd", Arch: "386"}:     true,
	{Os: "freebsd", Arch: "amd64"}:   true,
	{Os: "freebsd", Arch: "arm"}:     true,
	{Os: "linux", Arch: "386"}:       true,
	{Os: "linux", Arch: "amd64"}:     true,
	{Os: "linux", Arch: "arm"}:       true,
	{Os: "linux", Arch: "arm64"}:     true,
	{Os: "linux", Arch: "ppc64le"}:   true,
	{Os: "linux", Arch: "mips64"}:    true,
	{Os: "linux", Arch: "mips64le"}:  true,
	{Os: "linux", Arch: "riscv64"}:   true,
	{Os: "linux", Arch: "s390x"}:     true,
	{Os: "netbsd", Arch: "386"}:      true,
	{Os: "netbsd", Arch: "amd64"}:    true,
	{Os: "netbsd", Arch: "arm"}:      true,
	{Os: "openbsd", Arch: "386"}:     true,
	{Os: "openbsd", Arch: "amd64"}:   true,
	{Os: "openbsd", Arch: "arm"}:     true,
	{Os: "plan9", Arch: "386"}:       true,
	{Os: "plan9", Arch: "amd64"}:     true,
	{Os: "solaris", Arch: "amd64"}:   true,
	{Os: "windows", Arch: "386"}:     true,
	{Os: "windows", Arch: "amd64"}:   true,
}

func IsValidOsArch(Os string, Arch string) bool {
	// check for existence of this combo
	_, ok := validOsArches[OsArch{Os, Arch}]
	return ok
}

func GetPlatforms(imageName string) (*PlatformMap, error) {
	manifests, err := registry.FetchManifests(imageName)
	if err != nil {
		return nil, fmt.Errorf("error fetch manifests: %w", err)
	}
	var platforms = NewPlatformMap()

	for _, manifest := range *manifests {
		os := manifest.Platform.OS
		arch := manifest.Platform.Architecture
		platforms.Add(os, arch)
	}

	return platforms, nil
}
