package platform

import (
	"runtime"
)

type PlatformMap struct {
	platforms map[string]OsArch
}

func NewPlatformMap() *PlatformMap {
	return &PlatformMap{
		platforms: make(map[string]OsArch),
	}
}

func (pm *PlatformMap) Add(os, arch string) {
	if isValidOsArch(os, arch) {
		key := os + "/" + arch
		pm.platforms[key] = OsArch{Os: os, Arch: arch}
	}
}

func (pm *PlatformMap) Get() OsArch {
	userOS := runtime.GOOS
	userArch := runtime.GOARCH
	key := userOS + "/" + userArch

	if platform, exists := pm.platforms[key]; exists {
		return platform
	}

	if userArch == "arm64" {
		key = "linux" + "/" + userArch
	}

	if platform, exists := pm.platforms[key]; exists {
		return platform
	}

	return OsArch{Os: "linux", Arch: "amd64"}
}

func (pm *PlatformMap) GetString() string {
	platform := pm.Get()

	return platform.Os + "/" + platform.Arch
}
