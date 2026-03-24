package main

import (
	"github.com/gin-gonic/gin"
)

func health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

func registerRoutes(r *gin.Engine) {
	r.GET("/health", health)
	r.GET("/search-for-curies", SearchForCuries)
	r.GET("/get-canonical-curie-info", GetCurieInfo)
	r.GET("/download-from-pmc-tars", DownloadFromPMCTars)
}

func main() {
	r := gin.Default()
	registerRoutes(r)
	r.Run(":8550")
}
