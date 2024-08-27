package main

import (
	"flag"
	"log"

	"github.com/canonical/paascharmgen/internal/paascharm"
)

func main() {
	charmcraftDir := flag.String("c", ".", "charmcraft.yaml file directory")
	packageName := flag.String("p", "charmconfig", "name of the generated package")
	outputFile := flag.String("o", "charmconfig/charmconfig.go", "output file")
	flag.Parse()

	log.Printf("charmcraft.yaml location: %s\n", *charmcraftDir)
	log.Printf("package name: %s\n", *packageName)
	log.Printf("output file: %s\n", *outputFile)

	err := paascharm.Generate(*charmcraftDir, *packageName, *outputFile)
	if err != nil {
		log.Fatal(err)
	}
}
