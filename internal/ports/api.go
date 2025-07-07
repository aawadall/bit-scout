package ports

import "github.com/aawadall/bit-scout/internal/models"

// SearchQuery represents a search request (placeholder, expand as needed)
type SearchQuery struct {
	Query string
	// Add more fields as needed (filters, pagination, etc.)
}

// SearchResults represents search results (placeholder, expand as needed)
type SearchResults struct {
	Documents []models.Document
	// Add more fields as needed (scores, pagination, etc.)
}

// Stats represents system or index statistics (placeholder, expand as needed)
type Stats struct {
	NumDocuments int
	// Add more fields as needed (uptime, memory usage, etc.)
}

// APIPort defines the interface for API adapters (driven port)
// This allows plugging in different API implementations (e.g., GraphQL, REST)
type APIPort interface {
	// Name returns the name/type of the API (e.g., "GraphQL", "REST")
	Name() string
	// Start launches the API server (blocking or non-blocking)
	Start() error
	// Stop gracefully shuts down the API server
	Stop() error

	// Search executes a search query and returns results.
	Search(query SearchQuery) (SearchResults, error)
	// Stats returns statistics about the system or index.
	Stats() (Stats, error)
	// Index manually adds a document to the index.
	Index(doc models.Document) error
}
