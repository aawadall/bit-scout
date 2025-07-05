package features

import (
	"fmt"
	"strconv"
	"strings"
)

// ConfigBuilder provides a fluent interface for building extractor configurations
type ConfigBuilder struct {
	config ExtractorConfig
}

// NewConfigBuilder creates a new configuration builder
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: ExtractorConfig{
			Enabled:    true,
			Weight:     1.0,
			Parameters: make(map[string]interface{}),
			FeatureMap: make(map[string]string),
			Normalize:  true,
			Vectorize:  true,
		},
	}
}

// Enabled sets whether the extractor is enabled
func (b *ConfigBuilder) Enabled(enabled bool) *ConfigBuilder {
	b.config.Enabled = enabled
	return b
}

// Weight sets the global weight for all features from this extractor
func (b *ConfigBuilder) Weight(weight float64) *ConfigBuilder {
	b.config.Weight = weight
	return b
}

// Parameter sets a specific parameter for the extractor
func (b *ConfigBuilder) Parameter(key string, value interface{}) *ConfigBuilder {
	b.config.Parameters[key] = value
	return b
}

// Parameters sets multiple parameters at once
func (b *ConfigBuilder) Parameters(params map[string]interface{}) *ConfigBuilder {
	for key, value := range params {
		b.config.Parameters[key] = value
	}
	return b
}

// FeatureMap sets a mapping from internal feature names to output names
func (b *ConfigBuilder) FeatureMap(mapping map[string]string) *ConfigBuilder {
	b.config.FeatureMap = mapping
	return b
}

// MapFeature maps a single internal feature name to an output name
func (b *ConfigBuilder) MapFeature(internal, output string) *ConfigBuilder {
	b.config.FeatureMap[internal] = output
	return b
}

// Normalize sets whether to normalize numeric features
func (b *ConfigBuilder) Normalize(normalize bool) *ConfigBuilder {
	b.config.Normalize = normalize
	return b
}

// Vectorize sets whether to include features in vector representation
func (b *ConfigBuilder) Vectorize(vectorize bool) *ConfigBuilder {
	b.config.Vectorize = vectorize
	return b
}

// Build returns the final configuration
func (b *ConfigBuilder) Build() ExtractorConfig {
	return b.config
}

// PresetConfigs provides common configuration presets
type PresetConfigs struct{}

// NewPresetConfigs creates a new preset configurations helper
func NewPresetConfigs() *PresetConfigs {
	return &PresetConfigs{}
}

// Minimal creates a minimal configuration with only basic features
func (p *PresetConfigs) Minimal() ExtractorConfig {
	return NewConfigBuilder().
		Weight(0.5).
		Parameters(map[string]interface{}{
			"include_content_features":   false,
			"include_timestamp_features": false,
		}).
		Build()
}

// Standard creates a standard configuration with most features enabled
func (p *PresetConfigs) Standard() ExtractorConfig {
	return NewConfigBuilder().
		Weight(1.0).
		Parameters(map[string]interface{}{
			"include_content_features":   true,
			"include_timestamp_features": true,
		}).
		Build()
}

// Comprehensive creates a comprehensive configuration with all features and high weight
func (p *PresetConfigs) Comprehensive() ExtractorConfig {
	return NewConfigBuilder().
		Weight(1.5).
		Parameters(map[string]interface{}{
			"include_content_features":    true,
			"include_timestamp_features":  true,
			"include_path_features":       true,
			"include_permission_features": true,
		}).
		Build()
}

// Custom creates a custom configuration from a string specification
func (p *PresetConfigs) Custom(spec string) (ExtractorConfig, error) {
	builder := NewConfigBuilder()

	parts := strings.Split(spec, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		keyValue := strings.SplitN(part, "=", 2)
		if len(keyValue) != 2 {
			return ExtractorConfig{}, fmt.Errorf("invalid config spec format: %s", part)
		}

		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])

		switch key {
		case "enabled":
			builder.Enabled(value == "true")
		case "weight":
			if weight, err := parseFloat(value); err == nil {
				builder.Weight(weight)
			} else {
				return ExtractorConfig{}, fmt.Errorf("invalid weight value: %s", value)
			}
		case "normalize":
			builder.Normalize(value == "true")
		case "vectorize":
			builder.Vectorize(value == "true")
		default:
			// Treat as parameter
			builder.Parameter(key, value)
		}
	}

	return builder.Build(), nil
}

// RegistryConfig holds configuration for the entire feature registry
type RegistryConfig struct {
	Extractors map[string]ExtractorConfig // Configuration for each extractor
	Global     GlobalConfig               // Global settings
}

// GlobalConfig holds global configuration for the feature registry
type GlobalConfig struct {
	DefaultWeight float64                // Default weight for extractors
	DefaultParams map[string]interface{} // Default parameters for extractors
	LogLevel      string                 // Logging level for feature extraction
}

// NewRegistryConfig creates a new registry configuration
func NewRegistryConfig() *RegistryConfig {
	return &RegistryConfig{
		Extractors: make(map[string]ExtractorConfig),
		Global: GlobalConfig{
			DefaultWeight: 1.0,
			DefaultParams: make(map[string]interface{}),
			LogLevel:      "info",
		},
	}
}

// SetExtractorConfig sets configuration for a specific extractor
func (rc *RegistryConfig) SetExtractorConfig(name string, config ExtractorConfig) {
	rc.Extractors[name] = config
}

// GetExtractorConfig gets configuration for a specific extractor
func (rc *RegistryConfig) GetExtractorConfig(name string) (ExtractorConfig, bool) {
	config, exists := rc.Extractors[name]
	return config, exists
}

// ApplyToRegistry applies this configuration to a feature registry
func (rc *RegistryConfig) ApplyToRegistry(registry *FeatureRegistry) error {
	// Apply global defaults
	for name := range registry.extractors {
		config, exists := rc.Extractors[name]
		if !exists {
			// Use global defaults
			config = ExtractorConfig{
				Enabled:    true,
				Weight:     rc.Global.DefaultWeight,
				Parameters: make(map[string]interface{}),
				FeatureMap: make(map[string]string),
				Normalize:  true,
				Vectorize:  true,
			}

			// Copy global default parameters
			for key, value := range rc.Global.DefaultParams {
				config.Parameters[key] = value
			}
		}

		if err := registry.Configure(name, config); err != nil {
			return fmt.Errorf("failed to configure extractor %s: %w", name, err)
		}
	}

	return nil
}

// Helper function to parse float values
func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
