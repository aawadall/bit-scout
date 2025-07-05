package features

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/aawadall/bit-scout/internal/models"
	"github.com/stretchr/testify/assert"
)

// Example demonstrating basic feature extraction usage
func ExampleFeatureRegistry_basic() {
	// Create a temporary test file
	tempDir := os.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "This is a test document with some content.\nIt has multiple lines."

	// Write the test file
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		fmt.Printf("Error creating test file: %v\n", err)
		return
	}
	defer os.Remove(testFile) // Clean up

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

	// Create a test document
	doc := models.Document{
		ID:     "test-doc-1",
		Text:   testContent,
		Source: testFile,
		Meta:   make(map[string]string),
	}

	// Extract features
	featureSets, err := registry.ExtractAll(doc)
	if err != nil {
		fmt.Printf("Error extracting features: %v\n", err)
		return
	}

	// Print results
	for _, featureSet := range featureSets {
		fmt.Printf("Document: %s\n", featureSet.DocumentID)
		fmt.Printf("Features extracted: %d\n", len(featureSet.Features))
		fmt.Printf("Vector length: %d\n", len(featureSet.Vector))

		// Print some key features
		if feature, exists := featureSet.Features["filename"]; exists {
			fmt.Printf("Filename: %v\n", feature.Value)
		}
		if feature, exists := featureSet.Features["word_count"]; exists {
			fmt.Printf("Word count: %v\n", feature.Value)
		}
		if feature, exists := featureSet.Features["line_count"]; exists {
			fmt.Printf("Line count: %v\n", feature.Value)
		}
	}

	// Output:
	// Document: test-doc-1
	// Features extracted: 25
	// Vector length: 20
	// Filename: test.txt
	// Word count: 12
	// Line count: 2
}

// Example demonstrating configuration presets
func ExampleFeatureRegistry_presets() {
	// Create a feature registry
	registry := NewFeatureRegistry()

	// Create and register a filesystem extractor
	fsExtractor := NewFilesystemExtractor()
	registry.Register(fsExtractor)

	// Use preset configurations
	presets := NewPresetConfigs()

	// Minimal configuration
	minimalConfig := presets.Minimal()
	registry.Configure("filesystem", minimalConfig)

	fmt.Printf("Minimal config - Weight: %f, Enabled: %v\n",
		minimalConfig.Weight, minimalConfig.Enabled)

	// Standard configuration
	standardConfig := presets.Standard()
	registry.Configure("filesystem", standardConfig)

	fmt.Printf("Standard config - Weight: %f, Enabled: %v\n",
		standardConfig.Weight, standardConfig.Enabled)

	// Comprehensive configuration
	comprehensiveConfig := presets.Comprehensive()
	registry.Configure("filesystem", comprehensiveConfig)

	fmt.Printf("Comprehensive config - Weight: %f, Enabled: %v\n",
		comprehensiveConfig.Weight, comprehensiveConfig.Enabled)

	// Output:
	// Minimal config - Weight: 0.500000, Enabled: true
	// Standard config - Weight: 1.000000, Enabled: true
	// Comprehensive config - Weight: 1.500000, Enabled: true
}

// Example demonstrating custom configuration
func ExampleFeatureRegistry_custom() {
	// Create a feature registry
	registry := NewFeatureRegistry()

	// Create and register a filesystem extractor
	fsExtractor := NewFilesystemExtractor()
	registry.Register(fsExtractor)

	// Use custom configuration from string
	presets := NewPresetConfigs()
	customConfig, err := presets.Custom("weight=2.0,normalize=false,vectorize=true")
	if err != nil {
		fmt.Printf("Error parsing custom config: %v\n", err)
		return
	}

	registry.Configure("filesystem", customConfig)

	fmt.Printf("Custom config - Weight: %f, Normalize: %v, Vectorize: %v\n",
		customConfig.Weight, customConfig.Normalize, customConfig.Vectorize)

	// Output:
	// Custom config - Weight: 2.000000, Normalize: false, Vectorize: true
}

// Example demonstrating feature mapping
func ExampleFeatureRegistry_featureMapping() {
	// Create a temporary test file
	tempDir := os.TempDir()
	testFile := filepath.Join(tempDir, "example.txt")
	testContent := "Another test document"

	// Write the test file
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		fmt.Printf("Error creating test file: %v\n", err)
		return
	}
	defer os.Remove(testFile) // Clean up

	// Create a feature registry
	registry := NewFeatureRegistry()

	// Create and register a filesystem extractor
	fsExtractor := NewFilesystemExtractor()
	registry.Register(fsExtractor)

	// Configure with feature mapping
	config := NewConfigBuilder().
		Weight(1.0).
		MapFeature("filename", "name").
		MapFeature("file_size", "size").
		MapFeature("word_count", "words").
		Build()

	registry.Configure("filesystem", config)

	// Create a test document
	doc := models.Document{
		ID:     "test-doc-2",
		Text:   testContent,
		Source: testFile,
		Meta:   make(map[string]string),
	}

	// Extract features
	featureSets, err := registry.ExtractAll(doc)
	if err != nil {
		fmt.Printf("Error extracting features: %v\n", err)
		return
	}

	// Print mapped features
	for _, featureSet := range featureSets {
		if feature, exists := featureSet.Features["name"]; exists {
			fmt.Printf("Mapped filename: %v\n", feature.Value)
		}
		if feature, exists := featureSet.Features["size"]; exists {
			fmt.Printf("Mapped file size: %v\n", feature.Value)
		}
		if feature, exists := featureSet.Features["words"]; exists {
			fmt.Printf("Mapped word count: %v\n", feature.Value)
		}
	}

	// Output:
	// Mapped filename: example.txt
	// Mapped file size: 20
	// Mapped word count: 3
}

