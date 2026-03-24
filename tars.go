package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func cleanID(pmcID string) string {
	if strings.Contains(":", pmcID) {
		pmcID, _, _ = strings.Cut(pmcID, ":")
	}

	if !strings.Contains("PMC", pmcID) {
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
	pmcID = cleanID(pmcID)

	suffix := fmt.Sprintf("%v/%v.tar.xz", pmcID[9:], pmcID)
	tarPath := filepath.Join(pmcTars, suffix)

	_, err := os.Stat(tarPath)
	if os.IsNotExist(err) {
		c.JSON(506, gin.H{"error": err.Error(), "cause": "The specified PMC tar package hasn't been downloaded yet."})
		return
	}

	return
}
