package main

import (
	"database/sql"
	"os"
	"strings"

	_ "github.com/duckdb/duckdb-go/v2"
	"github.com/gin-gonic/gin"
)

var datassert = os.Getenv("DATASSERT_PATH")

func getDB() (*sql.DB, error) {
	db, err := sql.Open("duckdb", datassert)
	if err != nil {
		return nil, err
	}

	return db, nil
}

type CurieResult struct {
	CURIE          string `json:"CURIE"`
	PREFERRED_NAME string `json:"PREFERRED_NAME"`
	CATEGORY_NAME  string `json:"CATEGORY_NAME"`
	NCBI_TAXON_ID  int    `json:"NCBI_TAXON_ID,omitempty"`
}

func SearchForCuries(c *gin.Context) {
	username := c.Query("username")
	apiKey := c.Query("api-key")

	if !HypatiaAuth(c, username, apiKey) {
		return
	}

	term := c.Query("term")
	if term == "" {
		c.JSON(400, gin.H{"error": "'term' is a required API parameter"})
		return
	}

	term = strings.ToLower(term)

	db, err := getDB()
	if err != nil {
		c.JSON(503, gin.H{"error": err.Error()})
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
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error(), "cause": "The given term doesn't resolve to a valid curie. Try equivalent terms."})
		return
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

func GetCurieInfo(c *gin.Context) {
	username := c.Query("username")
	apiKey := c.Query("api-key")

	if !HypatiaAuth(c, username, apiKey) {
		return
	}

	curie := c.Query("curie")
	if curie == "" {
		c.JSON(400, gin.H{"error": "'curie' is a required API parameter"})
		return
	}

	db, err := getDB()
	if err != nil {
		c.JSON(503, gin.H{"error": err.Error()})
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
	LIMIT 1;
	`

	cu := CurieResult{}
	row := db.QueryRow(query, curie)
	err = row.Scan(&cu.CURIE, &cu.PREFERRED_NAME, &cu.CATEGORY_NAME, &cu.NCBI_TAXON_ID)

	if err != nil {
		c.JSON(404, gin.H{"error": err.Error(), "cause": "The given curie doesn't resolve to a cannonical curie and should be ommited. Check /search-for-curies/ for cannonical alternatives."})
		return
	}

	c.JSON(200, gin.H{"curie": cu})
}
