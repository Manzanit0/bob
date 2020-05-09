package app

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"

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
	err := cmd.Run()

	return errors.Wrap(err, "unable to clone repository")
}

func (a *App) Build(outputDir string) error {
	err := buildRepository(a.TempDir, a.Entry)
	if err != nil {
		return errors.Wrap(err, "error building repository")
	}

	binaryPath := filepath.Join(a.TempDir, a.Name)

	err = moveFile(binaryPath, outputDir)
	if err != nil {
		return errors.Wrap(err, "error moving artifact to output dir")
	}

	return nil
}

func (a *App) DeleteClone() error {
	cmd := exec.Command("rm", "-rf", a.TempDir)
	err := cmd.Run()

	return errors.Wrap(err, "unable to delete directory")
}

func buildRepository(path string, entryPoint string) error {
	cmd := exec.Command("go", "build", entryPoint)
	cmd.Dir = path

	err := cmd.Run()
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
