package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/cespare/xxhash/v2"
	_ "github.com/duckdb/duckdb-go/v2"
	"github.com/gin-gonic/gin"
)

var datassert = os.Getenv("DATASSERT_PATH")

func getDB(shard uint) (*sql.DB, error) {
	p := fmt.Sprintf("%v/data/%d.duckdb", datassert, shard)

	db, err := sql.Open("duckdb", p)
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

const shards = 16

func getShard(term string) uint {
	b := []byte(term)
	h := xxhash.Sum64(b)

	return uint(h) % shards
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
	shard := getShard(term)

	db, err := getDB(shard)
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
	WHERE S.SYNONYM = ?;
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

func cleanTaxon(taxonID string) string {
	if strings.Contains(taxonID, ":") {
		_, taxonID, _ = strings.Cut(taxonID, ":")
	}

	return taxonID
}

func SearchForGeneCuriesInNCBITaxon(c *gin.Context) {
	username := c.Query("username")
	apiKey := c.Query("api-key")

	if !HypatiaAuth(c, username, apiKey) {
		return
	}

	taxonID := c.Query("ncbi-taxon-id")
	if taxonID == "" {
		c.JSON(400, gin.H{"error": "'ncbi-taxon-id' is a required API parameter"})
		return
	}

	taxonID = cleanTaxon(taxonID)

	term := c.Query("term")
	if term == "" {
		c.JSON(400, gin.H{"error": "'term' is a required API parameter"})
		return
	}

	shard := getShard(term)

	db, err := getDB(shard)
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
	WHERE C.TAXON_ID = ? AND G.CATEGORY_NAME = 'Gene' AND S.SYNONYM = ?;
	`
	rows, err := db.Query(query, taxonID, term)
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

	curie = strings.ToLower(curie)
	shard := getShard(curie)

	db, err := getDB(shard)
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
		c.JSON(404, gin.H{"error": err.Error(), "cause": "The given curie doesn't resolve to a canonical curie and should be ommited. Check /search-for-curies for canonical alternatives."})
		return
	}

	c.JSON(200, gin.H{"curie": cu})
}

type TaxonResult struct {
	NCBI_TAXON_ID string `json:"CURIE"`
}

func GetTaxonIDFromName(c *gin.Context) {
	username := c.Query("username")
	apiKey := c.Query("api-key")

	if !HypatiaAuth(c, username, apiKey) {
		return
	}

	name := c.Query("organism-name")
	if name == "" {
		c.JSON(400, gin.H{"error": "'organism-name' is a required API parameter"})
		return
	}

	name = strings.ToLower(name)
	shard := getShard(name)

	db, err := getDB(shard)
	if err != nil {
		c.JSON(503, gin.H{"error": err.Error()})
		return
	}

	query := `
	SELECT C.CURIE
	FROM SYNONYMS S
	JOIN CURIES C ON S.CURIE_ID = C.CURIE_ID
	JOIN CATEGORIES G ON C.CATEGORY_ID = G.CATEGORY_ID
	WHERE S.SYNONYM = ?
		AND G.CATEGORY_NAME = 'OrganismTaxon'
		AND starts_with(C.CURIE, 'NCBITaxon:')
	LIMIT 1;
	`
	tr := TaxonResult{}

	row := db.QueryRow(query, name)
	err = row.Scan(&tr.NCBI_TAXON_ID)

	if err != nil {
		c.JSON(404, gin.H{"error": err.Error(), "cause": "The given organism name doesn't resolve to a valid NCBITaxon ID. Try equivalent terms."})
		return
	}

	taxonID := cleanTaxon(tr.NCBI_TAXON_ID)
	c.JSON(200, gin.H{"ncbi-taxon-id": taxonID})
}
