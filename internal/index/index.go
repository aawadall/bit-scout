package index

import (
	"github.com/aawadall/bit-scout/internal/models"
)

/* Index Interface */

type Index interface {
	// Configures the index
	Configure(config map[string]interface{}) error
	// Shows the current index configuration
	ShowConfig() (map[string]interface{}, error)
	// Adds document to current index
	AddDocument(models.Document) error
	// Adds multiple documents to current index
	AddDocuments([]models.Document) error
	// Searches for documents matching the query
	Search(query string) ([]models.Document, error)
	// Deletes a document from the index
	DeleteDocument(id string) error
	// Deletes multiple documents from the index
	DeleteDocuments([]string) error
	// Updates a document in the index
	UpdateDocument(id string, document models.Document) error
	// Updates multiple documents in the index
	UpdateDocuments([]models.Document) error
	// Closes the index
	Close() error
	// Flushes the index to disk
	Flush() error
	// Optimizes the index for faster search
	Optimize() error
	// Returns the number of documents in the index
	Count() (int, error)
	// Returns the size of the index in bytes
	Size() (int, error)
}
