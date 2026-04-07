package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"website-dm/internal/app"
)

// @title           Pet & Plant Store API
// @version         1.0
// @description     REST API для магазину товарів для тварин та рослин.
// @host            localhost:8080
// @BasePath        /api/v1
// @schemes         http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		if _, err := os.Stat("internal/config/config.yaml"); err == nil {
			cfgPath = "internal/config/config.yaml"
		}
	}
	if cfgPath == "" {
		_, thisFile, _, ok := runtime.Caller(0)
		if !ok {
			log.Fatal("не вдалося визначити шлях до конфіга")
		}
		cfgPath = filepath.Join(filepath.Dir(thisFile), "..", "internal", "config", "config.yaml")
	}

	if err := app.Run(cfgPath); err != nil {
		log.Fatal(err)
	}
}
