package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func cleanID(pmcID string) string {
	if strings.Contains(pmcID, ":") {
		_, pmcID, _ = strings.Cut(pmcID, ":")
	}

	if !strings.Contains(pmcID, "PMC") {
		pmcID = fmt.Sprintf("PMC%v", pmcID)
	}

	return pmcID
}

var pmcTars = os.Getenv("PMC_TARS_PATH")

func DownloadFromPMCTars(c *gin.Context) {
	username := c.Query("username")
	apiKey := c.Query("api-key")

	if !HypatiaAuth(c, username, apiKey) {
		return
	}

	pmcID := c.Query("pmc-id")
	if pmcID == "" {
		c.JSON(400, gin.H{"error": "'pmc-id' is a required API parameter"})
		return
	}

	pmcID = cleanID(pmcID)
	suffix := fmt.Sprintf("%v/%v.tar.xz", pmcID[9:], pmcID)
	tarPath := filepath.Join(pmcTars, suffix)

	file, err := os.Open(tarPath)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error(), "cause": "The specified PMC tar package hasn't been downloaded yet."})
		return
	}
	defer file.Close()

	disposition := fmt.Sprintf("attachment; filename=%v.tar.xz", pmcID)

	c.Header("Content-Disposition", disposition)
	c.Header("Content-Type", "application/octet-stream")
	io.Copy(c.Writer, file)
}
