package paascharm

import (
	"fmt"
	"log"
	"strings"
	"unicode"
)

// From integration key to the go Name
var IntegrationsNames = map[string]string{
	"mongodb":    "MongoDB",
	"mysql":      "MySQL",
	"postgresql": "PostgreSQL",
	"redis":      "Redis",
	"s3":         "S3",
	"saml":       "SAML",
}

// Contains the prefixes for database integrations
var DatabaseIntegrationsPrefixes = map[string]string{
	"mongodb":    "APP_MONGODB_",
	"mysql":      "APP_MYSQL_",
	"postgresql": "APP_POSTGRESQL_",
	"redis":      "APP_REDIS_",
}

type GoStructsData struct {
	PackageName             string
	Options                 []Option
	HasDatabaseIntegrations bool
	Integrations            map[string]Integration
}

type Option struct {
	GoVarName  string
	GoVarType  string
	EnvVarName string
}

type Integration struct {
	Name           string
	GoName         string
	Interface      string
	Optional       bool
	IsDatabase     bool
	DatabasePrefix string
}

func NewGoStructsData(packageName string, charmcraft CharmcraftYamlConfig) (goStructs GoStructsData, err error) {
	goStructs = GoStructsData{
		PackageName: packageName,
	}

	for key, value := range charmcraft.Config.Options {
		err, varType := buildGoVarType(value)
		if err != nil {
			return goStructs, err
		}
		option := Option{
			GoVarName:  buildGoVarName(key),
			GoVarType:  varType,
			EnvVarName: buildEnvVarName(key),
		}
		goStructs.Options = append(goStructs.Options, option)
	}

	goStructs.Integrations = make(map[string]Integration)
	for key, value := range charmcraft.Requires {
		integration := Integration{
			Name:      key,
			Interface: value.Interface,
			Optional:  value.Optional,
		}

		if goName, ok := IntegrationsNames[key]; ok {
			integration.GoName = goName
			if databasePrefix, okDatabase := DatabaseIntegrationsPrefixes[key]; okDatabase {
				goStructs.HasDatabaseIntegrations = true
				integration.IsDatabase = true
				integration.DatabasePrefix = databasePrefix
			}
			goStructs.Integrations[key] = integration
		} else {
			log.Printf("Skipping unknown integration %s\n", key)
		}
	}

	return
}

func buildGoVarName(configName string) (result string) {
	parts := strings.Split(configName, "-")
	for _, part := range parts {
		partRunes := []rune(part)
		if len(partRunes) > 0 {
			partRunes[0] = unicode.ToUpper(partRunes[0])
			result += string(partRunes)
		}
	}
	return result
}

func buildGoVarType(configOption CharmcraftConfigOption) (err error, result string) {
	switch configOption.Type {
	case "bool":
		result = "bool"
	case "boolean":
		result = "bool"
	case "float":
		result = "float64"
	case "int":
		result = "int"
	case "secret":
		result = "string"
	case "string":
		result = "string"
	default:
		err = fmt.Errorf("Unknown type for config option of type: %s", configOption.Type)
		return
	}
	// If it is nil by default, we may want to differentiate a default value from no environment variable
	if configOption.Default == nil {
		result = "*" + result
	}
	return
}

func buildEnvVarName(configName string) string {
	result := "APP_" + configName
	result = strings.ReplaceAll(result, "-", "_")
	result = strings.ToUpper(result)
	return result
}