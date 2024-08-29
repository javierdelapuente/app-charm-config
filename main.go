package main

import (
	"flag"
	"log"

	"github.com/canonical/paascharmgen/internal/paascharm"
)

const (
	defaultCharmcraftLocation = "./"
	defaultPackageName        = "appconfig"
	defaultOutputFile         = "appconfig.go"
)

func main() {
	charmcraftDir := flag.String("c", defaultCharmcraftLocation, "charmcraft.yaml file directory.")
	packageName := flag.String("p", defaultPackageName, "name of the generated package.")
	outputFile := flag.String("o", defaultOutputFile, "output file. Overwrites the previous file if it exists")
	flag.Parse()

	err := paascharm.CreateGoStructs(*charmcraftDir, *packageName, *outputFile)
	if err != nil {
		flag.Usage()
		log.Fatal(err)
	}
}
