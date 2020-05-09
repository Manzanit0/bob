package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"log"

	"github.com/manzanit0/bob/pkg/app"
	"github.com/pkg/errors"
)

type config struct {
	repositoryURL   string
	repositoryEntry string
	outputDir       string
}

func main() {
	conf, err := processArgs()
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	tempDir, err := ioutil.TempDir("", "*")
	if err != nil {
		log.Fatal(err)
	}

	a := app.New(conf.repositoryURL, conf.repositoryEntry, tempDir)

	err = a.Clone()
	if err != nil {
		log.Fatal(err)
	}

	err = a.Build(conf.outputDir)
	if err != nil {
		// Don't Fatal - we want to clean up the clone.
		log.Print(err)
	}

	err = a.DeleteClone()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nBob, The Builder has finished!")
}

func processArgs() (config, error) {
	repositoryURL := flag.String("repository", "NaR", "repository to build")
	entryPoint := flag.String("entry", ".", "entry point for project build")
	outputDir := flag.String("out", "/Users/manzanit0/Desktop", "output directory")
	flag.Parse()

	if repositoryURL == nil || *repositoryURL == "NaR" {
		return config{}, errors.New("no repository passed")
	}

	if entryPoint == nil {
		return config{}, errors.New("error parsing entry point")
	}

	c := config{
		repositoryURL:   *repositoryURL,
		repositoryEntry: *entryPoint,
		outputDir:       *outputDir, // TODO check it's a valid path
	}

	return c, nil
}
