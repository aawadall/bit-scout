# Feature Extraction System

The feature extraction system provides a configurable, extensible framework for extracting various features from documents. It supports multiple extractor types, configuration presets, and feature mapping.

## Overview

The feature extraction system consists of:

- **FeatureExtractor Interface**: Defines the contract for all feature extractors
- **FeatureRegistry**: Manages multiple extractors and their configurations
- **FilesystemExtractor**: Extracts filesystem-related features from documents
- **Configuration System**: Provides flexible configuration options including presets and custom configurations

## Core Concepts

### Feature
A single extracted characteristic from a document:
```go
type Feature struct {
    Name   string      // Feature identifier
    Value  interface{} // Feature value (string, number, bool, etc.)
    Type   string      // Feature type ("string", "number", "boolean", "vector")
    Weight float64     // Feature weight (default: 1.0)
}
```

### FeatureSet
A collection of features extracted from a document:
```go
type FeatureSet struct {
    DocumentID string             // Document ID
    Features   map[string]Feature // Map of feature name to feature
    Vector     []float64          // Vector representation of features
}
```

### ExtractorConfig
Configuration for a feature extractor:
```go
type ExtractorConfig struct {
    Enabled     bool                   // Whether extractor is enabled
    Weight      float64                // Global weight for all features
    Parameters  map[string]interface{} // Extractor-specific parameters
    FeatureMap  map[string]string      // Feature name mapping
    Normalize   bool                   // Normalize numeric features
    Vectorize   bool                   // Include in vector representation
}
```

## Usage

### Basic Usage

```go
// Create a feature registry
registry := NewFeatureRegistry()

// Create and register a filesystem extractor
fsExtractor := NewFilesystemExtractor()
registry.Register(fsExtractor)

// Configure the extractor
config := NewConfigBuilder().
    Weight(1.0).
    Vectorize(true).
    Build()

registry.Configure("filesystem", config)

// Extract features from a document
doc := models.Document{
    ID:     "doc-1",
    Text:   "Document content",
    Source: "/path/to/file.txt",
}

featureSets, err := registry.ExtractAll(doc)
if err != nil {
    log.Fatal(err)
}

// Process extracted features
for _, featureSet := range featureSets {
    fmt.Printf("Document: %s\n", featureSet.DocumentID)
    fmt.Printf("Features: %d\n", len(featureSet.Features))
    fmt.Printf("Vector length: %d\n", len(featureSet.Vector))
}
```

### Configuration Presets

The system provides common configuration presets:

```go
presets := NewPresetConfigs()

// Minimal configuration (basic features only)
minimalConfig := presets.Minimal()

// Standard configuration (most features enabled)
standardConfig := presets.Standard()

// Comprehensive configuration (all features, high weight)
comprehensiveConfig := presets.Comprehensive()

// Apply to registry
registry.Configure("filesystem", standardConfig)
```

### Custom Configuration

Create custom configurations using the builder pattern:

```go
config := NewConfigBuilder().
    Weight(2.0).
    Enabled(true).
    Normalize(false).
    Vectorize(true).
    Parameter("include_content_features", true).
    Parameter("include_timestamp_features", false).
    MapFeature("filename", "name").
    MapFeature("file_size", "size").
    Build()

registry.Configure("filesystem", config)
```

### String-based Configuration

Parse configurations from strings:

```go
presets := NewPresetConfigs()
config, err := presets.Custom("weight=2.0,normalize=false,vectorize=true")
if err != nil {
    log.Fatal(err)
}

registry.Configure("filesystem", config)
```

### Feature Mapping

Map internal feature names to custom output names:

```go
config := NewConfigBuilder().
    MapFeature("filename", "name").
    MapFeature("file_size", "size").
    MapFeature("word_count", "words").
    Build()

registry.Configure("filesystem", config)

// Now features will be available as "name", "size", "words" instead of
// "filename", "file_size", "word_count"
```

### Batch Processing

Extract features from multiple documents efficiently:

```go
docs := []models.Document{/* ... */}
batchResults, err := registry.ExtractAllBatch(docs)
if err != nil {
    log.Fatal(err)
}

// batchResults is [][]*FeatureSet - one slice per extractor
for extractorIdx, featureSets := range batchResults {
    for _, featureSet := range featureSets {
        // Process each feature set
    }
}
```

## Filesystem Extractor

The `FilesystemExtractor` extracts various filesystem-related features:

### Basic File Information
- `filename`: File name
- `extension`: File extension
- `path`: Full file path
- `directory`: Directory containing the file

### File Size Features
- `file_size`: File size in bytes
- `file_size_kb`: File size in kilobytes
- `file_size_mb`: File size in megabytes

