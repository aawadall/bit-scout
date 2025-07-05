package features

import (
	"fmt"

	"github.com/aawadall/bit-scout/internal/models"
	"github.com/rs/zerolog/log"
)

// Feature represents a single extracted feature from a document
type Feature struct {
	Name   string      // Name/identifier of the feature
	Value  interface{} // The feature value (can be string, number, bool, etc.)
	Type   string      // Type of the feature (e.g., "string", "number", "boolean", "vector")
	Weight float64     // Optional weight for the feature (default: 1.0)
}

// FeatureSet represents a collection of features extracted from a document
type FeatureSet struct {
	DocumentID string             // ID of the document these features belong to
	Features   map[string]Feature // Map of feature name to feature
	Vector     []float64          // Optional vector representation of all features
}

// ExtractorConfig holds configuration for a feature extractor
type ExtractorConfig struct {
	Enabled    bool                   // Whether this extractor is enabled
	Weight     float64                // Global weight for all features from this extractor
	Parameters map[string]interface{} // Extractor-specific parameters
	FeatureMap map[string]string      // Optional mapping of internal feature names to output names
	Normalize  bool                   // Whether to normalize numeric features
	Vectorize  bool                   // Whether to include features in vector representation
}

// FeatureExtractor defines the interface for extracting features from documents
type FeatureExtractor interface {
	// Name returns the name of this extractor
	Name() string

	// Configure sets the configuration for this extractor
	Configure(config ExtractorConfig) error

	// GetConfig returns the current configuration
	GetConfig() ExtractorConfig

	// Extract extracts features from a single document
	Extract(doc models.Document) (*FeatureSet, error)

	// ExtractBatch extracts features from multiple documents (for efficiency)
	ExtractBatch(docs []models.Document) ([]*FeatureSet, error)

	// GetSupportedFeatures returns a list of feature names this extractor can produce
	GetSupportedFeatures() []string

	// Validate checks if the extractor is properly configured
	Validate() error
}

// FeatureRegistry manages multiple feature extractors
type FeatureRegistry struct {
	extractors map[string]FeatureExtractor
	configs    map[string]ExtractorConfig
}

// NewFeatureRegistry creates a new feature registry
func NewFeatureRegistry() *FeatureRegistry {
	return &FeatureRegistry{
		extractors: make(map[string]FeatureExtractor),
		configs:    make(map[string]ExtractorConfig),
	}
}

// Register adds a feature extractor to the registry
func (r *FeatureRegistry) Register(extractor FeatureExtractor) error {
	name := extractor.Name()
	if _, exists := r.extractors[name]; exists {
		return fmt.Errorf("extractor %s already registered", name)
	}

	r.extractors[name] = extractor
	log.Info().Msgf("Registered feature extractor: %s", name)
	return nil
}

// Configure sets configuration for a specific extractor
func (r *FeatureRegistry) Configure(extractorName string, config ExtractorConfig) error {
	extractor, exists := r.extractors[extractorName]
	if !exists {
		return fmt.Errorf("extractor %s not found", extractorName)
	}

	if err := extractor.Configure(config); err != nil {
		return fmt.Errorf("failed to configure extractor %s: %w", extractorName, err)
	}

	r.configs[extractorName] = config
	log.Info().Msgf("Configured extractor %s with %d parameters", extractorName, len(config.Parameters))
	return nil
}

// ExtractAll extracts features from a document using all enabled extractors
func (r *FeatureRegistry) ExtractAll(doc models.Document) ([]*FeatureSet, error) {
	var results []*FeatureSet

	for name, extractor := range r.extractors {
		config := r.configs[name]
		if !config.Enabled {
			continue
		}

		featureSet, err := extractor.Extract(doc)
		if err != nil {
			log.Warn().Err(err).Msgf("Failed to extract features from %s using %s", doc.ID, name)
			continue
		}

		results = append(results, featureSet)
	}

	return results, nil
}

// ExtractAllBatch extracts features from multiple documents using all enabled extractors
func (r *FeatureRegistry) ExtractAllBatch(docs []models.Document) ([][]*FeatureSet, error) {
	var results [][]*FeatureSet

	for name, extractor := range r.extractors {
		config := r.configs[name]
		if !config.Enabled {
			continue
		}

		featureSets, err := extractor.ExtractBatch(docs)
		if err != nil {
			log.Warn().Err(err).Msgf("Failed to extract features using %s", name)
			continue
		}

		results = append(results, featureSets)
	}

	return results, nil
}

// GetExtractor returns a specific extractor by name
func (r *FeatureRegistry) GetExtractor(name string) (FeatureExtractor, bool) {
	extractor, exists := r.extractors[name]
	return extractor, exists
}

// ListExtractors returns all registered extractor names
func (r *FeatureRegistry) ListExtractors() []string {
	var names []string
	for name := range r.extractors {
		names = append(names, name)
	}
	return names
}

// GetEnabledExtractors returns names of all enabled extractors
func (r *FeatureRegistry) GetEnabledExtractors() []string {
	var names []string
	for name, config := range r.configs {
		if config.Enabled {
			names = append(names, name)
		}
	}
	return names
}
