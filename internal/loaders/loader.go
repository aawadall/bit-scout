package loaders

import (
	"github.com/aawadall/bit-scout/internal/models"
)

/*
Interface for corpus loader, and registry of loaders.
*/

// CorpusLoader defines the interface for loading documents from a source.
type CorpusLoader interface {
	// Load loads documents.
	Load() ([]models.Document, error)
}
