package engine

import (
	"github.com/aawadall/bit-scout/internal/ports"
)

/**
 * Search Engine - Core
 **/

// EngineCore defines the core of the search engine using hexagonal architecture.
// It exposes ports for registering adapters such as indexes, loaders, configuration, persistence, feature extractors, and cluster management.
type EngineCore struct {
	// Index registry: maps index names to index implementations
	indexes map[string]ports.IndexPort

	// Loader registry: maps loader names to loader implementations
	loaders map[string]ports.LoaderPort

	// Configuration registry: holds configuration for various components
	configs map[string]ports.ConfigPort

	// Persistence registry: maps persistence names to persistence adapters
	persistence map[string]ports.PersistencePort

	// Feature extractor registry: maps extractor names to feature extractor adapters
	featureExtractors map[string]ports.FeatureExtractorPort

	// Cluster management port (optional, for future extension)
	clusterManager ports.ClusterManagerPort

	// API port (only one supported for now)
	api ports.APIPort
}

// NewEngineCore creates a new EngineCore with empty registries.
func NewEngineCore() *EngineCore {
	return &EngineCore{
		indexes:           make(map[string]ports.IndexPort),
		loaders:           make(map[string]ports.LoaderPort),
		configs:           make(map[string]ports.ConfigPort),
		persistence:       make(map[string]ports.PersistencePort),
		featureExtractors: make(map[string]ports.FeatureExtractorPort),
	}
}

// RegisterIndex registers an index adapter.
func (e *EngineCore) RegisterIndex(name string, index ports.IndexPort) {
	e.indexes[name] = index
}

// RegisterLoader registers a loader adapter.
func (e *EngineCore) RegisterLoader(name string, loader ports.LoaderPort) {
	e.loaders[name] = loader
}

// RegisterConfig registers a configuration adapter.
func (e *EngineCore) RegisterConfig(name string, config ports.ConfigPort) {
	e.configs[name] = config
}

// RegisterPersistence registers a persistence adapter.
func (e *EngineCore) RegisterPersistence(name string, persistence ports.PersistencePort) {
	e.persistence[name] = persistence
}

// RegisterFeatureExtractor registers a feature extractor adapter.
func (e *EngineCore) RegisterFeatureExtractor(name string, extractor ports.FeatureExtractorPort) {
	e.featureExtractors[name] = extractor
}

// SetClusterManager sets the cluster manager adapter.
func (e *EngineCore) SetClusterManager(manager ports.ClusterManagerPort) {
	e.clusterManager = manager
}

// RegisterAPI registers an API adapter (only one supported for now)
func (e *EngineCore) RegisterAPI(api ports.APIPort) {
	e.api = api
}
