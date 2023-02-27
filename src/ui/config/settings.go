package config

import (
	"github.com/my-cloud/ruthenium/src/file"
	"path/filepath"
)

type Settings struct {
	ParticlesPerToken uint64
}

func NewSettings(configurationPath string) (*Settings, error) {
	path := filepath.Join(configurationPath, "settings.json")
	parser := file.NewJsonParser(path)
	var settings Settings
	err := parser.Parse(&settings)
	return &settings, err
}
