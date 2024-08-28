package paascharm_test

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/canonical/paascharmgen/internal/paascharm"
)

func TestCreateGoStructsWithEmptyCharmcraftYaml(t *testing.T) {
	tmpdir := t.TempDir()
	os.WriteFile(path.Join(tmpdir, "charmcraft.yaml"), []byte(""), 0644)
	packageName := "myconfigcharm"
	outputFile := path.Join(tmpdir, packageName, "config.go")

	err := paascharm.CreateGoStructs(tmpdir, packageName, outputFile)
	if err != nil {
		t.Errorf("Error creating go Structs %v\n", err)
	}

	// Check that the outputfile exists and containes the correct package name
	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Error("Error reading output file")
	}
	strData := string(data)
	if !strings.Contains(strData, "package myconfigcharm") {
		t.Error("Output file does not contain package name")
	}
}
