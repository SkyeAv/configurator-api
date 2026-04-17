package main

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

const macOSAppleSiliconZshInstallScript string = "/ssd2/sgoetz/tablassist-suite-macos-apple-silicon-zsh-installer.zsh"

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
