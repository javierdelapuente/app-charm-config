package paascharm

import (
	"gopkg.in/yaml.v3"
)

func ParseCharmcraftYaml(yamlData []byte) (charmcraft CharmcraftYamlConfig, err error) {
	err = yaml.Unmarshal(yamlData, &charmcraft)
	if err != nil {
		return charmcraft, err
	}
	return
}

type CharmcraftYamlConfig struct {
	Config   CharmcraftConfig
	Requires map[string]CharmcraftIntegration
}

type CharmcraftConfig struct {
	Options map[string]CharmcraftConfigOption
}

type CharmcraftConfigOption struct {
	Type        string
	Default     interface{}
	Description string
}

type CharmcraftIntegration struct {
	Interface string
	Optional  bool
}
