package main

import "github.com/gin-gonic/gin"

func DownloadFromPMCTars(c *gin.Context) {
	username := c.Query("username")
	apiKey := c.Query("api-key")

	if !HypatiaAuth(c, username, apiKey) {
		return
	}
}
