package main

import (
	"database/sql"
	"os"

	_ "github.com/duckdb/duckdb-go/v2"
	"github.com/gin-gonic/gin"
)

var p = os.Getenv("DATASSERT_PATH")

func getDB() (*sql.DB, error) {
	db, err := sql.Open("duckdb", p)
	if CheckErr(err) {
		return nil, err
	}

	return db, nil
}

func searchForCuries(c *gin.Context) {
	db, err := getDB()
	if CheckErr(err) {
		c.JSON(502, gin.H{"error": err.Error()})
	}

	c.JSON(200, "ok")
}
