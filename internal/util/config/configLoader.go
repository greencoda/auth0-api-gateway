package config

import (
	"errors"

	"github.com/greencoda/confiq"
	yaml_loader "github.com/greencoda/confiq/loaders/yaml"
)

var ErrNoConfigSet = errors.New("configSet is nil")

type ConfigFilename string

func LoadConfigYAML(configFilename ConfigFilename) (*confiq.ConfigSet, error) {
	configSet := confiq.New()

	if err := configSet.Load(
		yaml_loader.Load().FromFile(string(configFilename)),
	); err != nil {
		return nil, err
	}

	return configSet, nil
}

func LoadConfigFromSetWithPrefix[Config any](configSet *confiq.ConfigSet, prefix string) (*Config, error) {
	if configSet == nil {
		return nil, ErrNoConfigSet
	}

	var (
		config        = new(Config)
		decodeOptions = confiq.DecodeOptions{
			confiq.AsStrict(),
		}
	)

	if prefix != "" {
		decodeOptions = append(decodeOptions, confiq.FromPrefix(prefix))
	}

	err := configSet.Decode(config, decodeOptions...)
	if err != nil {
		return nil, err
	}

	return config, nil
}
