package main

import (
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
)

func main() {
	ASCIIloader()
	Separator()

	mode, pseudo, token := setVarEnvs()
	Separator()

	url := setURL(mode, pseudo, token)
	Separator()

	repositories := getRepositories(url, token)
	var wg sync.WaitGroup
	wg.Add(1)
	go writeCSV(pseudo, repositories, &wg)

	for _, repo := range repositories {
		wg.Add(1)
		fmt.Printf("Cloning %s...\n", repo.Name)
		go cloneAndPullRepo(pseudo, repo, token, &wg)
	}
	wg.Wait()
	Separator()

	zipRepositoriesFolder(pseudo)
	Separator()

	fmt.Printf("All repositories have been cloned.\n")
	Separator()

	downloadPath := "/download"
	fiberApp := fiber.New()
	fiberApp.Get(downloadPath, func(c *fiber.Ctx) error {
		return c.Download(fmt.Sprintf("../assets/%s.zip", pseudo))
	})
	fmt.Printf("Exposing the zip to be download at path: %s\n", downloadPath)
	fiberApp.Listen(":3000")
}
