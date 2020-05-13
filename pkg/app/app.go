package app

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type App struct {
	Name      string
	RemoteURL string
	Entry     string
	TempDir   string
}

func New(remoteURL, entry, tempDir string) App {
	_, repoName := path.Split(remoteURL)

	// Needs to be inexistent dir to be able to clone
	tmp := filepath.Join(tempDir, repoName)

	return App{Name: repoName, RemoteURL: remoteURL, Entry: entry, TempDir: tmp}
}

func (a *App) Clone() error {
	cmd := exec.Command("git", "clone", a.RemoteURL, a.TempDir)

	err := run(cmd)
	if err != nil {
		return errors.Wrap(err, "unable to clone repository")
	}

	return nil
}

func (a *App) Build(outputDir, targetOS, targetArch string) error {
	// TODO validate targetOS and targetArch
	err := getAllDependencies(a.TempDir)
	if err != nil {
		return errors.Wrap(err, "error fetching dependencies for repository")
	}

	d, _ := filepath.Split(a.Entry)
	d = strings.Trim(d, "/")
	compilationDir := filepath.Join(a.TempDir, d)

	err = buildRepository(compilationDir, a.Name, targetOS, targetArch)
	if err != nil {
		return errors.Wrap(err, "error building repository")
	}

	binaryPath := filepath.Join(compilationDir, a.Name)

	err = moveFile(binaryPath, outputDir)
	if err != nil {
		return errors.Wrap(err, "error moving artifact to output dir")
	}

	return nil
}

func (a *App) DeleteClone() error {
	cmd := exec.Command("rm", "-rf", a.TempDir)

	err := run(cmd)
	if err != nil {
		return errors.Wrap(err, "unable to delete repository")
	}

	return nil
}

func getAllDependencies(path string) error {
	cmd := exec.Command("go", "get", "./...")
	cmd.Dir = path

	err := run(cmd)
	if err != nil {
		return errors.Wrap(err, "unable to fetch dependencies")
	}

	return nil
}

func buildRepository(path, outName, targetOS, targetArch string) error {
	cmd := exec.Command("go", "build", "-i", "-o", outName, ".")
	cmd.Dir = path

	cmd.Env = append(os.Environ(), "GOOS="+targetOS)
	cmd.Env = append(cmd.Env, "GOARCH="+targetArch)

	err := run(cmd)
	if err != nil {
		return errors.Wrap(err, "unable to build repository")
	}

	return nil
}

func moveFile(filePath, outPath string) error {
	_, fileName := filepath.Split(filePath)

	fileOutPath := filepath.Join(outPath, fileName)

	return os.Rename(filePath, fileOutPath)
}

// run simply runs the command and in case of error, composes the error message
// with the combination of the the exit status and stderr.
func run(cmd *exec.Cmd) error {
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	fmt.Printf("\ncmd %s", cmd.String())
	fmt.Printf("\ncmd path %s\n", cmd.Dir)

	err := cmd.Run()
	if err != nil {
		return errors.New(fmt.Sprint(err) + ": " + stderr.String())
	}

	return nil
}
