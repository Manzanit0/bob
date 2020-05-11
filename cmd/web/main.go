package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/manzanit0/bob/pkg/app"
)

func main() {
	r := gin.Default()

	var phtml = flag.String("html", "static", "path to the template folder")

	flag.Parse()

	r.LoadHTMLGlob(fmt.Sprintf("%s/*", *phtml))

	r.POST("/build", buildHandler)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	err := r.Run(getPort())
	if err != nil {
		panic(err)
	}
}

func getPort() string {
	// Heroku shenanigans
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}

	return ":8080"
}

type BuildRequest struct {
	RepositoryURL   string `json:"url"`
	RepositoryEntry string `json:"entry_point"`
	TargetOS        string `json:"target_os"`
	TargetArch      string `json:"target_arch"`
}

func buildHandler(c *gin.Context) {
	var b BuildRequest
	c.BindJSON(&b)

	if b.RepositoryURL == "" {
		log.Printf("[WARNING] Request made with empty URL")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid repository url"})
		return
	}

	// TODO validate targetOS and targetArch

	tempDir, err := ioutil.TempDir("", "*")
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	a := app.New(b.RepositoryURL, b.RepositoryEntry, tempDir)

	err = a.Clone()
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer a.DeleteClone()

	outDir, err := ioutil.TempDir("", "*")
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = a.Build(outDir, b.TargetOS, b.TargetArch)
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// To download the file as a binary
	filePath := filepath.Join(outDir, a.Name)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+a.Name)
	c.Header("Content-Type", "application/octet-stream")
	c.File(filePath)
}
