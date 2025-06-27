package corpus

import "github.com/rs/zerolog/log"

// Document represents a single document loaded from a corpus source.
type Document struct {
	ID     string
	Text   string
	Source string            // Source of the document (e.g., file path, URL)
	Vector []float64         // Vector representation of the document
	Meta   map[string]string // Optional metadata (e.g., filename, tags)
}

// Print the document
func (d *Document) Print() {
	log.Info().Msgf("Document ID: %s", d.ID)
	log.Info().Msgf("Document Text: %s", d.Text)
	log.Info().Msgf("Document Source: %s", d.Source)
	log.Info().Msgf("Document Metadata:")
	for key, value := range d.Meta {
		log.Info().Msgf("  %s: %s", key, value)
	}
}