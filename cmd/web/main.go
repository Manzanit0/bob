package main

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manzanit0/bob/pkg/app"
)

func main() {
	r := gin.Default()

	r.POST("/build", buildHandler)

	err := r.Run()
	if err != nil {
		panic(err)
	}
}

type BuildRequest struct {
	RepositoryURL   string `json:"url"`
	RepositoryEntry string `json:"entry_point"`
}

func buildHandler(c *gin.Context) {
	var b BuildRequest
	c.BindJSON(&b)

	if b.RepositoryURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid repository url"})
		return
	}

	tempDir, err := ioutil.TempDir("", "*")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	a := app.New(b.RepositoryURL, b.RepositoryEntry, tempDir)

	err = a.Clone()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer a.DeleteClone()

	outDir, err := ioutil.TempDir("", "*")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = a.Build(outDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
