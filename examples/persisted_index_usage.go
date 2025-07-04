package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aawadall/bit-scout/internal/index"
	"github.com/aawadall/bit-scout/internal/models"
	"github.com/rs/zerolog/log"
)

func main() {
	// Example 1: Create a new index with database and load existing data (creates DB if needed)
	fmt.Println("=== Example 1: Loading existing index from database (creates if needed) ===")

	index1, err := index.NewPersistedSimpleIndexWithDatabaseAndLoad("./data/index.db")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create index with database")
	}
	defer index1.Close()

	// Check if database was empty
	isEmpty, err := index1.IsDatabaseEmpty()
	if err != nil {
		log.Error().Err(err).Msg("Failed to check if database is empty")
	} else {
		fmt.Printf("Database was empty: %v\n", isEmpty)
	}

	// Get database statistics
	stats, err := index1.GetDatabaseStats()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get database stats")
	} else {
		fmt.Printf("Database stats: %+v\n", stats)
	}

	// Example 2: Create a new index and add documents
	fmt.Println("\n=== Example 2: Creating new index and adding documents ===")

	index2, err := index.NewPersistedSimpleIndexWithDatabase("./data/new_index.db")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create new index")
	}
	defer index2.Close()

	// Configure the index
	config := map[string]interface{}{
		"max_results": 10,
		"dimensions":  []string{"fileSize", "lastModified", "fileExtension"},
	}
	if err := index2.Configure(config); err != nil {
		log.Error().Err(err).Msg("Failed to configure index")
	}

	// Add some sample documents
	documents := []models.Document{
		{
			ID:     "doc1",
			Text:   "This is a sample document about Go programming",
			Source: "/path/to/file1.go",
			Vector: []float64{0.1, 0.2},
			Meta: map[string]string{
				"fileExtension": "go",
				"fileSize":      "1024",
				"filename":      "file1.go",
			},
		},
		{
			ID:     "doc2",
			Text:   "Another document about database systems",
			Source: "/path/to/file2.md",
			Vector: []float64{0.3, 0.4},
			Meta: map[string]string{
				"fileExtension": "md",
				"fileSize":      "2048",
				"filename":      "file2.md",
			},
		},
	}

	// Add documents (this will be persisted asynchronously)
	if err := index2.AddDocuments(documents); err != nil {
		log.Error().Err(err).Msg("Failed to add documents")
	}

	// Search immediately (works from memory)
	results, err := index2.Search("Go programming")
	if err != nil {
		log.Error().Err(err).Msg("Failed to search")
	} else {
		fmt.Printf("Found %d documents matching 'Go programming'\n", len(results))
		for _, doc := range results {
			fmt.Printf("  - %s: %s\n", doc.ID, doc.Text[:50])
		}
	}

	// Advanced search
	results, err = index2.Search("fileExtension=go")
	if err != nil {
		log.Error().Err(err).Msg("Failed to search")
	} else {
		fmt.Printf("Found %d documents with fileExtension=go\n", len(results))
	}

	// Get index statistics
	count, err := index2.Count()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get count")
	} else {
		fmt.Printf("Index contains %d documents\n", count)
	}

	// Example 3: Reloading the same index later
	fmt.Println("\n=== Example 3: Reloading index from database ===")

	// Close the current index
	index2.Close()

	// Reopen and reload the same database
	index3, err := index.NewPersistedSimpleIndexWithDatabaseAndLoad("./data/new_index.db")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to reload index")
	}
	defer index3.Close()

	// Verify data was loaded
	count, err = index3.Count()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get count")
	} else {
		fmt.Printf("Reloaded index contains %d documents\n", count)
	}

	// Search again to verify data is available
	results, err = index3.Search("database")
	if err != nil {
		log.Error().Err(err).Msg("Failed to search")
	} else {
		fmt.Printf("Found %d documents matching 'database'\n", len(results))
	}

	fmt.Println("\n=== All examples completed successfully! ===")

	// Example 4: First-time database creation
	fmt.Println("\n=== Example 4: First-time database creation ===")
	demonstrateFirstTimeCreation()
}

// demonstrateFirstTimeCreation shows how the system handles first-time database creation
func demonstrateFirstTimeCreation() {
	// This will create the database and directory structure if they don't exist
	dbPath := "./data/first_time.db"

	// Remove the database file if it exists to simulate first-time creation
	os.Remove(dbPath)
	os.RemoveAll(filepath.Dir(dbPath))

	fmt.Printf("Creating new database at: %s\n", dbPath)

	index, err := index.NewPersistedSimpleIndexWithDatabaseAndLoad(dbPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create first-time database")
	}
	defer index.Close()

	// Verify it's empty
	isEmpty, err := index.IsDatabaseEmpty()
	if err != nil {
		log.Error().Err(err).Msg("Failed to check database")
	} else {
		fmt.Printf("New database is empty: %v\n", isEmpty)
	}

	// Add some initial data
	doc := models.Document{
		ID:     "first_doc",
		Text:   "This is the first document in the new database",
		Source: "/first/path",
		Vector: []float64{0.1, 0.2},
		Meta:   map[string]string{"type": "first"},
	}

	if err := index.AddDocument(doc); err != nil {
		log.Error().Err(err).Msg("Failed to add first document")
	} else {
		fmt.Println("Added first document to new database")
	}

	// Verify data is searchable
	results, err := index.Search("first document")
	if err != nil {
		log.Error().Err(err).Msg("Failed to search")
	} else {
		fmt.Printf("Found %d documents in new database\n", len(results))
	}

	fmt.Println("First-time database creation completed successfully!")
}

