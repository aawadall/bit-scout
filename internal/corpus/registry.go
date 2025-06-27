package corpus

import (
	"fmt"

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

// Load All
func (r *LoaderRegistry) LoadAll(source string) ([]Document, error) {
	loader, ok := r.Get(source)
	if !ok {
		err := fmt.Errorf("loader not found: %s", source)
		log.Error().Msgf("LoadAll: %s", err)
		return nil, err
	}
	documents, err := loader.Load(source)
	if err != nil {
		log.Error().Msgf("LoadAll: %s", err)
		return nil, err
	}
	return documents, nil
}