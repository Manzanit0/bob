package main

import (
	"errors"
	"flag"
	"fmt"
	"path"
)

type config struct {
	repositoryURL   string
	repositoryName  string
	repositoryEntry string
	outputDir       string
}

func main() {
	conf, err := processArgs()
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	err = doMagic(conf)
	if err != nil {
		fmt.Print(err.Error())
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

	_, s := path.Split(*repositoryURL)

	c := config{
		repositoryName:  s,
		repositoryURL:   *repositoryURL,
		repositoryEntry: *entryPoint,
		outputDir:       *outputDir, // TODO check it's a valid path
	}

	return c, nil
}

func doMagic(conf config) error {
	repositoryDir, err := getClonePath(conf.repositoryName)
	if err != nil {
		return err
	}

	err = cloneRepository(conf.repositoryURL, repositoryDir)
	if err != nil {
		return err
	}

	err = build(conf.repositoryName, conf.repositoryEntry, repositoryDir, conf.outputDir)
	if err != nil {
		// Don't return - we want to clean up, deleting the repo.
		fmt.Print(err.Error())
	}

	err = deleteRepository(repositoryDir)
	if err != nil {
		return err
	}

	return err
}
