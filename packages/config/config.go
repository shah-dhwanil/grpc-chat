package config

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func LoadConfig[T any](envPrefix string, configFilePath ...string) (*T, error) {
	k := koanf.New(".")
	// Load config from file
	for _, path := range configFilePath {
		if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
			return nil, fmt.Errorf("Error in loading config file (%v): %w",path, err)
		}
	}
	// Load config from environment variables
	if err := k.Load(env.Provider(".", env.Opt{
		Prefix: envPrefix,
		TransformFunc: func(k, v string) (string, any) {
			return strings.ToLower(strings.TrimPrefix(k, envPrefix)), v
		},
	}), nil); err != nil {
		return nil, fmt.Errorf("Error in loading Environment variables: %w", err)
	}
	payload := new(T)
	if err := k.Unmarshal("", payload); err != nil {
		return nil, fmt.Errorf("Error in unmarshalling config: %w", err)
	}
	return payload, nil
}