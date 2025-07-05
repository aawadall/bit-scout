package features

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aawadall/bit-scout/internal/models"
	"github.com/rs/zerolog/log"
)

// FilesystemExtractor extracts filesystem-related features from documents
type FilesystemExtractor struct {
	config ExtractorConfig
}

// NewFilesystemExtractor creates a new filesystem feature extractor
func NewFilesystemExtractor() *FilesystemExtractor {
	return &FilesystemExtractor{
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

// Name returns the name of this extractor
func (e *FilesystemExtractor) Name() string {
	return "filesystem"
}

// Configure sets the configuration for this extractor
func (e *FilesystemExtractor) Configure(config ExtractorConfig) error {
	e.config = config
	log.Debug().Msgf("FilesystemExtractor configured with enabled=%v, weight=%f", config.Enabled, config.Weight)
	return nil
}

// GetConfig returns the current configuration
func (e *FilesystemExtractor) GetConfig() ExtractorConfig {
	return e.config
}

// Extract extracts filesystem features from a single document
func (e *FilesystemExtractor) Extract(doc models.Document) (*FeatureSet, error) {
	if !e.config.Enabled {
		return &FeatureSet{
			DocumentID: doc.ID,
			Features:   make(map[string]Feature),
			Vector:     []float64{},
		}, nil
	}

	// Get file info from the source path
	info, err := os.Stat(doc.Source)
	if err != nil {
		return nil, err
	}

	features := make(map[string]Feature)

	// Extract basic file information
	features["filename"] = Feature{
		Name:   "filename",
		Value:  info.Name(),
		Type:   "string",
		Weight: e.config.Weight,
	}

	features["extension"] = Feature{
		Name:   "extension",
		Value:  filepath.Ext(info.Name()),
		Type:   "string",
		Weight: e.config.Weight,
	}

	features["path"] = Feature{
		Name:   "path",
		Value:  doc.Source,
		Type:   "string",
		Weight: e.config.Weight,
	}

	features["directory"] = Feature{
		Name:   "directory",
		Value:  filepath.Dir(doc.Source),
		Type:   "string",
		Weight: e.config.Weight,
	}

	// Extract file size features
	fileSize := info.Size()
	features["file_size"] = Feature{
		Name:   "file_size",
		Value:  fileSize,
		Type:   "number",
		Weight: e.config.Weight,
	}

	features["file_size_kb"] = Feature{
		Name:   "file_size_kb",
		Value:  fileSize / 1024,
		Type:   "number",
		Weight: e.config.Weight,
	}

	features["file_size_mb"] = Feature{
		Name:   "file_size_mb",
		Value:  float64(fileSize) / (1024 * 1024),
		Type:   "number",
		Weight: e.config.Weight,
	}

	// Extract timestamp features
	modTime := info.ModTime()
	features["modified_time"] = Feature{
		Name:   "modified_time",
		Value:  modTime.Format(time.RFC3339),
		Type:   "string",
		Weight: e.config.Weight,
	}

	features["modified_unix"] = Feature{
		Name:   "modified_unix",
		Value:  modTime.Unix(),
		Type:   "number",
		Weight: e.config.Weight,
	}

	features["modified_year"] = Feature{
		Name:   "modified_year",
		Value:  modTime.Year(),
		Type:   "number",
		Weight: e.config.Weight,
	}

	features["modified_month"] = Feature{
		Name:   "modified_month",
		Value:  int(modTime.Month()),
		Type:   "number",
		Weight: e.config.Weight,
	}

	features["modified_day"] = Feature{
		Name:   "modified_day",
		Value:  modTime.Day(),
		Type:   "number",
		Weight: e.config.Weight,
	}

	// Extract file mode features
	mode := info.Mode()
	features["is_directory"] = Feature{
		Name:   "is_directory",
		Value:  mode.IsDir(),
		Type:   "boolean",
		Weight: e.config.Weight,
	}

	features["is_regular_file"] = Feature{
		Name:   "is_regular_file",
		Value:  mode.IsRegular(),
		Type:   "boolean",
		Weight: e.config.Weight,
	}

	features["is_symlink"] = Feature{
		Name:   "is_symlink",
		Value:  mode&os.ModeSymlink != 0,
		Type:   "boolean",
		Weight: e.config.Weight,
	}

	features["is_executable"] = Feature{
		Name:   "is_executable",
		Value:  mode&0100 != 0,
		Type:   "boolean",
		Weight: e.config.Weight,
	}

	features["is_writable"] = Feature{
		Name:   "is_writable",
		Value:  mode&0200 != 0,
		Type:   "boolean",
		Weight: e.config.Weight,
	}

	features["is_readable"] = Feature{
		Name:   "is_readable",
		Value:  mode&0400 != 0,
		Type:   "boolean",
		Weight: e.config.Weight,
	}

	features["is_hidden"] = Feature{
		Name:   "is_hidden",
		Value:  strings.HasPrefix(info.Name(), "."),
		Type:   "boolean",
		Weight: e.config.Weight,
	}

	features["is_system"] = Feature{
		Name:   "is_system",
		Value:  mode&01000 != 0,
		Type:   "boolean",
		Weight: e.config.Weight,
	}

	features["is_archive"] = Feature{
		Name:   "is_archive",
		Value:  mode&02000 != 0,
		Type:   "boolean",
		Weight: e.config.Weight,
	}

	// Extract content-based features
	contentLength := len(doc.Text)
	features["content_length"] = Feature{
		Name:   "content_length",
		Value:  contentLength,
		Type:   "number",
		Weight: e.config.Weight,
	}

	features["line_count"] = Feature{
		Name:   "line_count",
		Value:  strings.Count(doc.Text, "\n") + 1,
		Type:   "number",
		Weight: e.config.Weight,
	}

	features["word_count"] = Feature{
		Name:   "word_count",
		Value:  len(strings.Fields(doc.Text)),
		Type:   "number",
		Weight: e.config.Weight,
	}

	// Extract path depth
	pathDepth := len(strings.Split(filepath.Clean(doc.Source), string(filepath.Separator)))
	features["path_depth"] = Feature{
		Name:   "path_depth",
		Value:  pathDepth,
		Type:   "number",
		Weight: e.config.Weight,
	}

	// Apply feature mapping if configured
	if len(e.config.FeatureMap) > 0 {
		mappedFeatures := make(map[string]Feature)
		for name, feature := range features {
			if mappedName, exists := e.config.FeatureMap[name]; exists {
				feature.Name = mappedName
				mappedFeatures[mappedName] = feature
			} else {
				mappedFeatures[name] = feature
			}
		}
		features = mappedFeatures
	}

	// Generate vector representation if requested
	var vector []float64
	if e.config.Vectorize {
		vector = e.generateVector(features)
	}

	featureSet := &FeatureSet{
		DocumentID: doc.ID,
		Features:   features,
		Vector:     vector,
	}

	log.Debug().Msgf("Extracted %d filesystem features from document %s", len(features), doc.ID)
	return featureSet, nil
}

// ExtractBatch extracts filesystem features from multiple documents
func (e *FilesystemExtractor) ExtractBatch(docs []models.Document) ([]*FeatureSet, error) {
	var results []*FeatureSet

	for _, doc := range docs {
		featureSet, err := e.Extract(doc)
		if err != nil {
			log.Warn().Err(err).Msgf("Failed to extract features from document %s", doc.ID)
			continue
		}
		results = append(results, featureSet)
	}

	log.Info().Msgf("Extracted filesystem features from %d documents", len(results))
	return results, nil
}

// GetSupportedFeatures returns a list of feature names this extractor can produce
func (e *FilesystemExtractor) GetSupportedFeatures() []string {
	return []string{
		"filename", "extension", "path", "directory",
		"file_size", "file_size_kb", "file_size_mb",
		"modified_time", "modified_unix", "modified_year", "modified_month", "modified_day",
		"is_directory", "is_regular_file", "is_symlink",
		"is_executable", "is_writable", "is_readable",
		"is_hidden", "is_system", "is_archive",
		"content_length", "line_count", "word_count",
		"path_depth",
	}
}

// Validate checks if the extractor is properly configured
func (e *FilesystemExtractor) Validate() error {
	if e.config.Weight < 0 {
		return fmt.Errorf("weight must be non-negative")
	}
	return nil
}

// generateVector creates a vector representation of the features
func (e *FilesystemExtractor) generateVector(features map[string]Feature) []float64 {
	var vector []float64

	// Add numeric features to vector
	numericFeatures := []string{
		"file_size", "file_size_kb", "file_size_mb",
		"modified_unix", "modified_year", "modified_month", "modified_day",
		"content_length", "line_count", "word_count", "path_depth",
	}

	for _, featureName := range numericFeatures {
		if feature, exists := features[featureName]; exists {
			if value, ok := feature.Value.(float64); ok {
				vector = append(vector, value*feature.Weight)
			} else if value, ok := feature.Value.(int64); ok {
				vector = append(vector, float64(value)*feature.Weight)
			} else if value, ok := feature.Value.(int); ok {
				vector = append(vector, float64(value)*feature.Weight)
			}
		}
	}

	// Add boolean features as 0/1 values
	booleanFeatures := []string{
		"is_directory", "is_regular_file", "is_symlink",
		"is_executable", "is_writable", "is_readable",
		"is_hidden", "is_system", "is_archive",
	}

	for _, featureName := range booleanFeatures {
		if feature, exists := features[featureName]; exists {
			if value, ok := feature.Value.(bool); ok {
				if value {
					vector = append(vector, feature.Weight)
				} else {
					vector = append(vector, 0.0)
				}
			}
		}
	}

	return vector
}