// Test demonstrating the complete feature extraction workflow
func TestFeatureExtractionWorkflow(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "This is a test file.\nIt has multiple lines.\nAnd some content."

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	assert.NoError(t, err)

	// Create a feature registry
	registry := NewFeatureRegistry()

	// Create and register a filesystem extractor
	fsExtractor := NewFilesystemExtractor()
	err = registry.Register(fsExtractor)
	assert.NoError(t, err)

	// Configure the extractor
	config := NewConfigBuilder().
		Weight(1.0).
		Vectorize(true).
		Build()

	err = registry.Configure("filesystem", config)
	assert.NoError(t, err)

	// Create a document
	doc := models.Document{
		ID:     "test-doc",
		Text:   testContent,
		Source: testFile,
		Meta:   make(map[string]string),
	}

	// Extract features
	featureSets, err := registry.ExtractAll(doc)
	assert.NoError(t, err)
	assert.Len(t, featureSets, 1)

	featureSet := featureSets[0]
	assert.Equal(t, "test-doc", featureSet.DocumentID)
	assert.Greater(t, len(featureSet.Features), 0)
	assert.Greater(t, len(featureSet.Vector), 0)

	// Verify specific features
	assert.Contains(t, featureSet.Features, "filename")
	assert.Contains(t, featureSet.Features, "word_count")
	assert.Contains(t, featureSet.Features, "line_count")
	assert.Contains(t, featureSet.Features, "content_length")

	// Verify feature values
	assert.Equal(t, "test.txt", featureSet.Features["filename"].Value)
	assert.Equal(t, 12, featureSet.Features["word_count"].Value) // "This is a test file. It has multiple lines. And some content."
	assert.Equal(t, 3, featureSet.Features["line_count"].Value)
	assert.Equal(t, len(testContent), featureSet.Features["content_length"].Value)

	// Verify vector generation
	assert.Greater(t, len(featureSet.Vector), 0)

	// Test batch extraction
	docs := []models.Document{doc}
	batchResults, err := registry.ExtractAllBatch(docs)
	assert.NoError(t, err)
	assert.Len(t, batchResults, 1)
	assert.Len(t, batchResults[0], 1)
}

// Test demonstrating configuration validation
func TestConfigurationValidation(t *testing.T) {
	// Test valid configuration
	config := NewConfigBuilder().
		Weight(1.0).
		Enabled(true).
		Build()

	fsExtractor := NewFilesystemExtractor()
	err := fsExtractor.Configure(config)
	assert.NoError(t, err)

	// Test invalid configuration (negative weight)
	invalidConfig := NewConfigBuilder().
		Weight(-1.0).
		Build()

	err = fsExtractor.Configure(invalidConfig)
	assert.NoError(t, err) // Configure should succeed

	err = fsExtractor.Validate()
	assert.Error(t, err) // But validation should fail
}

// Test demonstrating feature mapping
func TestFeatureMapping(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Test content"

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	assert.NoError(t, err)

	// Create a feature registry
	registry := NewFeatureRegistry()

	// Create and register a filesystem extractor
	fsExtractor := NewFilesystemExtractor()
	registry.Register(fsExtractor)

	// Configure with feature mapping
	config := NewConfigBuilder().
		MapFeature("filename", "name").
		MapFeature("file_size", "size").
		Build()

	registry.Configure("filesystem", config)

	// Create a test document
	doc := models.Document{
		ID:     "test-doc",
		Text:   testContent,
		Source: testFile,
		Meta:   make(map[string]string),
	}

	// Extract features
	featureSets, err := registry.ExtractAll(doc)
	assert.NoError(t, err)
	assert.Len(t, featureSets, 1)

	featureSet := featureSets[0]

	// Verify mapped features exist
	assert.Contains(t, featureSet.Features, "name")
	assert.Contains(t, featureSet.Features, "size")

	// Verify original features don't exist
	assert.NotContains(t, featureSet.Features, "filename")
	assert.NotContains(t, featureSet.Features, "file_size")

	// Verify mapped values are correct
	assert.Equal(t, "test.txt", featureSet.Features["name"].Value)
	assert.Equal(t, int64(len(testContent)), featureSet.Features["size"].Value)
}
