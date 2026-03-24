package main

import (
	"fmt"
	"os"
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
	return
}
