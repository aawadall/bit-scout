package loaders

import (
	"fmt"
	"github.com/aawadall/bit-scout/internal/models"
	"github.com/rs/zerolog/log"
)

// LoaderRegistry manages a set of CorpusLoader plugins.
type LoaderRegistry struct {
	loaders map[string]CorpusLoader
}

// NewLoaderRegistry creates a new LoaderRegistry.
func NewLoaderRegistry() *LoaderRegistry {
	return &LoaderRegistry{loaders: make(map[string]CorpusLoader)}
}

// Register adds a CorpusLoader implementation with a given name.
func (r *LoaderRegistry) Register(name string, loader CorpusLoader) {
	log.Info().Msgf("RegisterLoader: %s", name)
	r.loaders[name] = loader
}

// Get retrieves a registered CorpusLoader by name.
func (r *LoaderRegistry) Get(name string) (CorpusLoader, bool) {
	loader, ok := r.loaders[name]
	return loader, ok
}

// List returns the names of all registered loaders.
func (r *LoaderRegistry) List() []string {
	names := make([]string, 0, len(r.loaders))
	for name := range r.loaders {
		names = append(names, name)
	}
	return names
}

// LoadAll iterates over all registered loaders, calls Load on each with the provided source, and aggregates the results.
func (r *LoaderRegistry) LoadAll() ([]models.Document, error) {
	var allDocs []models.Document
	for name, loader := range r.loaders {
		docs, err := loader.Load()
		if err != nil {
			log.Error().Msgf("LoadAll: loader '%s' failed: %s", name, err)
			continue // skip this loader, but continue with others
		}
		allDocs = append(allDocs, docs...)
	}
	if len(allDocs) == 0 {
		return nil, fmt.Errorf("no documents loaded from any loader")
	}
	return allDocs, nil
}
