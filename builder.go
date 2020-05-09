package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

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

func build(repositoryName, repositoryEntry, repositoryDir, outputDir string) error {
	err := buildRepository(repositoryDir, repositoryEntry)
	if err != nil {
		return errors.Wrap(err, "error building repository")
	} else {
		binaryPath := filepath.Join(repositoryDir, repositoryName)
		err = moveFile(binaryPath, outputDir)
		if err != nil {
			return errors.Wrap(err, "error moving artifact to output dir")
		}
	}

	return nil
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

func moveFile(binaryPath, outPath string) error {
	_, binaryName := filepath.Split(binaryPath)

	binaryOutPath := filepath.Join(outPath, binaryName)

	return os.Rename(binaryPath, binaryOutPath)
}
