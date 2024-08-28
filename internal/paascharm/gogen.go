package paascharm

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"log"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

func CreateGoStructs(charmcraftDir string, packageName string, outputFile string) (err error) {
	yamlFile, err := os.ReadFile(path.Join(charmcraftDir, CharmcraftFileName))
	if err != nil {
		return fmt.Errorf("cannot read charmcraft.yaml file: %v", err)
	}

	charmcraftConfig, err := ParseCharmcraftYaml(yamlFile)
	if err != nil {
		return
	}

	goStructInfo, err := NewGoStructsData(packageName, charmcraftConfig)
	if err != nil {
		return
	}

	goStructs, err := GenerateGoStructs(goStructInfo)
	if err != nil {
		return fmt.Errorf("cannot generate go structs info: %v", err)
	}

	err = os.MkdirAll(filepath.Dir(outputFile), 0755)
	if err != nil {
		return fmt.Errorf("cannot create parent directories for output file %s: %v", outputFile, err)
	}

	err = os.WriteFile(outputFile, goStructs, 0644)
	if err != nil {
		return fmt.Errorf("cannot write output file %s: %v", outputFile, err)
	}
	log.Printf("Configuration written to file: %s\n", outputFile)

	return
}

//go:embed go.tmpl
var GoTemplate string

func GenerateGoStructs(goStructInfo GoStructsData) (goStructs []byte, err error) {
	tmpl, err := template.New("").Parse(GoTemplate)
	if err != nil {
		return
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, goStructInfo)
	if err != nil {
		return nil, fmt.Errorf("failed executing go template: %v", err)
	}

	goStructs, err = format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed formatting go code: %v", err)
	}

	return goStructs, nil
}
