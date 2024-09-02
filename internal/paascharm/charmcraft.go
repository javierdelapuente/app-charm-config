package paascharm

import (
	"io"

	"gopkg.in/yaml.v3"
)

func ParseCharmcraftYAML(reader io.Reader) (CharmcraftYAMLConfig, error) {
	decoder := yaml.NewDecoder(reader)
	var charmcraft CharmcraftYAMLConfig
	err := decoder.Decode(&charmcraft)
	if err != nil {
		return CharmcraftYAMLConfig{}, err
	}
	return charmcraft, nil
}

type CharmcraftYAMLConfig struct {
	Config   CharmcraftConfig
	Requires map[string]CharmcraftIntegration
}

type CharmcraftConfig struct {
	Options map[string]CharmcraftConfigOption
}

type CharmcraftConfigOption struct {
	Type        string
	Default     any
	Description string
}

type CharmcraftIntegration struct {
	Interface string
	Optional  bool
}
