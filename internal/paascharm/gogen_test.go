package paascharm_test

import (
	"fmt"
	"testing"

	"github.com/canonical/paascharmgen/internal/paascharm"
	"github.com/kylelemons/godebug/pretty"
)

func TestParseValidCharmcraftYaml(t *testing.T) {
	yamlData := `
name: go-app
type: charm
base: ubuntu@24.04
platforms:
  amd64:
  arm64:
  armhf:
  ppc64el:
  riscv64:
  s390x:
summary: A very short one-line summary of the Go application.
description: |
  A comprehensive overview of your Go application.
extensions:
  - go-framework
config:
  options:
    user-defined-str:
      type: string
      default: "hello"
      description: user-defined-str Description
    user-defined-int:
      type: int
      default: 100
      description: user-defined-int Description
    user-defined-bool:
      type: bool
      description: user-defined-bool Description
requires:
  mysql:
    interface: mysql_client
    optional: true
    limit: 1
  s3:
    interface: s3
    optional: false

parts: {0-git: {plugin: nil, build-packages: [git]}}
`
	charmcraft, err := paascharm.ParseCharmcraftYaml([]byte(yamlData))

	if err != nil {
		t.Errorf("Error parsing data")
	}

	expected := paascharm.CharmcraftYamlConfig{}
	expected.Config.Options = map[paascharm.ConfigOptionName]paascharm.ConfigOption{
		"user-defined-str":  {Type: "string", Default: "hello", Description: "user-defined-str Description"},
		"user-defined-int":  {Type: "int", Default: 100, Description: "user-defined-int Description"},
		"user-defined-bool": {Type: "bool", Default: nil, Description: "user-defined-bool Description"},
	}
	expected.Requires = map[paascharm.IntegrationName]paascharm.IntegrationConfig{
		"mysql": {Interface: "mysql_client", Optional: true},
		"s3":    {Interface: "s3", Optional: false},
	}

	if diff := pretty.Compare(charmcraft, expected); diff != "" {
		t.Errorf("charmcraft yaml is not correctly parsed. diff: %s\n", diff)
	}
}

func TestConfigOptionNames(t *testing.T) {
	var tests = []struct {
		configOptionName paascharm.ConfigOptionName
		goVarName        string
		envVarName       string
	}{
		{
			configOptionName: "user-config-option",
			goVarName:        "UserConfigOption",
			envVarName:       "APP_USER_CONFIG_OPTION",
		},
		{
			configOptionName: "useroption",
			goVarName:        "Useroption",
			envVarName:       "APP_USEROPTION",
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.configOptionName)
		t.Run(testname, func(t *testing.T) {
			if tt.configOptionName.GoVarName() != tt.goVarName {
				t.Errorf("Wrong GoVarName. Actual %s, expected %s", tt.configOptionName.GoVarName(), tt.goVarName)
			}
			if tt.configOptionName.EnvVarName() != tt.envVarName {
				t.Errorf("Wrong Env Var Name. Actual %s, expected %s", tt.configOptionName.EnvVarName(), tt.envVarName)
			}
		})
	}
}

func TestConfigOptionTypes(t *testing.T) {
	var tests = []struct {
		configOption paascharm.ConfigOption
		goType       string
	}{
		{
			configOption: paascharm.ConfigOption{
				Type:    "bool",
				Default: nil,
			},
			goType: "*bool",
		},
		{
			configOption: paascharm.ConfigOption{
				Type:    "boolean",
				Default: true,
			},
			goType: "bool",
		},
		{
			configOption: paascharm.ConfigOption{
				Type: "int",
			},
			goType: "*int",
		},
		{
			configOption: paascharm.ConfigOption{
				Type:    "string",
				Default: "mystring",
			},
			goType: "string",
		},
		{
			configOption: paascharm.ConfigOption{
				Type: "secret",
			},
			goType: "*string",
		},
		{
			configOption: paascharm.ConfigOption{
				Type: "unknown",
			},
			goType: "*string",
		},
	}

	for ix, tt := range tests {
		testname := fmt.Sprintf("%d", ix)
		t.Run(testname, func(t *testing.T) {
			if tt.configOption.GoType() != tt.goType {
				t.Errorf("Wrong GoType. Actual %s, expected %s", tt.configOption.GoType(), tt.goType)
			}

		})
	}
}
