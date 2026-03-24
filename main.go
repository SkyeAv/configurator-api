package main

import (
	"github.com/gin-gonic/gin"
)

func health(c *gin.Context) {
	c.JSON(200, gin.H{"ok": 200})
}

func registerRoutes(r *gin.Engine) {
	r.GET("/health/", health)
	r.GET("/search-for-curies/", SearchForCuries)
	r.GET("/get-cannonical-curie-info/", GetCurieInfo)
}

func main() {
	r := gin.Default()
	registerRoutes(r)
	r.Run(":8550")
}
