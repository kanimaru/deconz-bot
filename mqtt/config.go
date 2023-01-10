package mqtt

import (
	"github.com/PerformLine/go-stockutil/log"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

// AdditionalButtonConfig describes the paths where the button should be shown ex: /group_id/light_id
type AdditionalButtonConfig struct {
	Path map[string]ButtonConfig `yaml:"path"`
}

type ButtonType string

const TypeSingle ButtonType = "single"
const TypeToggle ButtonType = "toggle"
const TypeMulti ButtonType = "multi"

type ButtonConfig struct {
	Name         string                 `yaml:"name"`
	Topic        string                 `yaml:"topic"`
	Message      map[string]interface{} `yaml:"message"`
	Typ          ButtonType             `yaml:"typ"`
	ToggleValues []ButtonConfig         `yaml:"toggleValues"`
	MultiValues  []ButtonConfig         `yaml:"multiValues"`
}

func LoadConfigFile(path string) *AdditionalButtonConfig {
	if path == "" {
		return nil
	}
	reader, err := os.Open(path)
	if err != nil {
		log.Fatalf("Can't load %v cause of: %v", path, err)
	}
	config, err := LoadConfig(reader)
	if err != nil {
		log.Fatalf("Can't decode %v cause of: %v", path, err)
	}
	return config
}

func LoadConfig(reader io.Reader) (*AdditionalButtonConfig, error) {
	if reader == nil {
		return nil, nil
	}
	decoder := yaml.NewDecoder(reader)
	additionalButtonConfig := AdditionalButtonConfig{}
	err := decoder.Decode(&additionalButtonConfig)
	if err != nil {
		return nil, err
	}
	return &additionalButtonConfig, nil
}
