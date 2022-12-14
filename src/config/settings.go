package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Settings struct {
	ParticlesCount uint64
	GenesisAmount  uint64
}

func NewSettings(configurationPath string) (*Settings, error) {
	jsonFile, err := os.Open(configurationPath + "/settings.json")
	if err != nil {
		return nil, fmt.Errorf("unable to open settings configuration file: %w", err)
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read settings configuration file: %w", err)
	}
	if err = jsonFile.Close(); err != nil {
		return nil, fmt.Errorf("unable to close settings configuration file: %w", err)
	}
	var settings Settings
	if err = json.Unmarshal(byteValue, &settings); err != nil {
		return nil, fmt.Errorf("unable to unmarshal settings: %w", err)
	}
	return &settings, nil
}
