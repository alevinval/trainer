package home

import (
	"log"
	"os"
	"path/filepath"
)

var (
	// TrainerHome where trainer folders are kept
	TrainerHome string

	// ActivitiesPath where synchronised activities are kept
	ActivitiesPath string
)

func init() {
	userHome := os.Getenv("HOME")

	TrainerHome = filepath.Join(userHome, ".trainer")
	ActivitiesPath = filepath.Join(TrainerHome, "activities")
}

// Bootstrap ensures all required home folder and sub-folders exist.
func Bootstrap() {
	ensureDir(TrainerHome)
	ensureDir(ActivitiesPath)
}

func ensureDir(dir string) {
	fullPath := filepath.Join(dir)
	err := os.MkdirAll(fullPath, os.ModePerm)
	if err != nil {
		log.Printf("error during home folder bootstrap: %s\n", err)
	}
}
