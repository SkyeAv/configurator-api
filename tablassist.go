package main

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

var macOSAppleSiliconZshInstallScript string = os.Getenv("MAC_OS_APPLE_SILICON_ZSH_INSTALL_SCRIPT")

func MacOSAppleSiliconZshInstaller(c *gin.Context) {
	file, err := os.Open(macOSAppleSiliconZshInstallScript)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	c.Header("Content-Type", "application/octet-stream")
	io.Copy(c.Writer, file)
}
