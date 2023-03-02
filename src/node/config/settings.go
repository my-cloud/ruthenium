package config

import (
	"github.com/my-cloud/ruthenium/src/file"
	"path/filepath"
)

type Settings struct {
	GenesisAmountInParticles         uint64
	MaxOutboundsCount                int
	MinimalTransactionFee            uint64
	NetworkId                        string
	SynchronizationIntervalInSeconds int
	ValidationIntervalInSeconds      int
	VerificationsCountPerValidation  int
}

func NewSettings(configurationPath string) (*Settings, error) {
	path := filepath.Join(configurationPath, "settings.json")
	parser := file.NewJsonParser(path)
	var settings Settings
	err := parser.Parse(&settings)
	return &settings, err
}
