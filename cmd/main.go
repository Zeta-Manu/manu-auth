package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Zeta-Manu/manu-auth/config"
	"github.com/Zeta-Manu/manu-auth/internal/application"
)

func main() {
	cwd, _ := os.Getwd()
	filePath := filepath.Join(cwd, "config", "config.yaml")
	fmt.Println(filePath)

	appConfig, err := config.LoadConfig(filePath)
	if err != nil {
		log.Print("error")
	}

	application.NewApplication(*appConfig)
}
