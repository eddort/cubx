package registry

import (
	"fmt"
	"log"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func Manifest() {
	// Замените 'docker.io/library/alpine' на нужный вам образ
	imageRef, err := name.ParseReference("docker.io/library/node:lts-iron")
	if err != nil {
		log.Fatalf("Ошибка парсинга ссылки на образ: %v", err)
	}

	// Получаем манифест образа
	img, err := remote.Image(imageRef, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		log.Fatalf("Ошибка получения образа: %v", err)
	}

	manifest, err := img.Manifest()
	if err != nil {
		log.Fatalf("failed to retrieve image manifest: %v", err)
	}
	fmt.Println(manifest.Layers[0].URLs)
	platforms := make(map[string]bool)
	for _, layer := range manifest.Layers {
		// layer.Platform.Architecture
		// p := layer.Platform
		if err != nil {
			log.Printf("failed to get platform for layer: %v", err)
			continue
		}
		// fmt.Println(layer.Digest)

		imageRef, err := name.ParseReference("node@" + layer.Digest.String())
		fmt.Println("node@" + layer.Digest.String())
		if err != nil {
			log.Fatalf("Ошибка парсинга ссылки на образ: %v", err)
		}

		// Получаем манифест образа
		img, err := remote.Image(imageRef, remote.WithAuthFromKeychain(authn.DefaultKeychain))
		if err != nil {
			log.Printf("Ошибка получения образа: %v", err)
			continue
		}

		config, err := img.ConfigFile()
		if err != nil {
			log.Fatalf("failed to retrieve image manifest: %v", err)
		}

		fmt.Println(config.OS, config.Architecture)
		// platforms[fmt.Sprintf("%s/%s", p.OS, p.Architecture)] = true
	}

	fmt.Println("Список платформ и ОС для образа:", platforms)

	// // Получаем список всех платформ и ОС для образа
	// manifest, err := img.ConfigFile()
	// if err != nil {
	// 	log.Fatalf("Ошибка получения манифеста: %v", err)
	// }

	// fmt.Println(manifest)

	// // index, err := remote.Index(imageRef.Context().Registry)
	// if err != nil {
	// 	log.Fatalf("Ошибка получения индекса: %v", err)
	// }
	// repo, err := name.NewRepository("docker.io")
	// if err != nil {
	// 	log.Fatalf("Ошибка получения списка образов222: %v", err)
	// }
	// references, err := remote.List(repo)
	// if err != nil {
	// 	log.Fatalf("Ошибка получения списка образов1111: %v", err)
	// }
	// // Выводим информацию о платформах и ОС для каждого образа
	// for _, refStr := range references {
	// 	ref, err := name.ParseReference(refStr)
	// 	if err != nil {
	// 		log.Fatalf("Ошибка парсинга ссылки на образ: %v", err)
	// 	}
	// 	img, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	// 	if err != nil {
	// 		log.Fatalf("Ошибка получения образа: %v", err)
	// 	}
	// 	config, err := img.ConfigFile()
	// 	if err != nil {
	// 		log.Fatalf("Ошибка получения манифеста: %v", err)
	// 	}
	// 	fmt.Println(config.Platform)
	// }
	// Выводим информацию о платформах и ОС для каждого образа
	// for _, ref := range references {
	// 	img, err := index.Image(ref)
	// 	if err != nil {
	// 		log.Fatalf("Ошибка получения образа: %v", err)
	// 	}
	// 	manifest, err := img.ConfigFile()
	// 	if err != nil {
	// 		log.Fatalf("Ошибка получения манифеста: %v", err)
	// 	}
	// 	fmt.Printf("Платформа: %s, ОС: %s\n", manifest.Platform, manifest.Platform)
	// }
}
