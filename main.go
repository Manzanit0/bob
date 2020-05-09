package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
)

type config struct {
	repositoryURL        string
	repositoryName       string
	repositoryEntryPoint string
	outputDirectory      string
}

func main() {
	conf, err := processArgs()
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	destPath, err := getClonePath(conf.repositoryName)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	err = cloneRepository(conf.repositoryURL, destPath)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	err = buildRepository(destPath, conf.repositoryEntryPoint)
	if err != nil {
		fmt.Print(err.Error())
	} else {
		binaryPath := filepath.Join(destPath, conf.repositoryName)
		err = moveArtifact(binaryPath, conf.outputDirectory)
		if err != nil {
			fmt.Print(err.Error())
		}
	}

	err = deleteRepository(destPath)
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
		repositoryName:       s,
		repositoryURL:        *repositoryURL,
		repositoryEntryPoint: *entryPoint,
		outputDirectory:      *outputDir, // TODO check it's a valid path
	}

	return c, nil
}

func getClonePath(repositoryName string) (string, error) {
	h, err := os.UserHomeDir()
	if err != nil {
		fmt.Print(err.Error())
		return "", errors.Wrapf(err, "could not get user home dir")
	}

	// TODO do this in a temp dir
	return filepath.Join(h, "bob-workingdir", repositoryName), nil
}

func cloneRepository(url, destPath string) error {
	cmd := exec.Command("git", "clone", url, destPath)
	err := cmd.Run()

	return errors.Wrap(err, "unable to clone repository")
}

func deleteRepository(destPath string) error {
	cmd := exec.Command("rm", "-rf", destPath)
	err := cmd.Run()

	return errors.Wrap(err, "unable to delete repository")
}

func buildRepository(path string, entryPoint string) error {
	cmd := exec.Command("go", "build", entryPoint)
	cmd.Dir = path

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "unable to clone repository")
	}

	return nil
}

func moveArtifact(binaryPath, outPath string) error {
	_, binaryName := filepath.Split(binaryPath)

	binaryOutPath := filepath.Join(outPath, binaryName)

	return os.Rename(binaryPath, binaryOutPath)
}
