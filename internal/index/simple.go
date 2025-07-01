package index

import (
	"fmt"
	"strings"

	"github.com/aawadall/bit-scout/internal/models"
	"github.com/rs/zerolog/log"
)

// SimpleIndex is a basic in-memory index implementation
type SimpleIndex struct {
	documents map[string]models.Document
	config    map[string]interface{}
}

// NewSimpleIndex creates a new SimpleIndex instance
func NewSimpleIndex() *SimpleIndex {
	return &SimpleIndex{
		documents: make(map[string]models.Document),
		config:    make(map[string]interface{}),
	}
}

// Configure sets the index configuration
func (idx *SimpleIndex) Configure(config map[string]interface{}) error {
	idx.config = config
	log.Info().Msgf("SimpleIndex configured with %d settings", len(config))
	return nil
}

// ShowConfig returns the current index configuration
func (idx *SimpleIndex) ShowConfig() (map[string]interface{}, error) {
	// Return a copy of the config to prevent external modification
	configCopy := make(map[string]interface{})
	for key, value := range idx.config {
		configCopy[key] = value
	}
	log.Debug().Msgf("SimpleIndex config requested, returning %d settings", len(configCopy))
	return configCopy, nil
}

// AddDocument adds a single document to the index
func (idx *SimpleIndex) AddDocument(doc models.Document) error {
	idx.documents[doc.ID] = doc
	log.Debug().Msgf("Added document %s to index", doc.ID)
	return nil
}

// AddDocuments adds multiple documents to the index
func (idx *SimpleIndex) AddDocuments(docs []models.Document) error {
	for _, doc := range docs {
		if err := idx.AddDocument(doc); err != nil {
			return err
		}
	}
	log.Info().Msgf("Added %d documents to index", len(docs))
	return nil
}

// Search performs advanced query search with boolean operations and dimension filtering
func (idx *SimpleIndex) Search(query string) ([]models.Document, error) {
	if query == "" {
		return []models.Document{}, nil
	}

	// Try to parse as advanced query first
	parsedQuery, err := ParseQuery(query)
	if err == nil && len(parsedQuery.Conditions) > 0 {
		// Use advanced query evaluation
		return idx.searchAdvanced(parsedQuery)
	}

	// Fall back to simple text search for backward compatibility
	return idx.searchSimple(query)
}

// searchAdvanced performs search using parsed query conditions
func (idx *SimpleIndex) searchAdvanced(query *Query) ([]models.Document, error) {
	var results []models.Document

	for _, doc := range idx.documents {
		matches, err := query.Evaluate(doc)
		if err != nil {
			log.Warn().Msgf("Error evaluating query for document %s: %s", doc.ID, err)
			continue
		}

		if matches {
			results = append(results, doc)
		}
	}

	log.Info().Msgf("Advanced search for '%s' returned %d results", query.RawQuery, len(results))
	return results, nil
}

// searchSimple performs the original simple text search
func (idx *SimpleIndex) searchSimple(query string) ([]models.Document, error) {
	query = strings.ToLower(query)
	var results []models.Document

	for _, doc := range idx.documents {
		// Search in document text
		if strings.Contains(strings.ToLower(doc.Text), query) {
			results = append(results, doc)
			continue
		}

		// Search in metadata
		for key, value := range doc.Meta {
			if strings.Contains(strings.ToLower(key), query) ||
				strings.Contains(strings.ToLower(value), query) {
				results = append(results, doc)
				break
			}
		}

		// Search in source path
		if strings.Contains(strings.ToLower(doc.Source), query) {
			results = append(results, doc)
		}
	}

	log.Info().Msgf("Simple search for '%s' returned %d results", query, len(results))
	return results, nil
}

// DeleteDocument removes a document from the index
func (idx *SimpleIndex) DeleteDocument(id string) error {
	if _, exists := idx.documents[id]; !exists {
		return fmt.Errorf("document %s not found in index", id)
	}
	delete(idx.documents, id)
	log.Debug().Msgf("Deleted document %s from index", id)
	return nil
}

// DeleteDocuments removes multiple documents from the index
func (idx *SimpleIndex) DeleteDocuments(ids []string) error {
	for _, id := range ids {
		if err := idx.DeleteDocument(id); err != nil {
			return err
		}
	}
	log.Info().Msgf("Deleted %d documents from index", len(ids))
	return nil
}

// UpdateDocument updates an existing document in the index
func (idx *SimpleIndex) UpdateDocument(id string, doc models.Document) error {
	if _, exists := idx.documents[id]; !exists {
		return fmt.Errorf("document %s not found in index", id)
	}
	idx.documents[id] = doc
	log.Debug().Msgf("Updated document %s in index", id)
	return nil
}

// UpdateDocuments updates multiple documents in the index
func (idx *SimpleIndex) UpdateDocuments(docs []models.Document) error {
	for _, doc := range docs {
		if err := idx.UpdateDocument(doc.ID, doc); err != nil {
			return err
		}
	}
	log.Info().Msgf("Updated %d documents in index", len(docs))
	return nil
}

// Close performs cleanup operations
func (idx *SimpleIndex) Close() error {
	log.Info().Msg("SimpleIndex closed")
	return nil
}

// Flush writes the index to disk (not implemented for simple in-memory index)
func (idx *SimpleIndex) Flush() error {
	log.Info().Msg("SimpleIndex flush called (no-op for in-memory index)")
	return nil
}

// Optimize optimizes the index for faster search (not implemented for simple index)
func (idx *SimpleIndex) Optimize() error {
	log.Info().Msg("SimpleIndex optimize called (no-op for in-memory index)")
	return nil
}

// Count returns the number of documents in the index
func (idx *SimpleIndex) Count() (int, error) {
	return len(idx.documents), nil
}

// Size returns the approximate size of the index in bytes
func (idx *SimpleIndex) Size() (int, error) {
	size := 0
	for _, doc := range idx.documents {
		size += len(doc.ID)
		size += len(doc.Text)
		size += len(doc.Source)
		for key, value := range doc.Meta {
			size += len(key)
			size += len(value)
		}
		size += len(doc.Vector) * 8 // 8 bytes per float64
	}
	return size, nil
}
