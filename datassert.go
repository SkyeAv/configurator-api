package main

import (
	"database/sql"
	"os"
	"strings"

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

type CurieResult struct {
	CURIE          string
	PREFERRED_NAME string
	CATEGORY_NAME  string
	NCBI_TAXON_ID  int
}

func searchForCuries(c *gin.Context, term string) {
	term = strings.ToLower(term)

	db, err := getDB()
	if CheckErr(err) {
		c.JSON(502, gin.H{"error": err.Error()})
		return
	}

	query := `
	SELECT
		C.CURIE,
		C.PREFERRED_NAME,
		G.CATEGORY_NAME,
		C.TAXON_ID
	FROM SYNONYMS S
	JOIN CURIES C ON S.CURIE_ID = C.CURIE_ID
	JOIN CATEGORIES G ON C.CATEGORY_ID = G.CATEGORY_ID
	WHERE S.SYNONYM = ?
	LIMIT 50;
	`
	rows, err := db.Query(query, term)
	if CheckErr(err) {
		c.JSON(503, gin.H{"error": err.Error()})
	}

	defer rows.Close()
	curies := []CurieResult{}

	for rows.Next() {
		cu := CurieResult{}
		_ = rows.Scan(&cu.CURIE, &cu.PREFERRED_NAME, &cu.CATEGORY_NAME, &cu.NCBI_TAXON_ID)
		curies = append(curies, cu)
	}

	c.JSON(200, gin.H{"curies": curies})
}
