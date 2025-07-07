package ports

// FeatureExtractorPort defines the interface for feature extractor adapters (driven port)
type FeatureExtractorPort interface {
	ExtractFeatures(doc interface{}) (map[string]interface{}, error)
}
