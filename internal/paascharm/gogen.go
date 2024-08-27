package paascharm

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
	"unicode"

	"gopkg.in/yaml.v3"
)

const CharmcraftFileName string = "charmcraft.yaml"

type CharmcraftYamlConfig struct {
	Config struct {
		Options map[ConfigOptionName]ConfigOption
	}
	Requires RequiresMap
}

type ConfigOptionName string

func (cn ConfigOptionName) EnvVarName() string {
	result := "APP_" + string(cn)
	result = strings.ReplaceAll(result, "-", "_")
	result = strings.ToUpper(result)
	return result
}

func (cn ConfigOptionName) GoVarName() (result string) {
	parts := strings.Split(string(cn), "-")
	for _, part := range parts {
		partRunes := []rune(part)
		if len(partRunes) > 0 {
			partRunes[0] = unicode.ToUpper(partRunes[0])
			result += string(partRunes)
		}
	}
	return result
}

type ConfigOption struct {
	Type        string
	Default     interface{}
	Description string
}

func (co ConfigOption) GoType() (result string) {
	switch co.Type {
	case "bool":
		result = "bool"
	case "boolean":
		result = "bool"
	case "float":
		result = "float64"
	case "int":
		result = "int"
	case "secret":
		// TODO IS THIS CORRECT?
		result = "string"
	case "string":
		result = "string"
	default:
		log.Printf("Unknown type for a config option of type: %s. Returning string.\n", co.Type)
		result = "string"
	}
	// If it is nil by default, we may want to differentiate a default value from no environment variable
	if co.Default == nil {
		result = "*" + result
	}
	return
}

type RequiresMap map[IntegrationName]IntegrationConfig

func (m RequiresMap) Contains(key string) bool {
	_, ok := m[IntegrationName(key)]
	return ok
}

func (m RequiresMap) Get(key string) IntegrationConfig {
	return m[IntegrationName(key)]
}

func (m RequiresMap) HasDatabase() bool {
	for _, integration := range m {
		if integration.IsDatabase() {
			return true
		}
	}
	return false
}

type IntegrationName string

// All known integrations. Integrations not in paas-charmer should not be added to the generated code.
func (in IntegrationName) IsKnown() bool {
	switch in {
	case "redis", "mysql", "postgresql", "mongodb", "s3", "saml":
		return true
	}
	return false
}

func (in IntegrationName) GoName() (result string) {
	switch in {
	case "redis":
		result = "Redis"
	case "mysql":
		result = "MySQL"
	case "postgresql":
		result = "PostgreSQL"
	case "mongodb":
		result = "MongoDB"
	case "s3":
		result = "S3"
	case "saml":
		result = "SAML"
	default:
		log.Printf("Invalid integration name: %s")
		// TODO CRASH, SHOULD NOT GET HERE?
	}
	return
}

func (in IntegrationName) EnvPrefix() (result string) {
	switch in {
	case "redis":
		result = "APP_REDIS_"
	case "mysql":
		result = "APP_MYSQL_"
	case "postgresql":
		result = "APP_POSTGRESQL_"
	case "mongodb":
		result = "APP_MONGODB_"
	default:
		// Prefix is just for databases
		result = ""
	}
	return
}

type IntegrationConfig struct {
	Interface string
	Optional  bool
}

func (ic IntegrationConfig) IsDatabase() bool {
	switch ic.Interface {
	case "redis":
		return true
	case "mysql_client":
		return true
	case "postgresql_client":
		return true
	case "mongodb_client":
		return true
	}
	return false
}

type templateConfig struct {
	PackageName string
	Charmcraft  CharmcraftYamlConfig
}

func Generate(charmcraftDir string, packageName string, outputFile string) (err error) {
	yamlFile, err := os.ReadFile(path.Join(charmcraftDir, CharmcraftFileName))
	if err != nil {
		return fmt.Errorf("cannot read charmcraft.yaml file: %v", err)
	}

	charmcraftConfig, err := ParseCharmcraftYaml(yamlFile)
	if err != nil {
		return
	}

	templateConfig := templateConfig{packageName, charmcraftConfig}

	goStructs, err := GenerateGoStructs(templateConfig)
	if err != nil {
		return fmt.Errorf("cannot generate template: %v", err)
	}

	// TODO OVERRIDE BY DEFAULT AND CREATE DIRECTURE STRUCTURE?
	// TODO DRY MODE?
	err = os.WriteFile(outputFile, goStructs, 0644)
	if err != nil {
		return fmt.Errorf("cannot write output file %s: %v", outputFile, err)
	}

	return
}

func ParseCharmcraftYaml(yamlData []byte) (charmcraft CharmcraftYamlConfig, err error) {
	err = yaml.Unmarshal(yamlData, &charmcraft)
	if err != nil {
		return charmcraft, fmt.Errorf("cannot parse charmcraft.yaml file: %v", err)
	}
	return
}

//go:embed config.template
var GoTemplate string

func GenerateGoStructs(templateConfig templateConfig) (goStructs []byte, err error) {
	tmpl, err := template.New("").Parse(GoTemplate)

	if err != nil {
		fmt.Printf("Error Parsing Template file")
		return
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, templateConfig)
	if err != nil {
		return nil, fmt.Errorf("failed executing go template: %v", err)
	}

	fmt.Println(buf.String())

	srcFormatted, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed formatting go code: %v", err)
	}

	return srcFormatted, nil
}
