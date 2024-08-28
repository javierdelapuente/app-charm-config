package paascharm_test

import (
	"testing"

	"github.com/canonical/paascharmgen/internal/paascharm"
	"github.com/kylelemons/godebug/pretty"
)

func TestGenerateGoStructsFromCharmcraftYaml(t *testing.T) {
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
	packageName := "configpackage"
	charmcraft, err := paascharm.ParseCharmcraftYaml([]byte(yamlData))
	if err != nil {
		t.Errorf("Error parsing charmcraft.yaml")
	}
	goStructs, err := paascharm.NewGoStructsData(packageName, charmcraft)
	if err != nil {
		t.Errorf("Error parsing charmcraft.yaml")
	}
	expected := paascharm.GoStructsData{
		PackageName: "configpackage",
		Options: []paascharm.Option{
			{
				GoVarName:  "UserDefinedStr",
				GoVarType:  "string",
				EnvVarName: "APP_USER_DEFINED_STR",
			},
			{
				GoVarName:  "UserDefinedInt",
				GoVarType:  "int",
				EnvVarName: "APP_USER_DEFINED_INT",
			},
			{
				GoVarName:  "UserDefinedBool",
				GoVarType:  "*bool",
				EnvVarName: "APP_USER_DEFINED_BOOL",
			},
		},
		HasDatabaseIntegrations: true,
		Integrations: map[string]paascharm.Integration{
			"mysql": {
				Name:           "mysql",
				GoName:         "MySQL",
				Interface:      "mysql_client",
				Optional:       true,
				IsDatabase:     true,
				DatabasePrefix: "APP_MYSQL_",
			},
			"s3": {
				Name:       "s3",
				GoName:     "S3",
				Interface:  "s3",
				Optional:   false,
				IsDatabase: false,
			},
		},
	}

	if diff := pretty.Compare(goStructs, expected); diff != "" {
		t.Errorf("goStructs is not correctly generated. diff: %s\n", diff)
	}
}

// func TestConfigOptionNames(t *testing.T) {
// 	var tests = []struct {
// 		configOptionName paascharm.ConfigOptionName
// 		goVarName        string
// 		envVarName       string
// 	}{
// 		{
// 			configOptionName: "user-config-option",
// 			goVarName:        "UserConfigOption",
// 			envVarName:       "APP_USER_CONFIG_OPTION",
// 		},
// 		{
// 			configOptionName: "useroption",
// 			goVarName:        "Useroption",
// 			envVarName:       "APP_USEROPTION",
// 		},
// 	}

// 	for _, tt := range tests {
// 		testname := fmt.Sprintf("%s", tt.configOptionName)
// 		t.Run(testname, func(t *testing.T) {
// 			if tt.configOptionName.GoVarName() != tt.goVarName {
// 				t.Errorf("Wrong GoVarName. Actual %s, expected %s", tt.configOptionName.GoVarName(), tt.goVarName)
// 			}
// 			if tt.configOptionName.EnvVarName() != tt.envVarName {
// 				t.Errorf("Wrong Env Var Name. Actual %s, expected %s", tt.configOptionName.EnvVarName(), tt.envVarName)
// 			}
// 		})
// 	}
// }

// func TestConfigOptionTypes(t *testing.T) {
// 	var tests = []struct {
// 		configOption paascharm.ConfigOption
// 		goType       string
// 	}{
// 		{
// 			configOption: paascharm.ConfigOption{
// 				Type:    "bool",
// 				Default: nil,
// 			},
// 			goType: "*bool",
// 		},
// 		{
// 			configOption: paascharm.ConfigOption{
// 				Type:    "boolean",
// 				Default: true,
// 			},
// 			goType: "bool",
// 		},
// 		{
// 			configOption: paascharm.ConfigOption{
// 				Type: "int",
// 			},
// 			goType: "*int",
// 		},
// 		{
// 			configOption: paascharm.ConfigOption{
// 				Type:    "string",
// 				Default: "mystring",
// 			},
// 			goType: "string",
// 		},
// 		{
// 			configOption: paascharm.ConfigOption{
// 				Type: "secret",
// 			},
// 			goType: "*string",
// 		},
p// 		{
// 			configOption: paascharm.ConfigOption{
// 				Type: "unknown",
// 			},
// 			goType: "*string",
// 		},
// 	}

// 	for ix, tt := range tests {
// 		testname := fmt.Sprintf("%d", ix)
// 		t.Run(testname, func(t *testing.T) {
// 			if tt.configOption.GoType() != tt.goType {
// 				t.Errorf("Wrong GoType. Actual %s, expected %s", tt.configOption.GoType(), tt.goType)
// 			}

// 		})
// 	}
// }