### Timestamp Features
- `modified_time`: Last modification time (RFC3339 format)
- `modified_unix`: Last modification time (Unix timestamp)
- `modified_year`: Year of last modification
- `modified_month`: Month of last modification
- `modified_day`: Day of last modification

### File Mode Features
- `is_directory`: Whether file is a directory
- `is_regular_file`: Whether file is a regular file
- `is_symlink`: Whether file is a symbolic link
- `is_executable`: Whether file is executable
- `is_writable`: Whether file is writable
- `is_readable`: Whether file is readable
- `is_hidden`: Whether file is hidden (starts with ".")
- `is_system`: Whether file has system flag
- `is_archive`: Whether file has archive flag

### Content Features
- `content_length`: Length of document content
- `line_count`: Number of lines in content
- `word_count`: Number of words in content

### Path Features
- `path_depth`: Depth of file path

## Vector Generation

When `Vectorize` is enabled, the extractor generates a vector representation of features:

- **Numeric features**: Added directly to vector (weighted)
- **Boolean features**: Added as 0/1 values (weighted)
- **String features**: Not included in vector

Example vector structure:
```
[file_size, file_size_kb, file_size_mb, modified_unix, modified_year, 
 modified_month, modified_day, content_length, line_count, word_count, 
 path_depth, is_directory, is_regular_file, is_symlink, is_executable, 
 is_writable, is_readable, is_hidden, is_system, is_archive]
```

## Extending the System

### Creating Custom Extractors

Implement the `FeatureExtractor` interface:

```go
type MyExtractor struct {
    config ExtractorConfig
}

func (e *MyExtractor) Name() string {
    return "my_extractor"
}

func (e *MyExtractor) Configure(config ExtractorConfig) error {
    e.config = config
    return nil
}

func (e *MyExtractor) GetConfig() ExtractorConfig {
    return e.config
}

func (e *MyExtractor) Extract(doc models.Document) (*FeatureSet, error) {
    // Extract features from document
    features := make(map[string]Feature)
    
    // Add your feature extraction logic here
    
    return &FeatureSet{
        DocumentID: doc.ID,
        Features:   features,
        Vector:     e.generateVector(features),
    }, nil
}

func (e *MyExtractor) ExtractBatch(docs []models.Document) ([]*FeatureSet, error) {
    var results []*FeatureSet
    for _, doc := range docs {
        featureSet, err := e.Extract(doc)
        if err != nil {
            continue
        }
        results = append(results, featureSet)
    }
    return results, nil
}

func (e *MyExtractor) GetSupportedFeatures() []string {
    return []string{"feature1", "feature2", "feature3"}
}

func (e *MyExtractor) Validate() error {
    if e.config.Weight < 0 {
        return fmt.Errorf("weight must be non-negative")
    }
    return nil
}

func (e *MyExtractor) generateVector(features map[string]Feature) []float64 {
    // Generate vector representation
    return []float64{}
}
```

### Registering Custom Extractors

```go
registry := NewFeatureRegistry()
myExtractor := &MyExtractor{}
registry.Register(myExtractor)

config := NewConfigBuilder().Weight(1.0).Build()
registry.Configure("my_extractor", config)
```

## Configuration Examples

### Minimal Configuration
```go
config := NewConfigBuilder().
    Weight(0.5).
    Parameter("include_content_features", false).
    Parameter("include_timestamp_features", false).
    Build()
```

### Standard Configuration
```go
config := NewConfigBuilder().
    Weight(1.0).
    Parameter("include_content_features", true).
    Parameter("include_timestamp_features", true).
    Build()
```

### Comprehensive Configuration
```go
config := NewConfigBuilder().
    Weight(1.5).
    Parameter("include_content_features", true).
    Parameter("include_timestamp_features", true).
    Parameter("include_path_features", true).
    Parameter("include_permission_features", true).
    Build()
```

### Custom Feature Mapping
```go
config := NewConfigBuilder().
    Weight(1.0).
    MapFeature("filename", "name").
    MapFeature("file_size", "size").
    MapFeature("word_count", "words").
    MapFeature("line_count", "lines").
    Build()
```

## Error Handling

The system provides graceful error handling:

- **Extractor failures**: Individual extractor failures don't stop the entire process
- **Configuration errors**: Invalid configurations are caught during validation
- **File access errors**: Filesystem extractor handles missing or inaccessible files
- **Batch processing**: Failed documents are skipped, others continue processing

## Performance Considerations

- **Batch processing**: Use `ExtractBatch` for multiple documents
- **Feature selection**: Disable unused features to improve performance
- **Vector generation**: Disable vectorization if not needed
- **Configuration caching**: Reuse configurations when possible

## Testing

Run the tests to verify functionality:

```bash
go test ./internal/features/...
```

The test suite includes:
- Basic feature extraction workflow
- Configuration validation
- Feature mapping
- Batch processing
- Error handling scenarios 