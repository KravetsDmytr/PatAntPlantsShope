package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"website-dm/internal/app"
)

func main() {
	// У Docker после сборки бинарника исходники могут отсутствовать,
	// поэтому основной путь берем из env.
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		// Сценарий go run из корня проекта (или когда рабочая директория "такая как надо").
		if _, err := os.Stat("internal/config/config.yaml"); err == nil {
			cfgPath = "internal/config/config.yaml"
		}
	}
	if cfgPath == "" {
		// Фолбэк для go run из папки cmd/.
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
