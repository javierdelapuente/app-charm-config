package paascharm

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

// CreateGoStructs reads the charmcraft.yaml file from the charmcraftFile location
// and outputs the Go file with the proper structs in the outputFile location,
// using as package name packageName. It will override the outputFile file if
// it exists and create all the parent directories if missing.
func CreateGoStructs(charmcraftFile string, packageName string, outputFile string) (err error) {
	yamlFile, err := os.ReadFile(charmcraftFile)
	if err != nil {
		return fmt.Errorf("cannot read charmcraft.yaml file: %v", err)
	}

	charmcraftConfig, err := ParseCharmcraftYaml(yamlFile)
	if err != nil {
		return fmt.Errorf("error parsing charmcraft.yaml file: %v", err)
	}

	goStructInfo, err := NewGoStructsData(packageName, charmcraftConfig)
	if err != nil {
		return fmt.Errorf("cannot create go structs data: %v", err)
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

// Generate a []byte with the Go file containing the Go structs for a GoStructsData struct.
// The output code is formatted following gofmt style.
func GenerateGoStructs(goStructsData GoStructsData) (goStructs []byte, err error) {
	tmpl, err := template.New("").Parse(GoTemplate)
	if err != nil {
		return
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, goStructsData)
	if err != nil {
		return nil, fmt.Errorf("failed executing go template: %v", err)
	}

	goStructs, err = format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed formatting go code: %v", err)
	}

	return goStructs, nil
}
