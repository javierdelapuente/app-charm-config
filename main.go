package main

import (
	"flag"
	"log"

	"github.com/canonical/paascharmgen/internal/paascharm"
)

const (
	defaultCharmcraftLocation = "."
	defaultPackageName        = "charmconfig"
	defaultOutputFile         = "charmconfig/charmconfig.go"
)

func main() {
	charmcraftDir := flag.String("c", defaultCharmcraftLocation, "charmcraft.yaml file directory")
	packageName := flag.String("p", defaultPackageName, "name of the generated package")
	outputFile := flag.String("o", defaultOutputFile, "output file")
	flag.Parse()

	// TODO MAKE OUTPUTFILE MANDATORY?

	err := paascharm.CreateGoStructs(*charmcraftDir, *packageName, *outputFile)
	if err != nil {
		log.Fatal(err)
	}
}
